package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-playground/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-playground/ $GOPATH/src/github.com/vasili-v/grpc-playground/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-playground/stream && ls $GOPATH/src/github.com/vasili-v/grpc-playground/stream"

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-playground/stream"
)

func main() {
	log.Printf("connecting to %s...", server)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c, err := grpc.DialContext(ctx, server,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
	)
	if err != nil {
		log.Fatalf("couldn't dial to %s: %s", server, err)
	}
	defer c.Close()

	log.Printf("connected to %s", server)

	client := pb.NewStreamClient(c)
	err = communicate(client)
	if err != nil {
		log.Printf("failed to communicate with %s: %s", server, err)

		log.Printf("second try")
		err = communicate(client)
		if err != nil {
			log.Fatalf("failed to communicate with %s: %s", server, err)
		}
	}
}

func communicate(c pb.StreamClient) error {
	s, err := c.New(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't create new stream: %s", err)
	}

	defer func() {
		log.Printf("closing stream")
		err := s.CloseSend()
		if err != nil {
			log.Printf("couldn't close stream")
		} else {
			log.Printf("waiting for stream status")
			s.Recv()
		}
	}()

	log.Printf("made new stream")

	m := &pb.Message{
		Payload: []byte{0x00, 0x01, 0x02, 0x03},
	}
	for i := 1; i < 5; i++ {
		log.Printf("sending message %d...", i)
		err = s.Send(m)
		if err != nil {
			return fmt.Errorf("couldn't send message %d: %s", i, err)
		}

		log.Printf("waiting for response %d...", i)
		r, err := s.Recv()
		if err != nil {
			return fmt.Errorf("couldn't receive response %d: %s", i, err)
		}

		log.Printf("checking response %d...", i)
		for j, b := range r.Payload {
			if m.Payload[j] != b {
				return fmt.Errorf("expected 0x%02x at %d in response %d but got 0x%02x", m.Payload[j], j, i, b)
			}
		}

		log.Print("cooldown for 1s...")
		time.Sleep(time.Second)
	}

	return nil
}
