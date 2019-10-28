# Rancher Project Mapper

Rancher Project Mapper is a simple mutating webhook server that helps you automatically map newly-created namespaces to 
Rancher projects based on regex rules. It also allows you to configure regex rules to deny (or allow) creation of namespaces.

## Important Notes

After deploying this helm chart, you **must** create a ConfigMap in the location specified in Values.yaml
(keys `configmap.namespace` and `configmap.name`). 
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
      "defaults": {
        "mapping": {
          "cluster": "c-12345", // optional
          "project": "p-12345" // optional
        },
        "creation": "allow" // optional, default allow
      },
      "rules": {
        "mapping": [
        {
          "regex": "rpm-*",
          "cluster": "c-12345",
          "project": "p-12345"
        }
        ],
        "creation": [
        {
          "regex": "mpr-*",
          "action": "deny"
        }
        ]
      }
    }
```

You may have as many mapping rules as you'd like, as long as they follow the basic scheme:

```text
{
    "regex": "<valid regex>",
    "cluster": "<cluster id, e.g. c-12345>",
    "project": "<project id, e.g. p-12345>"
}
```

You may have as many creation rules as you'd like, as long as they follow the basic scheme:

```text
{
    "regex": "<valid regex>",
    "action": <"allow" || "deny">
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

