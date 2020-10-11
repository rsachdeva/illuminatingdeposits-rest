# Illuminating Deposits
# All command should be executed from the root directory (illuminatingdeposits) of the project 
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# REST API using JSON for Messages
# Docker Compose Deployment
 
### To start all services:
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.api.yml up --build

The --build option is there for any code changes.

### Then Migrate and set up seed data:
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

### Distributed Tracing with Docker Compose
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)  

### Profiling
http://localhost:4000/debug/pprof/

### Metrics
http://localhost:4000/debug/vars

##### To view logs of running services in a separate terminal:
docker-compose -f ./deploy/compose/docker-compose.api.yml logs -f --tail 1  

### Shutdown 

docker-compose -f ./deploy/compose/docker-compose.api.yml down  
docker-compose -f ./deploy/compose/docker-compose.seed.yml down

#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

# Push Images to Docker Hub

docker build -t rsachdeva/illuminatingdeposits.api:v0.1 -f ./build/Dockerfile.api .  
docker push rsachdeva/illuminatingdeposits.api:v0.1 (as an example)  
docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed .  
docker push rsachdeva/illuminatingdeposits.seed:v0.1 (as an example)  

# Kubernetes Deployment - WIP

kubectl apply -f deploy/kubernetes/zipkin-deployment.yaml   
kubectl apply -f deploy/kubernetes/zipkin-service.yaml  

kubectl apply -f deploy/kubernetes/ic-traefik-lb.yaml  
kubectl apply -f deploy/kubernetes/ingress.yaml  

Access Traefik Dashboard at [http://localhost:3000/dashboard/#/](http://localhost:3000/dashboard/#/)   

### Distributed Tracing with Kubernetes Ingress

Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io/zipkin/](http://zipkin.127.0.0.1.nip.io/zipkin)  

### Shutdown

kubectl delete -f deploy/kubernetes/.

# HTTP Client Requests:
See resteditorclient/HealthCRUD.http for request examples and sample response.
Use dev env for localhost or change for prod if running web service at different IP address

# TLS files
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl . && \
docker run -v $PWD/config/tls:/tls -ti tlscert:v0.1