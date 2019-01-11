#!/bin/bash

# Goal: Generate all the SSL certificates and keys to encript comunications using TLS (Requires openssl)

#Important info:

# Private files:
#    ca.key: Certificaty Authority (CA) private key file (This shouldn't be shared in production)
#    server.key: Server private key, password protected (This shouldn't be shared with users in production)
#    server.pem: Conversion of the server.key into a format gRPC likes (this shouldn't be shared)
#    server.crt Server certificate signed by the CA (This whould be sent back by the CA owner) - keep on server
# Share files:
#    ca.crt Certificate Authority trust certificate (this should be shared with users in production)
#    server.csr Server certificate signing request (This should be shared with with the CA owner)

# Changes on this CN's to match your hsts in your environment if needed.
SERVER_CN=localhost

# Step 1: Generate Certificate Authority + Tryst Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
