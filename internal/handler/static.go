package handler

import "net/http"

type StaticHandler struct {
	fileServer http.Handler
}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{
		fileServer: http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	}
}

func (h *StaticHandler) StaticFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=86400")
	h.fileServer.ServeHTTP(w, r)
}

func (h *StaticHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
