FROM golang:1.22-alpine AS builder 

WORKDIR /app 
COPY . .

WORKDIR /app/middleware 
RUN go mod download 

WORKDIR /app/example 
RUN go mod download 

RUN go build -o main main.go 

FROM alpine:latest AS runtime 

WORKDIR /app 
COPY --from=builder ./app/example/main . 
EXPOSE 8080 

CMD ["chmod", "+x", "main"] 
CMD ["./main"]