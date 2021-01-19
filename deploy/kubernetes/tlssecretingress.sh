# for internal use for tls-secret-ingress.yaml values
# run from project root
# sh ./deploy/kubernetes/tlssecretingress.sh
kubectl create secret tls illuminatingdeposits-rest-secret-tls --key conf/tls/serverkeyto.pem --cert conf/tls/servercrtto.pem