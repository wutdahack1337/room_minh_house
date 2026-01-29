package database

import "context"

type DatabaseService interface {
	Health(ctx context.Context) error
}
