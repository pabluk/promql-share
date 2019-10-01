package application

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Share struct {
	ID    string `json:"id"`
	Query string `json:"query"`
}

func NewShareHandler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var share Share
		_ = json.NewDecoder(r.Body).Decode(&share)
		logger.Info("Create new share for query: ", share.Query)

		share.ID = "abc123"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&share)
	}
}

func GetShareHandler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var share Share
		params := mux.Vars(r)
		logger.Info("Getting details for shareID: ", params["shareID"])
		share.ID = params["shareID"]
		share.Query = "graph?key=value"
		json.NewEncoder(w).Encode(&share)
	}
}

func GoToShareHandler(logger *logrus.Logger, base_url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var share Share
		params := mux.Vars(r)
		logger.Info("Redirection for shareID: ", params["shareID"])
		share.ID = params["shareID"]
		share.Query = "graph?key=value"
		http.Redirect(w, r, base_url+share.Query, 301)
	}
}
