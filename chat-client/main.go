package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "grpc/chat"
)

type chatServer struct{}

func (s *chatServer) StartStreaming(stream pb.Chat_StartStreamingServer) error {
	for {
		// 클라이언트로부터 메시지 스트림을 받음
		message, err := stream.Recv()
		if err == io.EOF {
			// 클라이언트가 스트림을 닫으면 종료
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("Received message from client: %s\n", message.GetContent())

		// 서버에서 응답을 생성하여 클라이언트에게 보냄
		response := &pb.ChatMessage{
			Content: "Server received message: " + message.GetContent(),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func main() {
	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// gRPC 클라이언트 생성
	client := pb.NewChatClient(conn)

	// 클라이언트 스트리밍을 위한 스트림 생성
	stream, err := client.StartStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// 서버로 메시지 스트림 전송
	go func() {
		messages := []string{"Hello", "How are you?", "I'm fine, thank you!"}

		for _, msg := range messages {
			request := &pb.ChatMessage{
				Content: msg,
				// 다른 필드를 추가하려면 여기에 추가하세요.
			}

			if err := stream.Send(request); err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			time.Sleep(1 * time.Second)
		}

		// 스트림 종료
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("Error closing stream: %v", err)
		}
	}()

	// 서버로부터 메시지 스트림 응답 수신
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// 서버가 스트림을 닫으면 종료
			break
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}

		fmt.Printf("Received response from server: %s\n", response.GetContent())
	}
}
