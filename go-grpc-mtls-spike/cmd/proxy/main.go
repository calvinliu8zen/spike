package main

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/calvin/grpc_spike/internal/protobuf"

	"github.com/calvin/grpc_spike/internal/utils"
)

const (
	hostAddress = "localhost:50001"
	network     = "tcp"

	serverAddress = "localhost:50002"

	clientCrtPath        = "../../data/out/localhost_50001.crt"
	clientPrivateKeyPath = "../../data/out/localhost_50001.key"
	hostCrtPath          = "../../data/out/localhost.crt"
	hostPrivateKeyPath   = "../../data/out/localhost.key"
	caCertPath           = "../../data/out//calvin.zendesk.com.crt"

	blank = ""
)

func proxy(ctx context.Context, req *pb.ProxyRequest) (*pb.ProxyResponse, error) {
	message := req.GetMessage()
	logger.Infof("[Proxy Service] received: %v", message)

	reply, err := proxyForward(message)
	if err != nil {
		return nil, err
	}

	logger.Infof("[Proxy Service] greeter repiled: %s ", reply)
	logger.Infof("[Proxy Service] repiled: %s ", reply)
	return &pb.ProxyResponse{Message: reply}, nil
}

func proxyForward(message string) (string, error) {
	certificate, certPool, err := utils.LoadTLSCreds(clientCrtPath, clientPrivateKeyPath, caCertPath)
	if err != nil {
		return blank, err
	}

	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   serverAddress,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return blank, err
	}

	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: message})
	if err != nil {
		return blank, err
	}

	return r.GetMessage(), nil
}

func main() {
	logger.Infof("[Proxy Service] starts at %s", hostAddress)
	listener, err := net.Listen(network, hostAddress)
	if err != nil {
		logger.Fatalf("[Proxy Service] failed to listen: %v", err)
	}

	certificate, _, err := utils.LoadTLSCreds(hostCrtPath, hostPrivateKeyPath, caCertPath)
	if err != nil {
		logger.Fatalf("[Proxy Service] failed to load certs: %v", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
	})

	server := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterProxyService(server, &pb.ProxyService{Forward: proxy})

	err = server.Serve(listener)
	if err != nil {
		logger.Fatalf("[Proxy Service] failed to start: %v", err)
	}
}
