# Rancher Project Mapper

A mutating webhook server that adds a project annotation to a namespace depending on regex rules.
It also helps you allow or deny creation of namespaces based on regex rules.
## Components

### Webhook Server

This code here. 

`go get && go build` to get it built.

Arguments:
```text
-v                      Log level (klog)
--namespace             Namespace containing the ConfigMap
--config-map            Name of the ConfigMap
--kubeconfig            The location of kubeconfig. Implies out-of-cluster.
--tls-cert-file         File that contains the TLS certificate for this server. PEM format.
--tls-private-key-file  File that contains the TLS key for this server. PEM format.
```

Can run inside the cluster or outside.

If running inside the cluster, the service account needs to be able to get a 
configmap called `rancher-project-mapper` from the `cattle-system` namespace.

If running outside the cluster, a `kubeconfig` file is required. The server will attempt
to intuit the location of your kubeconfig file (typically, `$HOME/.kube/config`), but
may require manual specification. 