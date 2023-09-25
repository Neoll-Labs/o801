/*
 license x
*/

package interfaces

import (
	"net/http"
)

// HandlerAPI interface for handlers for resources.
type HandlerAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}
