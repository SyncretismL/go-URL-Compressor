package server

import (
	"compressor/internal/compressor_grpc"
	"compressor/internal/urlData"
	"compressor/pkg/logger"
	"context"
)

type Server struct {
	logger   logger.Logger
	urlDatas urlData.URLDatas
}

func NewServer(newLogger logger.Logger, urlDatas urlData.URLDatas) *Server {
	return &Server{
		logger:   newLogger,
		urlDatas: urlDatas,
	}
}

func (s *Server) Create(c context.Context, request *compressor_grpc.CompressedURLRequest) (response *compressor_grpc.CompressedURLResponse, err error) {
	url := &urlData.URLData{
		URL: request.FullURL,
	}
	if err := s.urlDatas.SetURL(url); err != nil {
		s.logger.Debugf("%+s", err)
		return nil, err
	}
	url.URLCompressing()
	s.urlDatas.SetURLCompressed(url)

	output := url.URLCompressed
	response = &compressor_grpc.CompressedURLResponse{
		CompressedURL: output,
	}

	return response, nil

}

func (s *Server) Get(c context.Context, request *compressor_grpc.FullURLRequest) (response *compressor_grpc.FullURLResponse, err error) {
	url := &urlData.URLData{
		URLCompressed: request.CompressedURL,
	}
	if err := s.urlDatas.GetFullURL(url); err != nil {
		s.logger.Debugf("%+s", err)
		return nil, err
	}

	output := url.URL
	response = &compressor_grpc.FullURLResponse{
		FullURL: output,
	}
	return

}
