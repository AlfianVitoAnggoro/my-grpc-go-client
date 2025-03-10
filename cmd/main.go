package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/adapter/bank"
	"github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/adapter/hello"
	"github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/adapter/resiliency"
	dbank "github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/application/domain/bank"
	dresl "github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/application/domain/resiliency"
	"github.com/AlfianVitoAnggoro/my-grpc-go-client/internal/interceptor"
	resl_proto "github.com/AlfianVitoAnggoro/my-grpc-proto/protogen/go/resiliency"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cbreaker *gobreaker.CircuitBreaker

func init() {
	mybreaker := gobreaker.Settings{
		Name: "course-circuit-breaker",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			log.Printf("Circuit breaker failure is %v, requests is %v, means failure ratio : %v\n",
				counts.TotalFailures, counts.Requests, failureRatio)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		Timeout:     4 * time.Second,
		MaxRequests: 3,
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("Circuit breaker %v changed state, from %v to %v\n\n", name, from, to)
		},
	}

	cbreaker = gobreaker.NewCircuitBreaker(mybreaker)
}

func main() {
	// Logging Writer
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	var opts []grpc.DialOption

	// Disable TLS gRPC
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// opts = append(opts,
	// 	grpc.WithUnaryInterceptor(
	// 		grpc_retry.UnaryClientInterceptor(
	// 			grpc_retry.WithCodes(codes.Unknown, codes.Internal),
	// 			grpc_retry.WithMax(4),
	// 			grpc_retry.WithBackoff(grpc_retry.BackoffExponential(2*time.Second)),
	// 		),
	// 	),
	// )
	// opts = append(opts,
	// 	grpc.WithStreamInterceptor(
	// 		grpc_retry.StreamClientInterceptor(
	// 			grpc_retry.WithCodes(codes.Unknown, codes.Internal),
	// 			grpc_retry.WithMax(4),
	// 			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(3*time.Second)),
	// 		),
	// 	),
	// )

	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			interceptor.LogUnaryClientInterceptor(),
			interceptor.BasicUnaryClientInterceptor(),
			interceptor.TimeoutUnaryClientInterceptor(5*time.Second),
		),
	)
	opts = append(opts,
		grpc.WithChainStreamInterceptor(
			interceptor.LogStreamClientInterceptor(),
			interceptor.BasicClientStreamInterceptor(),
			interceptor.TimeoutStreamClientInterceptor(15*time.Second),
		),
	)

	// Connect to grpc server
	conn, err := grpc.NewClient("localhost:9090", opts...)

	if err != nil {
		log.Fatalln("Can't connect to grpc server", err)
	}

	// Close Connection after finish
	defer conn.Close()

	// Create HelloAdapter
	// helloAdapter, err := hello.NewHelloAdapter(conn)

	// if err != nil {
	// 	log.Fatalln("Can't create HelloAdapter", err)
	// }

	// bankAdapter, err := bank.NewBankAdapter(conn)

	// if err != nil {
	// 	log.Fatalln("Can not create BankAdapter :", err)
	// }

	resiliencyAdapter, err := resiliency.NewResiliencyAdapter(conn)

	if err != nil {
		log.Fatalln("Can not create ResiliencyAdapter :", err)
	}

	// Hello Adapter
	// Call SayHello
	// runSayHello(helloAdapter, "Alfian")

	// Call SayManyHellos
	// runSayManyHellos(helloAdapter, "Vito")

	// dummy data for client stream and bidirectional stream
	// names := []string{
	// 	"Alfian",
	// 	"Vito",
	// 	"Anggoro",
	// 	"Amelia",
	// 	"Tria",
	// 	"Agustina",
	// }

	// Call SayHelloToEveryone
	// runSayHelloToEveryone(helloAdapter, names)

	// Call SayHelloContinuous
	// runSayHelloContinuous(helloAdapter, names)

	// Bank
	// runGetCurrentBalance(bankAdapter, "7835697001") // Success
	// runGetCurrentBalance(bankAdapter, "7835697001xxx") //Failed

	// runFetchExchangeRates(bankAdapter, "USD", "IDR") //Success
	// runFetchExchangeRates(bankAdapter, "USD", "GBP") //Failed

	// runSummarizeTransactions(bankAdapter, "7835697001", 10) //Success
	// runSummarizeTransactions(bankAdapter, "7835697001xxx", 10) //Failed

	// runTransferMultiple(bankAdapter, "7835697001", "7835697003", 10)     //Success
	// runTransferMultiple(bankAdapter, "7835697001", "7835697003zzzz", 10) //Failed
	// runTransferMultiple(bankAdapter, "7835697001zzzz", "7835697003", 10) //Failed
	// runTransferMultiple(bankAdapter, "7835697001", "7835697003", 200)    //Failed

	// Resiliency
	// runUnaryResiliencyWithTimeout(resiliencyAdapter, 2, 8, []uint32{dresl.OK}, 5*time.Second) // give delayed 0 - 3 for success and 2 - 8 for error or success and 5+ for error

	// runServerStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 15*time.Second)

	// runClientStreamingResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 60*time.Second)
	// runBiDirectionalResiliencyWithTimeout(resiliencyAdapter, 0, 3, []uint32{dresl.OK}, 10, 60*time.Second)
	// runUnaryResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN, dresl.OK})
	// runServerStreamingResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN, dresl.OK})
	// runClientStreamingResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN}, 10)
	// runBiDirectionalResiliency(resiliencyAdapter, 0, 3, []uint32{dresl.UNKNOWN}, 10)

	// for i := 0; i < 300; i++ {
	// 	runUnaryResiliencyWithCircuitBreaker(resiliencyAdapter, 0, 0, []uint32{dresl.UNKNOWN, dresl.OK})
	// 	time.Sleep(time.Second)
	// }

	// Metadata
	// runUnaryResiliencyWithMetadata(resiliencyAdapter, 6, 10, []uint32{dresl.OK})
	runServerStreamingResiliencyWithMetadata(resiliencyAdapter, 1, 3, []uint32{dresl.OK})
	// runClientStreamingResiliencyWithMetadata(resiliencyAdapter, 0, 1, []uint32{dresl.OK}, 10)
	// runBiDirectionalResiliencyWithMetadata(resiliencyAdapter, 0, 1, []uint32{dresl.OK}, 10)
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

