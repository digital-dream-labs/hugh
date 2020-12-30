package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/status"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) RegisterEchoServer(srvr *grpc.Server, srv echo.EchoServer) *Server {
	return &Server{}
}

// UnaryEcho is unary echo.
func (s *Server) UnaryEcho(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: req.Message}, nil
}

// ServerStreamingEcho is server side streaming.
func (s *Server) ServerStreamingEcho(req *echo.EchoRequest, str echo.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

// ClientStreamingEcho is client side streaming.
func (s *Server) ClientStreamingEcho(req echo.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.NotFound, "unimplemented")
}

// BidirectionalStreamingEcho is bidi streaming.
func (s *Server) BidirectionalStreamingEcho(req echo.Echo_BidirectionalStreamingEchoServer) error {
	return status.Error(codes.NotFound, "unimplemented")
}
