package server

import (
	"context"

	types "github.com/wutdahack1337/introductory_capstone/backend/pkg/pb"
)

type Service struct {
	types.UnimplementedQueryServer
}

func (s *Service) Health(ctx context.Context, request *types.HealthRequest) (*types.HealthResponse, error) {
	return &types.HealthResponse{
		Status: "ok",
	}, nil
}
