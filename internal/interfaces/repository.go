/*
 license x
*/

package interfaces

import (
	"context"
)

// Repository interface create and get.
type Repository[T any] interface {
	Create(context.Context, T) (T, error)
	Get(context.Context, T) (T, error)
}
