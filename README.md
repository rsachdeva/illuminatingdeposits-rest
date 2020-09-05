Illuminating Deposits - CLI and API

![Illuminating Deposits Project Logo](logo.png)

# Docker Compose Based Deployment
 
### To start all services:
#### docker-compose -f ./deploy/docker-compose.api.yml up --build

The --build option is there for any code changes.

### Then Migrate and set up seed data:
#### export COMPOSE_IGNORE_ORPHANS=True
#### docker-compose -f ./deploy/docker-compose.seed.yml up --build

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

##### To view logs of running services in a separate terminal:
###### docker-compose -f ./deploy/docker-compose.api.yml logs -f --tail 1

##### To run HTTP requests:
See cmd/httpclient/editorsupport/HealthCRUD.http for examples.
Use dev env for localhost or change for prod if running web service at different IP address


### Shutdown 

#### docker-compose -f ./deploy/docker-compose.api.yml down
#### docker-compose -f ./deploy/docker-compose.seed.yml down

#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

# Push Images to Docker Hub

docker build -t rsachdeva/illuminatingdeposits.api:v0.1 -f ./build/Dockerfile.api .  

docker push rsachdeva/illuminatingdeposits.api:v0.1 (as an example) 

docker build -t rsachdeva/illuminatingdeposits.seed:v0.1 -f ./build/Dockerfile.seed . 

docker push rsachdeva/illuminatingdeposits.seed:v0.1 (as an example) 


(Development is WIP)