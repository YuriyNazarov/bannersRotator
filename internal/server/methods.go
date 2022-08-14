package server

import "net/http"

type Muxer struct {
	app Application
	mux *http.ServeMux
}

func NewMux(app Application) *Muxer {
	muxer := Muxer{
		app: app,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/banners", muxer.banners)
	mux.HandleFunc("/click", muxer.click)
	mux.HandleFunc("/banner", muxer.banner)

	muxer.mux = mux
	return &muxer
}

func (m *Muxer) banners(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (m *Muxer) click(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (m *Muxer) banner(w http.ResponseWriter, r *http.Request) {
	//todo
}
