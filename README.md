# Illuminating Deposits - Rest http json
###### All commands should be executed from the root directory (illuminatingdeposits-rest) of the project 
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# REST API using JSON for Messages

## Docker Compose Deployment

### Start postgres and tracing
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up 
```

### Then Migrate and set up seed data:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build
```

### To start all services without TLS:
Make sure DEPOSITS_REST_SERVICE_TLS=false in docker-compose.rest.server.yml
### To start all services with TLS:
Make sure DEPOSITS_REST_SERVICE_TLS=true in docker-compose.rest.server.yml
### And then execute:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.rest.server.yml up --build
``` 

The --build option is there for any code changes.

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

### Logs of running services (in a separate terminal):
docker-compose -f ./deploy/compose/docker-compose.rest.server.yml logs -f --tail 1  

### Distributed Tracing
Access [zipkin](https://zipkin.io/) service at [http://localhost:9411/zipkin/](http://localhost:9411/zipkin/)  

### Profiling
[http://localhost:4000/debug/pprof/](http://localhost:4000/debug/pprof/)

### Metrics
[http://localhost:4000/debug/vars](http://localhost:4000/debug/vars)  

### Shutdown 
```shell
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml down
docker-compose -f ./deploy/compose/docker-compose.rest.server.yml down
```

### Quick calculations with Same JSON output without actually invoking REST Http Method
Run at terminal:

```shell
docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate 
```

### To start only external db and trace service for working with local machine Editor/IDE:
Start postgres and tracing as usual
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up

### Then Migrate and set up seed data:
```shell
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.seed.yml up --build
````

Then Set the following env variables when starting directly running server: change as needed
And per your Editor/IDE:
```shell
export DEPOSITS_REST_SERVICE_TLS=true
export DEPOSITS_DB_DISABLE_TLS=true
export DEPOSITS_DB_HOST=127.0.0.1
export DEPOSITS_TRACE_URL=http://127.0.0.1:9411/api/v2/spans
go run ./cmd/server
```

### REST HTTP Services Endpoints Invoked:

#### Sanity test Client:
See    
cmd/sanitytestclient/main.go
The server side DEPOSITS_REST_SERVICE_TLS should be consistent and set for client also.
Uncomment any desired function request
Make sure to make email unique to avoid error.

### TLS files
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run -v $PWD/conf/tls:/tls tlscert:v0.1
``` 

To see openssl version being used in Docker:
```shell
docker build -t tlscert:v0.1 -f ./build/Dockerfile.openssl ./conf/tls && \
docker run -ti -v $PWD/conf/tls:/tls tlscert:v0.1 sh
```

You get a prompt
/tls

And enter version check
```shell
openssl version
```

### Troubleshooting
If for any reason no connection is happening from client to server or client hangs or server start up issues:
Run
```
ps aux | grep "go run"
ps aux | grep "go_build" 
```

to confirm is something else is already running


## Kubernetes Deployment - WIP

### Push Images to Docker Hub

```shell
docker build -t rsachdeva/illuminatingdeposits.rest.server:v0.1 -f ./build/Dockerfile.rest.server .  
docker push rsachdeva/illuminatingdeposits.rest.server:v0.1 
docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed .  
docker push rsachdeva/illuminatingdeposits.seed:v0.1  
``` 

### kubectl apply

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

# Version
v1.52