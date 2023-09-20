/*
 license x
*/

package api

import (
	"net/http"
)

type ServicesAPIRouter interface {
	GetOrCreateUser(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
}
