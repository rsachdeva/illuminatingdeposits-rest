Illuminating Deposits - CLI and API

![Illuminating Deposits Project Logo](logo.png)

# To start all services:
#### docker-compose -f docker-compose.api.yml --build

The --build option is there for any code changes.

## Then Migrate and set up seed data:
#### export COMPOSE_IGNORE_ORPHANS=True
#### docker-compose -f docker-compose.seed.yml up --build

# To view logs of running services:
#### docker-compose logs -f --tail 1

# To run HTTP requests:

See cmd/deltacli/httpreq/HealthCRUD.http
Use dev env for localhost or change for prod if running web service at different IP address


#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

docker build -f ./build/Dockerfile.calculate -t illumcalculate  . && \
docker run illumcalculate

WIP