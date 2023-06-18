package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "grpc/stream" // 실제 경로에 맞게 수정해주세요.
)

func main() {
	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// gRPC 클라이언트 생성
	client := pb.NewStreamClient(conn)

	// 클라이언트에서 스트림 수신
	stream, err := client.StartStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// 서버로부터 메시지 수신
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			// 서버가 스트림을 닫으면 종료
			break
		}
		if err != nil {
			log.Fatalf("Error receiving message from server: %v", err)
		}

		fmt.Printf("Received message from server: %s\n", message.GetContent())
	}
}
