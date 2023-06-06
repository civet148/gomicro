# gomicro

golang grpc+go-micro.v2 wrapper

# install protoc and golang plugins binaries 

  https://github.com/civet148/protoc-plugins

# install protoc from source code

- download proto compiler source code

```shell script
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protobuf-cpp-3.19.1.tar.gz
$ tar xvfz protobuf-cpp-3.19.0.tar.gz
$ cd protobuf-cpp-3.19.0
$ ./configure && sudo make && sudo make install
```

- install protoc-gen-go

```shell script
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

- install protoc-gen-micro v2
```shell script
$ go get -u github.com/micro/protoc-gen-micro/v2
```


# A&Q

## compile problem

- grpc imports

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

```go
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
```

- go-micro imports

```go
go: github.com/micro/go-micro/v2@v2.9.1 requires
        github.com/micro/cli/v2@v2.1.2: reading https://goproxy.io/github.com/micro/cli/v2/@v/v2.1.2.mod: 404 Not Found
        server response:
        not found: github.com/micro/cli/v2@v2.1.2: invalid version: git ls-remote -q origin in /data1/golang/pkg/mod/cache/vcs/2f5431eb5439e9d79f82a6d853348656f17b78125db9eda81300bc014d0f0a5d: exit status 128:                fatal: could not read Username for 'https://github.com': terminal prompts disabled
        Confirm the import path was entered correctly.
        If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
make: *** [Makefile:4: gomicro] Error 1
```

replace in go.mod file
```go
replace github.com/micro/cli/v2 => github.com/civet148/micro-cli/v2 v2.1.2
```
