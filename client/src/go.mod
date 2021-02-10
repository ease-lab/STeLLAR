module vhive-bench/client

go 1.15

replace (
	vhive-bench/client/benchmarking => ./benchmarking
	vhive-bench/client/benchmarking/networking => ./benchmarking/networking
	vhive-bench/client/benchmarking/networking/benchgrpc/proto_gen => ./benchmarking/networking/benchgrpc/proto_gen
	vhive-bench/client/benchmarking/visualization => ./benchmarking/visualization
	vhive-bench/client/benchmarking/writers => ./benchmarking/writers
	vhive-bench/client/setup => ./setup
	vhive-bench/client/util => ./util
)

require (
	github.com/ajstarks/svgo v0.0.0-20200725142600-7a3c8b57fecb // indirect
	github.com/aws/aws-sdk-go v1.37.6
	github.com/go-gota/gota v0.10.1
	github.com/golang/protobuf v1.4.3
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	gonum.org/v1/gonum v0.8.1
	gonum.org/v1/plot v0.8.0
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea // indirect
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
