package main

import (
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	"net/http"
	"regexp"
)

type NamespaceWatcher struct {
	config NamespaceWatcherConfig
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
