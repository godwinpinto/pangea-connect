# Connect for Pangea
A set of middleware plugins / libraries build on pangea to add pangea features to the most popular API frameworks and API gateways

### Development Pre-requisite
- Docker with Docker compose

### How to run
1. Visit [Pangea.cloud](https://pangea.cloud) and get Pangea IP Intel token and pangea domain
2. Setup the environment (.env) file in the root directory with the values from step 1
3. use docker compose with folder name to run the respective solution

```
docker compose up <folder_name>
```

### Test:

```sh
// Test for a service which is malicous. Fire 2-3 times before getting forbidden 
curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.251' \
--header 'Content-Type: application/json' \
'

// Test for a service which is regular. Fire 2-3 times to verify
curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.252' \
--header 'Content-Type: application/json' \
'

```


### Roadmap Status
Below is the roadmap for Connect's collection with Pangea services

### Golang Framework Status
|  Framework | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  gin gonic | &check;  | &cross;  | &cross;  | &cross;  |
|  fiber    | &check;  | &cross;  | &cross;  | &cross;  |
|  echo | &cross;  | &cross;  | &cross;  | &cross;  |
|  chi | &cross;  | &cross;  | &cross;  | &cross;  |

### Java Framework Status
|  Framework | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Spring boot | &check;  | &cross;  | &cross;  | &cross;  |
|  Spring webflux    | &check;  | &cross;  | &cross;  | &cross;  |
|  Quarkus | &cross;  | &cross;  | &cross;  | &cross;  |
|  Quarkus Reactive | &check;  | &cross;  | &cross;  | &cross;  |
|  Vertx | &cross;  | &cross;  | &cross;  | &cross;  |
|  Micronaut | &cross;  | &cross;  | &cross;  | &cross;  |


### NodeJS Framework Status
|  Framework | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Express.js | &cross;  | &cross;  | &cross;  | &cross;  |
|  NestJS    | &cross;  | &cross;  | &cross;  | &cross;  |
|  Fastify | &cross;  | &cross;  | &cross;  | &cross;  |
|  Feather.js | &cross;  | &cross;  | &cross;  | &cross;  |
|  Meteor.js | &cross;  | &cross;  | &cross;  | &cross;  |
|  Sails | &cross;  | &cross;  | &cross;  | &cross;  |


### API Gateway Status
|  Platform | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Kong | &check;  | &cross;  | &cross;  | &cross;  |
|  Krakend    | &check;  | &cross;  | &cross;  | &cross;  |
|  Spring Cloud | &check;  | &cross;  | &cross;  | &cross;  |
|  Tyk | &cross;  | &cross;  | &cross;  | &cross;  |
|  Traefik | &cross;  | &cross;  | &cross;  | &cross;  |
|  Gravitee | &cross;  | &cross;  | &cross;  | &cross;  |
|  Apache APISix | &cross;  | &cross;  | &cross;  | &cross;  |
|  Apigee | &cross;  | &cross;  | &cross;  | &cross;  |


### Load Balancer
|  Platform | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Nginx | &cross;  | &cross;  | &cross;  | &cross;  |

### Ingress Gateway
|  Platform | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Nginx | &cross;  | &cross;  | &cross;  | &cross;  |
