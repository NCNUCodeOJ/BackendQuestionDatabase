# build stage
FROM golang:1.16-alpine AS build-env
COPY . /src
RUN apk add build-base && cd /src && go build -o app

# final stage
FROM alpine:3
WORKDIR /app
COPY --from=build-env /src/app /app/app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
ENTRYPOINT ./app
EXPOSE 8080
HEALTHCHECK --timeout=5s CMD ./app ping