/*
 license x
*/

package db

import (
	"context"
)

// CRService interface create and get users
type CRService[T any] interface {
	Create(context.Context, string) (T, error)
	Get(context.Context, int) (T, error)
}
