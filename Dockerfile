FROM golang:1.24.4-bookworm AS builder
# RUN apt update -y
# RUN apt install -y unzip
# RUN PROTOC_VERSION=$(curl -s "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
# RUN wget -qO protoc.zip https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-$PROTOC_VERSION-linux-x86_64.zip
# RUN unzip -q protoc.zip bin/protoc -d /usr/local/
# RUN chmod a+x /usr/local/bin/protoc

WORKDIR /app
COPY . .
# RUN /app/install_protoc.sh
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# RUN go build -o app .
# ENTRYPOINT [ "entrypoint.sh" ]

# FROM ubuntu:latest

# WORKDIR /app
# COPY --from=builder /app/app .

# CMD ["./app"]