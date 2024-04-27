Create an env file

export PANGEA_INTEL_TOKEN=
export PANGEA_DOMAIN=
export PANGEA_IP_INTEL_TYPE=
export PANGEA_IP_INTEL_SCORE_THRESHOLD=


```
curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.251' \
--header 'Content-Type: application/json' \
--data '{
    "name":"Godwinz"
}'

curl --location --request GET 'http://localhost:8080/get' \
--header 'X-Forwarded-For: 190.28.74.252' \
--header 'Content-Type: application/json' \
--data '{
    "name":"Godwinz"
}'


```

```
cd plugin

//Docker targets
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.6.2 go build -buildmode=plugin -o krakend-pangea-connect.so .

//On prem installations
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.6.2-linux-generic go build -buildmode=plugin -o krakend-pangea-connect.so .

```
