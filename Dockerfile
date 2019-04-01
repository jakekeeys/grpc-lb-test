FROM golang:alpine as build
RUN apk add --no-cache git
WORKDIR /go/src/github.com/jakekeeys/grpc-lb-test
COPY . .
RUN go get -d -v .
RUN go build -v -o app .

FROM alpine
WORKDIR /service
COPY --from=build /go/src/github.com/jakekeeys/grpc-lb-test/app .
COPY --from=build /go/src/github.com/jakekeeys/grpc-lb-test/swagger.json .
ENTRYPOINT ["./app"]