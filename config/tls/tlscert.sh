#!/bin/bash
# This script is executed by Dockerfile.openssl

CN=drinnovations.us

# out ca.key
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096

# out ca.crt from ca.key  -- ca.crt used by client, hence shared
openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=${CN}"

# out server.pem format from server.key -- server.pem used by server, hence shared
openssl genrsa -passout pass:1111 -des3 -out server.key 4096
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem

# out server.csr from server.key
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${CN}"

# out server.crt from server.csr and ca.key and ca.crt -- server.crt used by server, hence shared
openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
