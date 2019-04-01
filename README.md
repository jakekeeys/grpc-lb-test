[![Docker Repository on Quay](https://quay.io/repository/jakekeeys/grpc-lb-test/status "Docker Repository on Quay")](https://quay.io/repository/jakekeeys/grpc-lb-test)

```$xslt
NAME:
   grpc-lb-test - A new cli application

USAGE:
   grpc-lb-test [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     server
     client
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Server
```$xslt
NAME:
   grpc-lb-test server -

USAGE:
   grpc-lb-test server [command options] [arguments...]

OPTIONS:
   --grpc-bind value  (default: ":8090") [$GRPC_BIND]
   --http-bind value  (default: ":8080") [$HTTP_BIND]
```

Client
```$xslt
NAME:
   grpc-lb-test client -

USAGE:
   grpc-lb-test client [command options] [arguments...]

OPTIONS:
   --grpc-server value  (default: ":8090") [$GRPC_SERVER]
   --http-bind value    (default: ":8081") [$HTTP_BIND]
   --interval-ms value  (default: 100) [$INTERVAL_MS]
```