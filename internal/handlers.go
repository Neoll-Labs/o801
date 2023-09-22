/*
 license x
*/

package internal

import (
	"net/http"
)

type Handlers interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}
