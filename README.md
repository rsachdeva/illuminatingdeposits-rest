# Illuminating Deposits
# All command should be executed from the root directory (illuminatingdeposits) of the project 
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# REST API using JSON for Messages
# Docker Compose Deployment
 
### To start all services without TLS:
Make sure DEPOSITS_WEB_SERVICE_SERVER_TLS=false in docker-compose.api.yml
And execute:
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.api.yml up --build
 
### To start all services with TLS:
Make sure DEPOSITS_WEB_SERVICE_SERVER_TLS=true in docker-compose.api.yml
And execute:
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.api.yml up --build

The --build option is there for any code changes.

### Then Migrate and set up seed data:
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

### Logs of running services (in a separate terminal):
docker-compose -f ./deploy/compose/docker-compose.api.yml logs -f --tail 1  

### Distributed Tracing
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)  

### Profiling
http://localhost:4000/debug/pprof/

### Metrics
http://localhost:4000/debug/vars

### Shutdown 

docker-compose -f ./deploy/compose/docker-compose.api.yml down  
docker-compose -f ./deploy/compose/docker-compose.seed.yml down

### Quick calculations with Same JSON output without actually invoking REST Http Method
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

### Interest Service REST HTTP Methods Invoked:

####Without TLS Client: 
See tools/resteditor/HealthCRUD.http for request examples and sample response.
Use dev env for localhost:3000  
Or go to cmd/client/main.go  
And uncomment any desired function request starting with "withoutTls..."
Make sure to make email unique to avoid error
See Base64EncodedString call line in the file

####TLS Client: 
Go to cmd/client/main.go  
And uncomment any desired function request starting with "tls..."
run cmd/client
Make sure to make email unique to avoid error
See Base64EncodedString call line in the file

# Push Images to Docker Hub

docker build -t rsachdeva/illuminatingdeposits.api:v0.1 -f ./build/Dockerfile.api .  
docker push rsachdeva/illuminatingdeposits.api:v0.1 
docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed .  
docker push rsachdeva/illuminatingdeposits.seed:v0.1  

# Kubernetes Deployment - WIP

kubectl apply -f deploy/kubernetes/traefik-ingress-daemonset-service.yaml 

kubectl apply -f deploy/kubernetes/zipkin-deployment.yaml   
kubectl apply -f deploy/kubernetes/zipkin-service.yaml   
kubectl apply -f deploy/kubernetes/zipkin-ingress.yaml  

kubectl apply -f deploy/kubernetes/postgres-config.yaml 
kubectl apply -f deploy/kubernetes/postgres-stateful.yaml  
kubectl apply -f deploy/kubernetes/postgres-service.yaml  
To connect external tool with postgres to see database internals use:
Use a connection string similar to:
jdbc:postgresql://127.0.0.1:30007/postgres
If still an issue you can try
kubectl port-forward service/postgres 5432:postgres
Now can easily connect using
jdbc:postgresql://localhost:5432/postgres

Access Traefik Dashboard at [http://localhost:3000/dashboard/#/](http://localhost:3000/dashboard/#/)   

### Distributed Tracing with Kubernetes Ingress

Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io/zipkin/](http://zipkin.127.0.0.1.nip.io/zipkin)  

### Shutdown

kubectl delete -f deploy/kubernetes/.

# TLS files
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl . && \
docker run -v $PWD/config/tls:/tls tlscert:v0.1

To see openssl version being used in Docker:
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl . && \
docker run -ti -v $PWD/config/tls:/tls tlscert:v0.1 sh

/tls # openssl version
OpenSSL 1.1.1g  21 Apr 2020

# CLIENT WIP
# curl 
# https://www.digitalocean.com/community/questions/how-to-ping-docker-container-from-another-container-by-name
# https://hub.docker.com/r/curlimages/curl
docker run --network=compose_deposits_shared_network -it -v "$PWD/config/tls:/tlscurl" curlimages/curl \
--request POST 'https://interestsvcserver:3000/v1/health' \
--header 'Content-Type: application/json' --cacert tlscurl/ca.crt


docker-compose -f ./deploy/compose/docker-compose.curl.yml up