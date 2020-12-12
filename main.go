package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/KarasWinds/tag-service/proto"
	"github.com/KarasWinds/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcPort string
var httpPort string

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "gRPC啟動通訊埠")
	flag.StringVar(&httpPort, "http_port", "9001", "http啟動通訊埠")
	flag.Parse()
}

func main() {
	errs := make(chan error)
	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errs <- err
		}
	}()
	log.Printf("Server start")
	select {
	case err := <-errs:
		log.Fatalf("Run Server err : %v", err)
	}
}

func RunHttpServer(port string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(
		"/ping",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`pong`))
		},
	)

	return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer(port string) error {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return s.Serve(lis)
}
