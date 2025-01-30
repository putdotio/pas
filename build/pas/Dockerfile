FROM golang:1.23.1

WORKDIR /pas
COPY go.* ./
RUN go mod download

COPY . ./
RUN go install

COPY ./build/pas/config.toml /etc/pas.toml

ENTRYPOINT ["pas", "-config", "/etc/pas.toml"]
