/*
 license x
*/

package api

import (
	"net/http"
)

type HandlerFuncAPI interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}
