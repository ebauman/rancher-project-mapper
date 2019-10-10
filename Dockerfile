FROM alpine:latest

COPY rancher-project-mapper /

EXPOSE 443

ENV RPM_NAMESPACE "cattle-system"
ENV RPM_CONFIGMAP "rancher-project-mapper"
ENV TLS_CERT_FILE "/certs/tls.crt"
ENV TLS_KEY_FILE "/certs/tls.key"
ENV RPM_LOGLEVEL "1"

CMD /rancher-project-mapper -v $RPM_LOGLEVEL --namespace $RPM_NAMESPACE --config-map $RPM_CONFIGMAP --tls-cert-file $TLS_CERT_FILE --tls-private-key-file $TLS_KEY_FILE