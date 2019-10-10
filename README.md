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