Create an env file

export PANGEA_INTEL_TOKEN=
export PANGEA_DOMAIN=
export PANGEA_IP_INTEL_TYPE=
export PANGEA_IP_INTEL_SCORE_THRESHOLD=



```
cd plugin

//Docker targets
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.6.2 go build -buildmode=plugin -o krakend-pangea-connect.so .

//On prem installations
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.6.2-linux-generic go build -buildmode=plugin -o krakend-pangea-connect.so .

```
