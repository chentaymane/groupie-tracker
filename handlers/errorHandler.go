package zone

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

// HandleError renders an error page with the given status and message
func HandleError(w http.ResponseWriter, status int, message string) {
	tmpl, err := template.ParseFiles("templates/err.html")
	if err != nil {
		io.WriteString(w, "Failed to load error template")
		return
	}

	data := struct {
		Message string
		Status  int
	}{
		Message: message,
		Status:  status,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
