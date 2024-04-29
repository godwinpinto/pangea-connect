FROM gradle:8.7.0-jdk21 AS build

WORKDIR /app
COPY . .

RUN gradle clean build --no-daemon

FROM openjdk:21-jdk-slim AS run

RUN adduser --system --group app-user

COPY --from=build --chown=app-user:app-user /app/build/libs/*-SNAPSHOT.jar app.jar

EXPOSE 8080
USER app-user

CMD ["java", "-jar", "app.jar"]