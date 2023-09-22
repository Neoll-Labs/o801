/*
 license x
*/

package internal

import (
	"context"
)

// Repository interface create and get
type Repository[T any] interface {
	Create(context.Context, string) (T, error)
	Get(context.Context, int) (T, error)
}
