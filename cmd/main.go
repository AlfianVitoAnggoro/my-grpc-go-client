package main

import (
	"context"
	"log"

	"github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/adapter/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Logging Writer
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	var opts []grpc.DialOption

	// Disable TLS gRPC
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Connect to grpc server
	conn, err := grpc.NewClient("localhost:9090", opts...)

	if err != nil {
		log.Fatalln("Can't connect to grpc server", err)
	}

	// Close Connection after finish
	defer conn.Close()

	// Create HelloAdapter
	helloAdapter, err := hello.NewHelloAdapter(conn)

	if err != nil {
		log.Fatalln("Can't create HelloAdapter", err)
	}

	// Call SayHello
	// runSayHello(helloAdapter, "Alfian")

	// Call SayManyHellos
	runSayManyHellos(helloAdapter, "Vito")
}

func runSayHello(adapter *hello.HelloAdapter, name string) {
	// Call SayHello
	greet, err := adapter.SayHello(context.Background(), name)

	if err != nil {
		log.Fatalln("Can't call sayHello", err)
	}

	log.Println(greet.Greet)

}

func runSayManyHellos(adapter *hello.HelloAdapter, name string) {
	adapter.SayManyHellos(context.Background(), name)
}
