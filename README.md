Illuminating Deposits - CLI and API

![Illuminating Deposits Project Logo](logo.png)

# To start all services:
#### docker-compose up --build

The --build option for any code changes.

## Then Migrate and set up seed data:
#### docker-compose -f docker-compose.seed.yml up --build

# To view logs of running services:
#### docker-compose logs -f --tail 1

# To run HTTP requests:

See cmd/deltacli/httpreq/HealthCRUD.http
Use dev env for localhost or change for prod if running web service at different IP address


#### As a Side note to run quick calculations with JSON output without HTTP 
Run at terminal:

sh ./dockerquickcal.sh

WIP