# Rancher Project Mapper

A mutating webhook server that adds a project annotation to a namespace depending on regex rules.

## Components

### Webhook Server

This code here. 

`go get && go build` to get it built.

Arguments:
```text
--in-cluster            Whether the server is running in a cluster or not.
--kubeconfig            Combined with --in-cluster=false, the location of kubeconfig.
--tls-cert-file         File that contains the TLS certificate for this server. PEM format.
--tls-private-key-file  File that contains the TLS key for this server. PEM format.
```

Can run inside the cluster or outside.

If running inside the cluster, the service account needs to be able to get a 
configmap called `rancher-project-mapper` from the `cattle-system` namespace.

If running outside the cluster, a `kubeconfig` file is required. The server will attempt
to intuit the location of your kubeconfig file (typically, `$HOME/.kube/config`), but
may require manual specification. 

### Mutating Webhook Configuration

An object in Kubernetes called `MutatingWebhookConfiguration`. 
Sets up communications between the kube-apiserver, and this mutating webhook server.

Example:

```text
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: "namespacewatcher.1eb100.net"
webhooks:
- name: "namespacewatcher.1eb100.net"
  rules:
  - apiGroups: [""]
    apiVersions: ["v1"]
    operations: ["CREATE"]
    resources: ["namespaces"]
    scope: "Cluster"
  clientConfig:
    caBundle: "LS0t..." <base64 of your ca certificate used to generate tls.crt and tls.key> 
    url: "https://192.168.1.100.xip.io/namespace"
  sideEffects: None
  timeoutSeconds: 10
```

### Configuration

Create a `ConfigMap` called `rancher-project-mapper` inside of the `cattle-system` namespace. 
This should be the namespace of your cluster, not the Rancher cluster itself.

Inside that ConfigMap, create the following:

```text
apiVersion: v1
kind: ConfigMap
metadata:
    name: rancher-project-mapper
namespace: cattle-system
data:
    config: |-
        {
            "rules": [
                {
                    "regex": "namespace(s)?",
                    "cluster": "c-12345",
                    "project": "p-12345"
                }
            ]
        }
```

You may have as many rules as you like. Rules are evaluated in the order they are defined.
Once a rule is matched, the processing will stop and the namespace will be assigned to the
corresponding project. You may optionally define a default rule to act as a catch-all.
This is done in the following manner:

```text
...
"rules": [
    ...,
    {
        "regex": "<regex>",
        "cluster": "c-12345",
        "project": "p-12345",
        "default": true
    }
]
...
```

It is recommended to make the default rule the final rule, but you may define it at any point
in the rules list. 