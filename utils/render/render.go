package render

import (
	"net/http"

	"github.com/go-chi/render"
)

type M map[string]interface{}


func JSON(w http.ResponseWriter, r *http.Request, statusCode int, body M) {
	// code, ok := body["code"].(int)

	render.JSON(w, r, body)
}

func PNG(w http.ResponseWriter, r *http.Request, statusCode int, v []byte) {
	render.Status(r, statusCode)
	w.Header().Set("Content-Type", "image/png")
	w.Write(v)
}

func WAV(w http.ResponseWriter, r *http.Request, statusCode int, v []byte) {
	render.Status(r, statusCode)
	w.Header().Set("Content-Type", "audio/wav")
	w.Write(v)
}
