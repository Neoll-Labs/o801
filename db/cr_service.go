package db

import (
	"context"
)

type CRService[T any] interface {
	Create(context.Context, string) (T, error)
	Get(context.Context, int) (T, error)
}
