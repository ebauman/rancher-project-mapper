FROM scratch

COPY rancher-namespace-watcher /
CMD ["/rancher-namespace-watcher", "--tls-cert-file", "tls.crt", "--tls-private-key-file", "tls.key"]