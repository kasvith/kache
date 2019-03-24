# build stage
FROM golang:alpine AS build-env
ADD .. /src
RUN go get github.com/magefile/mage && cd /src && mage vendor && mage check && mage kache

  # final stage
FROM alpine
WORKDIR /kache
COPY --from=build-env /src/build/kache-* /app/kache
ENTRYPOINT ./kache