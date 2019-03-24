# build stage
FROM golang:alpine AS build
ADD . /src
RUN go get github.com/magefile/mage && cd /src && mage vendor && mage check && mage kache

  # final stage
FROM alpine
WORKDIR /kache-home
COPY --from=build /src/build/kache-* /kache-home/kache
COPY --from=build /src/config /kache-home
EXPOSE 7088
ENTRYPOINT ./kache