package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	pb "grpc/calculater"
)

func main() {
	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// gRPC 클라이언트 생성
	client := pb.NewCalculatorClient(conn)

	// 클라이언트 스트리밍을 위한 스트림 생성
	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// 서버로 숫자 스트림을 전송
	for i := 0; i < 5; i++ {
		num := rand.Int31n(10)
		fmt.Printf("Sending number: %v\n", num)
		stream.Send(&pb.Number{Value: num})
		time.Sleep(1 * time.Second)
	}

	// 스트림 종료 및 응답 수신
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Printf("Average: %v\n", response.GetAverage())
}
