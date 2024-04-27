```
./gradlew bootRun


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