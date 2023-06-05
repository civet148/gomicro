module github.com/civet148/gomicro/v2

go 1.16

require (
	github.com/civet148/log v1.1.3
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro/v2 v2.9.1
	google.golang.org/protobuf v1.26.0
)

replace github.com/micro/go-micro/v2 => ./third-party/micro/v2/go-micro

replace github.com/micro/cli/v2 => ./third-party/micro/v2/cli
