# Reference: https://docs.docker.com/language/golang/build-images/
# https://github.com/GoogleContainerTools/skaffold/blob/bcbdfe043c2f334f919fa2e6ae06aed4a7578486/examples/helm-deployment/Dockerfile

## Build
# FROM golang:1.19-alpine3.17 as builder
FROM golang:1.19-alpine3.17 as builder

WORKDIR /code
# to add gcc back
RUN apk add build-base

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath -o /app

## Deploy - maybe want a distroless version
FROM alpine:3.17
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
EXPOSE 3000
COPY --from=builder /app .
CMD [ "./app" ]
