package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"net/http"
	"os"
	"path/filepath"
)

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

type admitFunc func(review v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	requestedAdmissionReview := v1beta1.AdmissionReview{}

	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := scheme.Codecs.UniversalDeserializer()

	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Error(err)
		responseAdmissionReview.Response = toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Error(err)
	}

	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

//func serveNamespace(w http.ResponseWriter, r *http.Request) {
//	serve(w, r, admitNamespace)
//}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func buildInClusterConfig() (*rest.Config) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func buildOutOfClusterConfig() (*rest.Config) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config;
}

func getNamespaceWatcherConfig(config *rest.Config) (*NamespaceWatcherConfig) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	cm, err := clientset.CoreV1().ConfigMaps("cattle-system").Get("rancher-namespace-watcher", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// pull out the config
	rawConfig := []byte(cm.Data["config"])
	nwConfig := NamespaceWatcherConfig{}

	if err = json.Unmarshal(rawConfig, &nwConfig); err != nil {
		panic(fmt.Sprintf("invalid configuration: %s", err.Error()))
	}

	return &nwConfig;
}

func main() {
	var config Config
	config.addFlags()
	inCluster := flag.Bool("in-cluster", false, "True if running inside the cluster (default), false otherwise.")
	flag.Parse()

	var restConfig *rest.Config

	if !*inCluster {
		restConfig = buildOutOfClusterConfig()
	} else {
		restConfig = buildInClusterConfig()
	}

	nwConfig := getNamespaceWatcherConfig(restConfig)

	namespaceWatcher := NamespaceWatcher{config: *nwConfig}

	http.HandleFunc("/namespace", namespaceWatcher.serve)

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: configTLS(config),
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err.Error())
	}
}
