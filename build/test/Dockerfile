FROM golang:1.23.1 AS builder

# Install mysql
RUN apt-get update && \
    apt-get -y install --no-install-recommends default-mysql-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /pas

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./build/test/entrypoint.sh /root/entrypoint.sh
COPY ./build/test/config.toml /etc/pas.toml
COPY . .

ENTRYPOINT ["/bin/bash", "/root/entrypoint.sh"]
