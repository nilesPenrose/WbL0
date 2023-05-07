FROM golang:1.18 AS go

# Update together with .gitlab-ci.yml
FROM golangci/golangci-lint:v1.51.2 AS linter

FROM go AS dev
ENV INSIDE_DEV_CONTAINER 1
COPY --from=linter     /usr/bin/golangci-lint        /usr/bin/
WORKDIR /app
