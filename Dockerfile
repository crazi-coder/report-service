# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /report-gs-endpoint -buildvcs=false

##
## Deploy
##
FROM gcr.io/distroless/static-debian11
ENV GIN_MODE=release
WORKDIR /

COPY --from=build  /report-gs-endpoint /report-service

EXPOSE 3001

# USER nonroot:nonroot

ENTRYPOINT ["/report-service", "run"]