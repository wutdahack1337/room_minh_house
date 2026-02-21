package server

import (
	"context"
	"fmt"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

type Service struct {
	types.UnimplementedQueryServer
	dbService *database.MongoService
}

func NewService(dbService *database.MongoService) *Service {
	return &Service{dbService: dbService}
}

func (s *Service) CheckHealth(ctx context.Context, request *types.CheckHealthRequest) (*types.CheckHealthResponse, error) {
	if err := s.dbService.CheckHealth(ctx); err != nil {
		return &types.CheckHealthResponse{Status: fmt.Sprintf("database error: %v", err)}, err
	}

	return &types.CheckHealthResponse{Status: "i'm fine"}, nil
}
