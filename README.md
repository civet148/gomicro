# gomicro

golang grpc+go-micro.v2 wrapper

# install protoc and golang plugins binaries 

  https://github.com/civet148/protoc-plugins

# install protoc from source code

- 1. download proto compiler source code

```shell script
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protobuf-cpp-3.19.1.tar.gz
$ tar xvfz protobuf-cpp-3.19.0.tar.gz
$ cd protobuf-cpp-3.19.0
$ ./configure && sudo make && sudo make install
```

- 2. install protoc-gen-go

```shell script
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

- 3. install protoc-gen-micro v2
```shell script
$ go get -u github.com/micro/protoc-gen-micro/v2
```


- 4. A&Q

- compile problem

```shell
go: finding module for package google.golang.org/grpc/naming
go: finding module for package google.golang.org/grpc/examples/helloworld/helloworld
go: found google.golang.org/grpc/examples/helloworld/helloworld in google.golang.org/grpc/examples v0.0.0-20230516222055-92e65c890c9a
go: finding module for package google.golang.org/grpc/naming
node-agent/pkg/client imports
        github.com/civet148/gomicro/v2 imports
        github.com/micro/go-micro/v2/registry/etcd imports
        github.com/coreos/etcd/clientv3 tested by
        github.com/coreos/etcd/clientv3.test imports
        github.com/coreos/etcd/integration imports
        github.com/coreos/etcd/proxy/grpcproxy imports
        google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.55.0), but does not contain package google.golang.org/grpc/naming
make: *** [Makefile:10: agent] Error 1
```

replace in go.mod file

```golang
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
```
