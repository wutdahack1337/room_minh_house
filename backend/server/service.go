package server

import (
	"context"
	"fmt"
	"regexp"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

// MM-YYYY
var periodPattern = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

func validatePeriod(period string) error {
	if !periodPattern.MatchString(period) {
		return fmt.Errorf("period must be in MM-YYYY format (e.g., \"08-2025\")")
	}
	return nil
}

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
