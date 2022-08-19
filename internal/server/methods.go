package server

import (
	"fmt"
	"net/http"
)

type Muxer struct {
	app Application
	mux *http.ServeMux
}

var successResponse = []byte(`{"success": true}`)

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
	switch r.Method {
	case http.MethodPost:
		m.addBanner(w, r)
	case http.MethodDelete:
		m.removeBanner(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *Muxer) addBanner(w http.ResponseWriter, r *http.Request) {
	request, err := parseBannerRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = m.app.AddBanner(request.BannerID, request.SlotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(successResponse)
}

func (m *Muxer) removeBanner(w http.ResponseWriter, r *http.Request) {
	request, err := parseBannerRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	err = m.app.DeleteBanner(request.BannerID, request.SlotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(successResponse)
}

func (m *Muxer) click(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	request, err := parseClickRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}
	err = m.app.RegisterClick(request.BannerID, request.SlotID, request.GroupID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(successResponse)
}

func (m *Muxer) banner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	request, err := parseShowRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	bannerID, err := m.app.GetBanner(request.SlotID, request.GroupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"banner_id":%d}`, bannerID)))
}
