/*
 license x
*/

package api

import (
	"net/http"
)

// HandlerFuncAPI interface for handlers for resources.
type HandlerFuncAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}
