package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc/calculater"
)

type calculatorServer struct {
	pb.UnimplementedCalculatorServer
}

func (s *calculatorServer) ComputeAverage(stream pb.Calculator_ComputeAverageServer) error {
	sum := 0
	count := 0

	for {
		// 클라이언트로부터 숫자 스트림을 받음
		num, err := stream.Recv()
		if err != nil {
			if err.Error() == "io.EOF" {
				// 스트림 종료 시 평균 계산 후 응답 전송
				average := float32(sum) / float32(count)
				response := &pb.AverageResponse{Average: average}
				return stream.SendAndClose(response)
			}
			log.Fatalf("Error receiving number: %v", err)
		}

		// 받은 숫자로 평균 계산을 위해 누적 값과 개수 업데이트
		sum += int(num.GetValue())
		count++
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

	// gRPC 서비스 등록
	pb.RegisterCalculatorServer(grpcServer, &calculatorServer{})

	// 서버 시작
	fmt.Println("Server started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
