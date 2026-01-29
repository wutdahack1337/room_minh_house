package server

import (
	"context"
	"fmt"

	"github.com/wutdahack1337/introductory_capstone/backend/internal/database"
	types "github.com/wutdahack1337/introductory_capstone/backend/pkg/pb"
)

type Service struct {
	types.UnimplementedQueryServer
	dbService database.DatabaseService
}

func (s *Service) Health(ctx context.Context, request *types.HealthRequest) (*types.HealthResponse, error) {
	if err := s.dbService.Health(ctx); err != nil {
		return &types.HealthResponse{
			Status: fmt.Sprintf("database error: %v", err),
		}, err
	}

	return &types.HealthResponse{
		Status: "ok",
	}, nil
}
