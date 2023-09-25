/*
 license x
*/

package interfaces

import (
	"github.com/nelsonstr/o801/internal/model"
)

// Cacheable interface for the types that can be stored at cache.
type Cacheable interface {
	model.UserView
}

type Cache[T Cacheable] interface {
	Set(key int64, value *T)
	Get(key int64) (T, bool)
	Len() int
}
