# Rancher Project Mapper

Rancher Project Mapper is a simple mutating webhook server that helps you automatically
map newly-created namespaces to Rancher projects based on regex rules. 

## Important Notes

After deploying this helm chart, you must create a `ConfigMap` with project mapper configuration.
By default, this ConfigMap should exist in the `cattle-system` namespace, and have the name 
`rancher-project-mapper`. Those values are overridable via chart configuration.
An example of this is here:

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
                "regex": "rpm-*",
                "cluster": "c-4mdbz",
                "project": "p-48p5n"
            }
        ]
    }
```

You may have as many rules as you'd like, as long as they follow the basic scheme:

```text
{
    "regex": "<valid regex>",
    "cluster": "<cluster id, e.g. c-12345>",
    "project": "<project id, e.g. p-12345>"
}
```

## Options

| Option | Purpose | Default |
| ------ | ------- | ------- |
| `replicaCount` | Number of replicas for the deployment | `1` |
| `image.repository` | Image repository | `ebauman/rancher-project-mapper` |
| `image.tag` | Tag of the image to be deployed | `<version corresponding to helm chart version>` |
| `image.pullPolicy` | Pull policy for the image | `IfNotPresent` |
| `nameOverride` | Overridden name of the chart | `""` |
| `fullNameOverride` | Overridden full name of the chart | `""` |
| `tls.mountpath` | Path of the secret mount location | `/certs` |
| `tls.files.crt` | Location of the `tls.crt` file  | `/certs/tls.crt` |
| `tls.files.key` | Location of the `tls.key` file | `/certs/tls.key` |
| `configmap.name` | Name of the ConfigMap holding project mapper configuration | `rancher-project-mapper` |
| `configmap.namespace` | Namespace containing the ConfigMap | `cattle-system` |
| `loglevel` | Log level of the application | `1` |
| `service.type` | Type of the service deployed | `ClusterIP` | 
| `service.port` | Port of the service deployed | `443` |
| `resources` | Any resource requests | `{}` |
| `nodeSelector` | Any node selector configurations | `{}` |
| `tolerations` | Any tolerations of node taints | `[]` |
| `affinity` | Any node affinity rules | `{}` |

It is generally *not* recommended to adjust the service port nor type. When the 
`MutatingWebhookConfiguration` is installed, it references the service created from this 
helm chart. That service is expected to be an `https://` service, thus any adjustment
in port or service configuration may mess up the `MutatingWebhookConfiguration`.

