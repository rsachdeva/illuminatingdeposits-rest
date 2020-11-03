#!/bin/bash
# This script is executed by Dockerfile.openssl

CN=localhost

# out ca.key
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096

openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=${CN}"

# out server.pem format from server.key -- server.pem used by server, hence shared
openssl genrsa -passout pass:1111 -des3 -out server.key 4096
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem

# out server.csr from server.key
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${CN}"

# out server.crt from server.csr and ca.key and ca.crt -- server.crt used by server, hence shared
openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt


#SERVER_CN=localhost
#
## Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
#openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
#openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"
#
## Step 2: Generate the Server Private Key (server.key)
#openssl genrsa -passout pass:1111 -des3 -out server.key 4096
#
## Step 3: Get a certificate signing request from the CA (server.csr)
#openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"
#
## Step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt
#openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
#
## Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
#echo 'now for server-key.pem'
#openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server-key.pem
#echo 'now for server-cert.pem'
##https://stackoverflow.com/questions/4691699/how-to-convert-crt-to-pem
#openssl x509 -in server.crt -out server-cert.pem -outform PEM
##openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.crt -out server-cert.pem
