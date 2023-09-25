/*
 license x
*/

package api

import (
	"context"
)

// ServiceAPI interface create and get.
type ServiceAPI[T any] interface {
	Create(context.Context, T) (T, error)
	Get(context.Context, T) (T, error)
}
