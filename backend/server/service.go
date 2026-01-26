package server

import (
	"context"

	querytypes "github.com/wutdahack1337/introductory_capstone/backend/types/proto"
	v1 "github.com/wutdahack1337/introductory_capstone/backend/types/proto"
)

type Service struct {
	v1.UnimplementedQueryServer
}

func (s *Service) Health(ctx context.Context, request *querytypes.HealthRequest) (*querytypes.HealthResponse, error) {
	return &querytypes.HealthResponse{
		Status: "ok",
	}, nil
}
