package system

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	//TODO check db connectivity
	w.Write([]byte("up and running"))
}
