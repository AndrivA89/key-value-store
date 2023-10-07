FROM golang:1.17 as build
COPY . /go/src/key-value-store
WORKDIR /go/src/key-value-store/cmd
RUN go mod download github.com/stretchr/testify@v1.8.1
RUN GOARCH=amd64 go build -o kvs

FROM alpine:3.14
COPY --from=build /go/src/key-value-store/cmd .
EXPOSE 8080
EXPOSE 50051
RUN chmod +x kvs
CMD ["./kvs"]
