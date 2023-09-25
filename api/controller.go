/*
 license x
*/

package api

import (
	"context"
)

// Repository interface create and get.
type Repository[T any] interface {
	Create(context.Context, T) (T, error)
	Fetch(context.Context, T) (T, error)
}
