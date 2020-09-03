package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	pb "github.com/calvin/grpc_spike/internal/protobuf"
	"github.com/calvin/grpc_spike/internal/utils"

	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	hostAddress = "localhost:50002"
	network     = "tcp"

	hostCrtPath    = "../../data/out/localhost_50002.crt"
	privateKeyPath = "../../data/out/localhost_50002.key"
	caCertPath     = "../../data/out//calvin.zendesk.com.crt"
)

func sayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Infof("[Greeter Service] received: %v", req.GetName())
	message := fmt.Sprintf("Welcome %s", req.GetName())
	logger.Infof("[Greeter Service] replied: %v", message)
	return &pb.HelloReply{Message: message}, nil
}

func main() {
	logger.Infof("[Greeter Service] starts at %s", hostAddress)

	listener, err := net.Listen(network, hostAddress)
	if err != nil {
		logger.Fatalf("[Greeter Service] failed to listen: %v", err)
	}

	certificate, certPool, err := utils.LoadTLSCreds(hostCrtPath, privateKeyPath, caCertPath)
	if err != nil {
		logger.Fatalf("[Proxy Service] failed to load certs: %v", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	server := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGreeterService(server, &pb.GreeterService{SayHello: sayHello})

	err = server.Serve(listener)
	if err != nil {
		logger.Fatalf("[Greeter Service] failed to start: %v", err)
	}
}
