FROM golang:1.23 AS builder

WORKDIR /app
COPY . .


RUN go mod tidy \
    && go get github.com/99designs/gqlgen/codegen/config \
    && go get github.com/99designs/gqlgen/internal/imports \
    && go get github.com/99designs/gqlgen \
    && go install github.com/99designs/gqlgen@latest \
    && go run github.com/99designs/gqlgen generate \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


#RUN go build -o main .


FROM alpine:latest
WORKDIR /app


RUN apk add --no-cache go


#COPY --from=builder /app/main .
COPY . .

EXPOSE 8081 50051 8080

#CMD ["go", "run", "main.go"]
CMD go run main.go;go run internal/server_proto.go;go run internal/client_proto.go
