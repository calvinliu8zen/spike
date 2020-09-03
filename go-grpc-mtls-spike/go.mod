module github.com/calvin/grpc_spike

go 1.15

require (
	github.com/golang/protobuf v1.4.2
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.33.0-dev.0.20200828165940-d8ef479ab79a
