# Build container
FROM golang:1.21 AS build

ARG VERSION

COPY . /home/src
WORKDIR /home/src

RUN CGO_ENABLED=0 go build -v -o substrate-publisher -ldflags="-w -s"

# Runtime container
FROM alpine as substrate-publisher

RUN apk add --update --no-cache \
    ca-certificates \
    bash

WORKDIR /home/app
COPY --from=build /home/src/substrate-publisher .
RUN chmod +x /home/app/substrate-publisher

CMD ["./substrate-publisher", "start"]
