package sample

import (
	"context"
)

type Service struct {
}

func (s *Service) SampleRPC(ctx context.Context, request *SampleRequest) (*SampleResponse, error) {
	return &SampleResponse{}, nil
}