func runSayHelloToEveryone(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloToEveryone(context.Background(), names)
}

func runSayHelloContinuous(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloContinuous(context.Background(), names)
}

// Bank
func runGetCurrentBalance(adapter *bank.BankAdapter, acct string) {
	bal, err := adapter.GetCurrentBalance(context.Background(), acct)

	if err != nil {
		log.Fatalln("Failed to call GetCurrentBalance :", err)
	}

	log.Println(bal)
}

func runFetchExchangeRates(adapter *bank.BankAdapter, fromCur string, toCur string) {
	adapter.FetchExchangeRates(context.Background(), fromCur, toCur)
}

func runSummarizeTransactions(adapter *bank.BankAdapter, acct string, numDummyTransactions int) {
	var tx []dbank.Transaction

	for i := 1; i <= numDummyTransactions; i++ {
		ttype := dbank.TransactionTypeIn

		if i%3 == 0 {
			ttype = dbank.TransactionTypeOut
		}

		t := dbank.Transaction{
			Amount:          float64(rand.Intn(500) + 10),
			TransactionType: ttype,
			Notes:           fmt.Sprintf("Dummy transaction %v", i),
		}

		tx = append(tx, t)
	}

	adapter.SummarizeTransactions(context.Background(), acct, tx)
}

func runTransferMultiple(adapter *bank.BankAdapter, fromAcct string, toAcct string,
	numDummyTransactions int) {
	var trf []dbank.TransferTransaction

	for i := 1; i <= numDummyTransactions; i++ {
		tr := dbank.TransferTransaction{
			FromAccountNumber: fromAcct,
			ToAccountNumber:   toAcct,
			Currency:          "USD",
			Amount:            float64(rand.Intn(200) + 5),
		}

		trf = append(trf, tr)
	}

	adapter.TransferMultiple(context.Background(), trf)
}

// Resiliency
// With Time Out
func runUnaryResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := adapter.UnaryResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatalln("Failed to call UnaryResiliency :", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.ServerStreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes)
}

func runClientStreamingResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.ClientStreamingResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliencyWithTimeout(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int, timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	adapter.BiDirectionalResiliency(ctx, minDelaySecond, maxDelaySecond, statusCodes, count)
}

func runUnaryResiliency(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	res, err := adapter.UnaryResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatalln("Failed to call UnaryResiliency :", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) {
	adapter.ServerStreamingResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
}

func runClientStreamingResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.ClientStreamingResiliency(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliency(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.BiDirectionalResiliency(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

// Circuit Breaker

func runUnaryResiliencyWithCircuitBreaker(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) {
	cbreakerRes, cbreakerErr := cbreaker.Execute(
		func() (interface{}, error) {
			return adapter.UnaryResiliency(context.Background(), minDelaySecond, maxDelaySecond, statusCodes)
		},
	)

	if cbreakerErr != nil {
		log.Println("Failed to call UnaryResiliency :", cbreakerErr)
	} else {
		log.Println(cbreakerRes.(*resl_proto.ResiliencyResponse).DummyString)
	}
}

func runUnaryResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	res, err := adapter.UnaryResiliencyWithMetadata(context.Background(),
		minDelaySecond, maxDelaySecond, statusCodes)

	if err != nil {
		log.Fatalln("Failed to call UnaryResiliencyWithMetadata :", err)
	}

	log.Println(res.DummyString)
}

func runServerStreamingResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32) {
	adapter.ServerStreamingResiliencyWithMetadata(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes)
}

func runClientStreamingResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.ClientStreamingResiliencyWithMetadata(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}

func runBiDirectionalResiliencyWithMetadata(adapter *resiliency.ResiliencyAdapter,
	minDelaySecond int32, maxDelaySecond int32, statusCodes []uint32,
	count int) {
	adapter.BiDirectionalResiliencyWithMetadata(context.Background(), minDelaySecond,
		maxDelaySecond, statusCodes, count)
}
