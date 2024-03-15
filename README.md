# Create OpenAPI from gRPC Gateway Service

This is an example gRPC service that has two services, one a gRPC service, and another that's exposed via a gRPC Gateway as REST.

The REST exposure definition is used to then define an OpenAPI Spec with a protoc tool called [gnostic](https://github.com/google/gnostic).


## Install all the tools

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
```

## Build the protos

Since this uses the http googleapis protos, clone that into a dir and copy out some files.

```
git clone git@github.com:googleapis/googleapis.git
mkdir -p bookstore/proto/google/api
cp googleapis/google/api/http.proto bookstore/proto/google/api
cp googleapis/google/api/annotations.proto bookstore/proto/google/api
```

Then, generate the code from protos and generate the OpenAPI file.

This should create the `bookstore/pb` protobuf generated go code in from the protobuf definition `proto` dir and also create an openapi.yaml file.

```
cd bookstore
go mod tidy

protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=. --openapi_out=.
```

## Test

```
go run *.go
2024/03/14 22:16:27 gRPC server started on port 8080
2024/03/14 22:16:27 gRPC-gateway started on port 8090
```

Then in another terminal

```
$ grpcurl -plaintext localhost:8080 bookstore.Inventory.GetBooks
{
  "books": [
    {
      "title": "Book 1",
      "author": "Author 1",
      "pages": 100
    },
    {
      "title": "Book 2",
      "author": "Author 2",
      "pages": 200
    }
  ]
}
$ curl localhost:8090/v1/echo -d '{"value":"hello you big boy"}'
{"value":"hello you big boy"}

```


## Thanks to

* https://github.com/google/gnostic “A compiler for APIs described by the OpenAPI Specification with plugins for code generation and other API support tasks.”
* https://github.com/google/gnostic/tree/main/cmd/protoc-gen-openapi “... protoc plugin that generates an OpenAPI description for a REST API that corresponds to a Protocol Buffer service.”
* https://grpc.io/docs/languages/go/basics/
* https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/
* https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/#using-protoc
