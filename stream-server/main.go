package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc/stream"
)

type streamServer struct {
	pb.UnimplementedStreamServer
}

func (s *streamServer) StartStreaming(stream pb.Stream_StartStreamingServer) error {
	// 클라이언트로 메시지 스트림을 보냄
	go func() {
		messages := []string{"Hello", "How are you?", "I'm fine, thank you!"}

		for _, msg := range messages {
			response := &pb.StreamMessage{
				Content: msg,
			}

			if err := stream.Send(response); err != nil {
				log.Fatalf("Error sending message to client: %v", err)
			}
		}
	}()

	// 클라이언트에서 메시지 수신
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			// 클라이언트가 스트림을 닫으면 종료
			break
		}
		if err != nil {
			log.Fatalf("Error receiving message from client: %v", err)
		}

		fmt.Printf("Received message from client: %s\n", message.GetContent())
	}

	return nil
}

func main() {
	// gRPC 서버 리스너 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()

	// Stream 서비스 등록
	pb.RegisterStreamServer(grpcServer, &streamServer{})

	// 서버 시작
	fmt.Println("Server started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
