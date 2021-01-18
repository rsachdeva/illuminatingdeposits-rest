# Illuminating Deposits - Rest http json

###### All commands should be executed from the root directory (illuminatingdeposits-rest) of the project 
(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>

# REST API using JSON for Messages
# Features include:
- Golang (Go)  REST Http Service requests with json for Messages
- TLS for all requests
- Integration and Unit tests run in parallel
- Coverage Result for key packages
- Postgres DB health check service
- User Management service with Postgres for user creation
- JWT generation for Authentication
- JWT Authentication for Interest Calculations
- 30daysInterest for a deposit is called Delta
- Delta is for
    - each deposit
    - each bank with all deposits
    - all banks!
- Sanity test client included
- Docker support
- Docker compose deployment for development

# Docker Compose Deployment

### Start postgres and tracing services
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

# Kubernetes Deployment Manually -WIP
(for Better control; For Local Setup tested with Docker Desktop latest version with Kubernetes Enabled)

### Make docker images and Push Images to Docker Hub

```shell
docker build -t rsachdeva/illuminatingdeposits.rest.server:v1.3.60 -f ./build/Dockerfile.rest.server .  
docker build -t rsachdeva/illuminatingdeposits.seed:v1.3.60 -f ./build/Dockerfile.seed .  

docker push rsachdeva/illuminatingdeposits.rest.server:v1.3.60
docker push rsachdeva/illuminatingdeposits.seed:v1.3.60
``` 

### Start postgres service

```shell
kubectl apply -f deploy/kubernetes/postgres-env.yaml 
kubectl apply -f deploy/kubernetes/postgres.yaml
```

### Then Migrate and set up seed data manually for more control initially:
First should see in logs
database system is ready to accept connections
```kubectl logs pod/postgres-deposits-0```
And then execute migration/seed data for manual control when getting started:
```shell
kubectl apply -f deploy/kubernetes/seed.yaml
```
And if status for ```kubectl get pod``` 
shows completed for seed pod, optionally can be deleted:
```shell
kubectl delete -f deploy/kubernetes/seed.yaml
```
To connect external tool with postgres to see database internals use:
Use a connection string similar to:
jdbc:postgresql://127.0.0.1:30007/postgres
If still an issue you can try
kubectl port-forward service/postgres 5432:postgres
Now can easily connect using
jdbc:postgresql://localhost:5432/postgres


### Installing Ingress controller
Using helm to install nginx ingress controller
```shell
brew install helm
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
```
and then use 
```helm install ingress-nginx ingress-nginx/ingress-nginx```
to install ingress controller

```
### Start tracing service
```shell
kubectl apply -f deploy/kubernetes/zipkin.yaml
```
##### Distributed Tracing with Kubernetes Ingress
Access [zipkin](https://zipkin.io/) service at [http://zipkin.127.0.0.1.nip.io](http://zipkin.127.0.0.1.nip.io)
Sort Newest First and Click Find Traces

### Illuminating deposists Rest server in Kubernetes!
```shell
kubectl apply -f deploy/kubernetes/rest-server.yaml
```
And see logs using 
```kubectl logs -l app=restserversvc -f```

### Remove all resources / Shutdown

```shell
kubectl delete -f ./deploy/kubernetes/.
helm uninstall ingress-nginx
```

# Sanity test Client -REST HTTP Services Endpoints Invoked Externally:
Use env as
```export DEPOSITS_REST_SERVICE_ADDRESS=restserversvc.127.0.0.1.nip.io```
for Kubernetes Ingress 
otherwise use
```export DEPOSITS_REST_SERVICE_ADDRESS=localhost:3000```
Similarly,
The server side DEPOSITS_REST_SERVICE_TLS should be consistent and set for client also.
```export DEPOSITS_REST_SERVICE_TLS=false```

Example:
```shell
export GODEBUG=x509ignoreCN=0
export DEPOSITS_REST_SERVICE_TLS=false
export DEPOSITS_REST_SERVICE_ADDRESS=restserversvc.127.0.0.1.nip.io
go run ./cmd/sanitytestclient
```

With this Sanity test client, you will be able to:
- get status of Prostres DB
- add a new user
- JWT generation for Authentication
- JWT Authentication for Interest Delta Calculations for each deposit; each bank with all deposits and all banks
Quickly confirms Sanity check for the Envirinment set up with Kubernetes/Docker. 
There are also separate Integration and Unit tests.
  
# Running Integration/Unit tests
Tests are designed to run in parallel with its own test server and docker based postgres db using dockertest.
To run all tests with coverages reports for focussed packages:
Run following only once as tests use this image; so faster:
```shell 
docker pull postgres:11.1-alpine
``` 
And then run the following with coverages for key packages concerned:
```shell
go test -v -count=1 -covermode=count -coverpkg=./userauthn/...,./usermgmt/...,./postgreshealth/...,./interestcal/... -coverprofile cover.out ./... && go tool cover -func cover.out
go test -v -count=1 -covermode=count -coverpkg=./userauthn/...,./usermgmt/...,./postgreshealth/...,./interestcal/... -coverprofile cover.out ./... && go tool cover -html cover.out
```
Coverage Result for key packages:  
**total:	(statements)	96.3%**  

To run a single test - no coverage:
```shell 
go test -v -count=1 -run=TestServiceServer_CreateUser ./usermgmt/...
```
To run a single test - with coverage:
```shell 
go test -v -count=1 -covermode=count -coverpkg=./usermgmt -coverprofile cover.out -run=TestServiceServer_CreateUser ./usermgmt/... && go tool cover -func cover.out
```
The -v is for Verbose output: log all tests as they are run. Search "FAIL:" in parallel test output here to see reason for failure
in case any test fails.
Just to run all easily with verbose ouput:
```shell
go test -v ./... 
```
The -count=1 is mainly to not use caching and can be added as follows if needed for
any go test command:
```shell 
go test -v -count=1 ./...
```
See Editor specifcs to see Covered Parts in the Editor.
#### Test Docker containers for Postgresdb
Docker containers are mostly auto removed. This is done by passing true to testserver.InitRestServer(ctx, t, false)
in your test.
If you want to examine postgresdb data for a particular test, you can temporarily
set allowPurge as false in testserver.InitRestHttpServer(ctx, t, false) for your test.
Then after running specific failed test connect to postgres db in the docker container using any db ui.
As an example, if you want coverage on a specific package and run a single test in a package with verbose output:
```shell 
go test -v -count=1 -covermode=count -coverpkg=./usermgmt -coverprofile cover.out -run=TestServiceServer_CreateUser ./usermgmt/... && go tool cover -func cover.out
```
Any docker containers still running after tests should be manually removed:
```shell 
docker ps
docker stop $(docker ps -qa)
docker rm -f $(docker ps -qa)
```
And if mongodb not connecting for tests: (reference: https://www.xspdf.com/help/52284027.html)
```shell 
docker volume rm $(docker volume ls -qf dangling=true)
```

# TLS files
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

Check version using command:
```shell
openssl version
```

# Editor/IDE development without docker/docker compose/kubernetes as described above
To start only external db and trace service for working with local machine:  
Start postgres and tracing as usual
export COMPOSE_IGNORE_ORPHANS=True && \
docker-compose -f ./deploy/compose/docker-compose.external-db-trace-only.yml up

##### Then Migrate and set up seed data:
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

# Troubleshooting
If for any reason no connection is happening from client to server or client hangs or server start up issues:
Run
```
ps aux | grep "go run"
ps aux | grep "go_build" 
```

to confirm is something else is already running

# Version
v1.3.60