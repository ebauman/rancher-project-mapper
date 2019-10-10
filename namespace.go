package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	"log"
	"net/http"
	"regexp"
)

type NamespaceWatcher struct {
	config NamespaceWatcherConfig
}

func (nw *NamespaceWatcher) configWatcher(clientset *kubernetes.Clientset, namespace string, cmName string) {
	// first try to get the config
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(cmName, metav1.GetOptions{})
	if err != nil {
		// namespace might not exist, just log it
		klog.Error(fmt.Sprintf("error getting configmap %s: %s", cmName, err))
	}

	config, err := _parseConfigMap(cm.Data["config"])
	if err != nil {
		klog.Error(fmt.Sprintf("error parsing config: %s", err))
	} else {
		nw.config = *config
	}

	watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	resultChan := watcher.ResultChan()

	klog.Info(fmt.Sprintf("watching for configmap changes for map: %s in namespace: %s", cmName, namespace))

	for event := range resultChan {
		cm, ok := event.Object.(*corev1.ConfigMap) // pull out the configmap
		if !ok {
			klog.Error("unexpected type passed to configmap watcher")
		}

		switch event.Type {
		case watch.Modified:
			fallthrough
		case watch.Added:
			// new configmap. is it what we're expecting?
			if cm.Name == cmName {
				// yes, it is. Try and pull config
				config, err := _parseConfigMap(cm.Data["config"])
				if err != nil {
					// config was not what we expected. no bueno.
					klog.Error(err)
				}

				nw.config = *config // set the new config
			}
		case watch.Deleted:
			if cm.Name == cmName {
				// our config just got deleted. probably
				// should tell someone
				klog.Error("config removed: %s", cm.Name)
				nw.config = NamespaceWatcherConfig{} // blanking out
			}
		}
	}
}

func _parseConfigMap(data string) (*NamespaceWatcherConfig, error) {
	rawConfig := []byte(data)
	nwConfig := NamespaceWatcherConfig{}

	if err := json.Unmarshal(rawConfig, &nwConfig); err != nil {
		return nil, err
	}

	return &nwConfig, nil
}

func (nw *NamespaceWatcher) admitNamespace(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("Admitting namespace")

	nsResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

	if ar.Request.Resource != nsResource {
		klog.Errorf("expect resource to be %s, got %s", nsResource, ar.Request.Resource)
		return nil
	}

	raw := ar.Request.Object.Raw
	ns := corev1.Namespace{}
	deserializer := scheme.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &ns); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	// go through and attempt to match rules.
	// if we find one, mutate
	for i := 0; i < len(nw.config.Rules); i++ {
		reg, err := regexp.Compile(nw.config.Rules[i].Regex)
		if err != nil {
			klog.Error("unable to compile regex for rule %s, regex was %s", i, nw.config.Rules[i].Regex)
			continue
		}
		matched := reg.Match([]byte(ns.Name))
		if matched {
			// we have found a regex match, terminate and mutate
			patch := `[{"op": "add", "path": "/metadata/annotations", "value": {"field.cattle.io/projectId": "%s:%s"}}]`
			reviewResponse.Patch = []byte(fmt.Sprintf(patch, nw.config.Rules[i].Cluster, nw.config.Rules[i].Project))
			break
		}
	}

	return &reviewResponse
}

func (nw *NamespaceWatcher) serve(w http.ResponseWriter, r *http.Request) {
	serve(w, r, nw.admitNamespace)
}
