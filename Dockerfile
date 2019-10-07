FROM scratch

COPY rancher-project-mapper /

EXPOSE 443

CMD ["/rancher-project-mapper", "--tls-cert-file", "tls.crt", "--tls-private-key-file", "tls.key"]