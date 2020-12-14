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
### To start all services with TLS:
Make sure DEPOSITS_WEB_SERVICE_SERVER_TLS=true in docker-compose.api.yml
### And then execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.api.yml up --build
``` 

The --build option is there for any code changes.


### Then Migrate and set up seed data:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build
````

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

### Logs of running services (in a separate terminal):
docker-compose -f ./deploy/compose/docker-compose.api.yml logs -f --tail 1  

### Distributed Tracing
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)  

### Profiling
[http://localhost:4000/debug/pprof/](http://localhost:4000/debug/pprof/)

### Metrics
[http://localhost:4000/debug/vars](http://localhost:4000/debug/vars)  

### Shutdown 
```shell
docker-compose -f ./deploy/compose/docker-compose.api.yml down  
docker-compose -f ./deploy/compose/docker-compose.seed.yml down 
```

### Quick calculations with Same JSON output without actually invoking REST Http Method
Run at terminal:

```shell
docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate 
```

### To start only external db and trace service for working with Editor/IDE:
Execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up --build
```

Set the following env variables when starting directly running server: change as needed
And per your Editor/IDE:
```shell
export DEPOSITS_WEB_SERVICE_SERVER_TLS=true
export DEPOSITS_DB_DISABLE_TLS=true
export DEPOSITS_DB_HOST=127.0.0.1
export DEPOSITS_TRACE_URL=http://127.0.0.1:9411/api/v2/spans
go run github.com/rsachdeva/illuminatingdeposits-rest/tools/dbcli migrate  (only once)
go run github.com/rsachdeva/illuminatingdeposits-rest/tools/dbcli seed     (only once)
go run github.com/rsachdeva/illuminatingdeposits-rest/cmd/server
```

### Interest Service REST HTTP Methods Invoked:

#### Without TLS Sanity test Client:   
See tools/resteditor/HealthCRUD.http for request examples and sample response.
Use dev env for localhost:3000  
Or go to cmd/sanitytestclient/main.go  
And uncomment any desired function request starting with "withoutTls..."
Make sure to make email unique to avoid error
See Base64EncodedString call line in the file

#### TLS Sanity test Client:  
Go to cmd/sanitytestclient/main.go  
And uncomment any desired function request starting with "tls..."
run cmd/sanitytestclient
Make sure to make email unique to avoid error
See Base64EncodedString call line in the file

# Push Images to Docker Hub
```shell
docker build -t rsachdeva/illuminatingdeposits.api:v0.1 -f ./build/Dockerfile.api .  
docker push rsachdeva/illuminatingdeposits.api:v0.1 
docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed .  
docker push rsachdeva/illuminatingdeposits.seed:v0.1  
``` 

# Kubernetes Deployment - WIP

```shell
kubectl apply -f deploy/kubernetes/traefik-ingress-daemonset-service.yaml 

kubectl apply -f deploy/kubernetes/zipkin-deployment.yaml   
kubectl apply -f deploy/kubernetes/zipkin-service.yaml   
kubectl apply -f deploy/kubernetes/zipkin-ingress.yaml  

kubectl apply -f deploy/kubernetes/postgres-config.yaml 
kubectl apply -f deploy/kubernetes/postgres-stateful.yaml  
kubectl apply -f deploy/kubernetes/postgres-service.yaml  
``` 

To connect external tool with postgres to see database internals use:
Use a connection string similar to:
jdbc:postgresql://127.0.0.1:30007/postgres
If still an issue you can try
kubectl port-forward service/postgres 5432:postgres
Now can easily connect using
jdbc:postgresql://localhost:5432/postgres

Access Traefik Dashboard at [http://localhost:3000/dashboard/#/](http://localhost:3000/dashboard/#/)   

### Distributed Tracing with Kubernetes Ingress

Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io/zipkin](http://zipkin.127.0.0.1.nip.io/zipkin)

### Shutdown

kubectl delete -f deploy/kubernetes/.

# TLS files
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl . && \
docker run -v $PWD/config/tls:/tls tlscert:v0.1
``` 

To see openssl version being used in Docker:
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl . && \
docker run -ti -v $PWD/config/tls:/tls tlscert:v0.1 sh

/tls # openssl version
OpenSSL 1.1.1g  21 Apr 2020
```

# Version
v1.0