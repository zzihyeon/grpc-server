package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc/chat"
)

type chatServer struct {
	pb.UnimplementedChatServer
}

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
	// gRPC 서버 리스너 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()

	// Chat 서비스 등록
	pb.RegisterChatServer(grpcServer, &chatServer{})

	// 서버 시작
	fmt.Println("Server started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
