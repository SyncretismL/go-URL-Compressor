package server

import (
	"compressor/internal/compressor_grpc"
	"compressor/internal/urlData/mocks"
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	mockURLData := &mocks.URLDatas{}
	mockURLData.On("SetURL", mock.AnythingOfType("*urlData.URLData")).Return(nil)
	mockURLData.On("SetURLCompressed", mock.AnythingOfType("*urlData.URLData")).Return(nil)
	mockURLData.On("GetFullURL", mock.AnythingOfType("*urlData.URLData")).Return(nil)

	compressor_grpc.RegisterCompressingServiceServer(s, &Server{urlDatas: mockURLData})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
		url  string
		res  string
		err  error
	}{
		{
			"valid request with non negative amount",
			"http://hhh.com/data/phonenumber2",
			"http://hhh.com/data/aaaaaaaaaa",
			nil,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := compressor_grpc.NewCompressingServiceClient(conn).Create(context.Background(), &compressor_grpc.CompressedURLRequest{FullURL: tt.url})
			if err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			if response.CompressedURL != tt.res {
				t.Error("error: expected", tt.res, "received", response)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}

}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		url  string
		res  string
		err  error
	}{
		{
			"valid request with non negative amount",
			"http://hhh.com/data/aaaaaaaaaa",
			"http://hhh.com/data/phonenumber2",
			nil,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := compressor_grpc.NewCompressingServiceClient(conn).Get(context.Background(), &compressor_grpc.FullURLRequest{CompressedURL: tt.url})
			if err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			if response.FullURL != tt.res {
				t.Error("error: expected", tt.res, "received", response)
			}

			if err != nil && errors.Is(err, tt.err) {
				t.Error("error: expected", tt.err, "received", err)
			}
		})
	}

}
