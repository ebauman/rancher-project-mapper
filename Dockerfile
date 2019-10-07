FROM scratch

COPY rancher-project-mapper /

EXPOSE 443

CMD ["/rancher-project-mapper", "--tls-cert-file", "/certs/tls.crt", "--tls-private-key-file", "/certs/tls.key"]