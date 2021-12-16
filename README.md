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
