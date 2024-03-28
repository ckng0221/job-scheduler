# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS builder
ARG service

WORKDIR /app

COPY apps/${service}/go.mod apps/${service}/go.sum ./
RUN go mod download

COPY apps/${service}/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /${service}


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS production
ARG service
ENV ENV=production

WORKDIR /

# Get shell
COPY --from=builder /bin/sh /bin/sh

COPY --from=builder /${service} /go-app

EXPOSE 8000

# USER nonroot:nonroot

ENTRYPOINT ["/go-app"]

# Eg.
# docker build --build-arg="service=api" -f ./apps/Dockerfile.service . -t jobscheduler-api
# docker run -p 8000:8000 --env-file ./apps/api/.env jobscheduler-api

# docker build --build-arg="service=worker" -f ./apps/Dockerfile.service . -t jobscheduler-worker
# docker run --env-file ./apps/worker/.env jobscheduler-worker
