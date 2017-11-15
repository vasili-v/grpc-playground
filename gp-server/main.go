package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-playground/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-playground/ $GOPATH/src/github.com/vasili-v/grpc-playground/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-playground/stream && ls $GOPATH/src/github.com/vasili-v/grpc-playground/stream"

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-playground/stream"
)

func handler(in *pb.Message) *pb.Message {
	return &pb.Message{
		Payload: in.Payload,
	}
}

type server struct{}

func (s *server) New(stream pb.Stream_NewServer) error {
	log.Printf("got new stream")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("coudn't receive message: %s", err)
			return err
		}

		err = stream.Send(handler(in))
		if err != nil {
			log.Printf("coudn't send message: %s", err)
			return err
		}
	}

	log.Printf("stream depleted")
	return nil
}

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("couldn't start tcp listener on %s: %s", address, err)
	}
	log.Printf("listening for tcp connections on %s", address)

	p := grpc.NewServer()
	pb.RegisterStreamServer(p, &server{})

	log.Printf("starting server on %s...", address)
	err = p.Serve(ln)
	if err != nil {
		log.Fatalf("couldn't start gRPC server on %s: %s", address, err)
	}
}
