Illuminating Deposits - CLI and API

![Illuminating Deposits Project Logo](logo.png)

# Local Development
 
### To start all services:
#### docker-compose -f docker-compose.api.yml up --build

The --build option is there for any code changes.

### Then Migrate and set up seed data:
#### export COMPOSE_IGNORE_ORPHANS=True
#### docker-compose -f docker-compose.seed.yml up --build

COMPOSE_IGNORE_ORPHANS is there for 
docker compose [setting](https://docs.docker.com/compose/reference/envvars/#compose_ignore_orphans).

##### To view logs of running services in a separate terminal:
###### docker-compose -f docker-compose.api.yml logs -f --tail 1

##### To run HTTP requests:
See cmd/httpclient/editorsupport/HealthCRUD.http for examples.
Use dev env for localhost or change for prod if running web service at different IP address


### Shutdown 

#### docker-compose -f docker-compose.api.yml down
#### docker-compose -f docker-compose.seed.yml down

#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

WIP