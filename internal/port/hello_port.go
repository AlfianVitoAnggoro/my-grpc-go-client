package port

import (
	"context"

	"github.com/AlfianVitoAnggoro/my-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
	SayManyHellos(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[hello.HelloResponse], error)
	SayHelloToEveryone(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[hello.HelloRequest, hello.HelloResponse], error)
	SayHelloContinuous(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[hello.HelloRequest, hello.HelloResponse], error)
}
