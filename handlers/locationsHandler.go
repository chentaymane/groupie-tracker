package zone

import (
	"net/http"
)

func HandlerLocations(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	
}
