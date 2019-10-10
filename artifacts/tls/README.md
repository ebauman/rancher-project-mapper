# Sample TLS Certs

Used for testing this server. Mount `tls.crt` and `tls.key` into the Docker container, or pass
them as arguments to a locally-running instance of `rancher-project-mapper`.

Then, b64 encode the `ca.pem` file and pass that to the `caBundle` argument
for the MutatingWebhookConfiguration. 3

## Files

| File | Purpose | 
|------|---------|
| `ca.key` | Private key for certificate authority. Password: `rancher-project-mapper` |
| `ca.pem` | Public certificate for CA |
| `ca.srl` | Serial for CA |
| `README.md` | This file :) |
| `tls.crt` | TLS certificate, valid for `rancher-project-mapper.cattle-system.svc` and `rancher-project-mapper.cattle-system` |
| `tls.csr` | CSR for requesting `tls.crt` |
| `tls.ext` | ext file for openssl containing SAN information |
| `tls.key` | Private key for `tls.crt` |

## Methods

### CA Generation

Create the private key for the CA:

`openssl genrsa -des3 -out ca.key 2048`

Used password: `rancher-project-mapper`

Create the public cert for the CA:

`openssl req -x509 -new -nodes -key ca.key -sha256 -days 730 -out ca.pem`

### TLS Cert Generation

Create the private key:

`openssl genrsa -out tls.key 2048`

Create the CSR:

`openssl req -new -key tls.key -out tls.csr`

Create the ext file:

```text
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = rancher-project-mapper.cattle-system.svc
DNS.2 = rancher-project-mapper.cattle-system
```

Sign the cert:

`openssl x509 -req -in tls.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out tls.crt -days 730 -sha256 -extfile tls.ext`