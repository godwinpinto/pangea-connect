# Connect for Pangea
A set of middleware plugins / libraries build on Pangea to add pangea security to the most popular API frameworks and API gateways thus acting as a drop-in centralized security solution. 

### Development Pre-requisite
- Docker with Docker compose

### How to run
1. Visit [Pangea.cloud](https://pangea.cloud) and get Pangea IP Intel token and pangea domain
2. Clone the repo
```sh
git clone https://github.com/godwinpinto/pangea-connect.git 
cd pangea-connect 
```
3. Setup the environment (rename example.env to .env) file in the root directory with the values from step 1
4. Use docker compose with folder name to run the respective solution

```
docker compose up redis httpbin <folder_name>
```

### Test:

```sh
# Test for a service which is malicous. Fire 2-3 times before getting forbidden 
curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.251' \
--header 'Content-Type: application/json'

# Test for a service which is regular. Fire 2-3 times to verify, should always result in success response
curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.252' \
--header 'Content-Type: application/json'

```
> Note on testing: Kong works on Port 8000



### Roadmap Status
Below is the roadmap for Connect's collection with Pangea services

> Note: An improved approach could also be to publish the plugins on central packaged registeries for Teams adopting as is usage.
> IP Intel service = IP Reputation

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
|  Express.js | &check;  | &cross;  | &cross;  | &cross;  |
|  NestJS    | &cross;  | &cross;  | &cross;  | &cross;  |
|  Fastify | &cross;  | &cross;  | &cross;  | &cross;  |
|  Feather.js | &cross;  | &cross;  | &cross;  | &cross;  |
|  Meteor.js | &cross;  | &cross;  | &cross;  | &cross;  |
|  Sails | &cross;  | &cross;  | &cross;  | &cross;  |


### Rust Framework Status
|  Framework | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Actix-web | &cross;  | &cross;  | &cross;  | &cross;  |
|  Axum    | &cross;  | &cross;  | &cross;  | &cross;  |
|  Gloo | &cross;  | &cross;  | &cross;  | &cross;  |

### Python Framework Status
|  Framework | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Django | &cross;  | &cross;  | &cross;  | &cross;  |
|  Flask    | &cross;  | &cross;  | &cross;  | &cross;  |
|  Fast API | &cross;  | &cross;  | &cross;  | &cross;  |


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

### Ingress Controller
|  Platform | IP Intel  | Secure Audit Log | Embargo | File Intel |
|---|---|---|---|---|
|  Nginx | &cross;  | &cross;  | &cross;  | &cross;  |

### Improvement needed in existing plugins
- Revisiting the solution by comparing other middleware / plugin (doing it right) for each platform / framework
- Unit Tests, Performance Test, etc
- Basic development hygience (code comments, code analysis)
- CI pipelines
- Extending Connect to other Pangea services that can be centralized and easy to configure
