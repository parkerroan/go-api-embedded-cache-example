package webserver

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/storage"
)

const port = "8080"

type WebServer struct {
	StorageClient *storage.StorageClient
	Router        *mux.Router
}

func NewWebServer(storage *storage.StorageClient) *WebServer {
	r := mux.NewRouter()
	http.Handle("/", r)

	server := &WebServer{
		StorageClient: storage,
		Router:        r,
	}
	server.Register()
	return server
}

// Register registers the routes for the webserver
func (s *WebServer) Register() {
	s.Router.HandleFunc("/items/{id}", s.getItemsHandler).Methods("GET")
	s.Router.HandleFunc("/items", s.upsertItemsHandler).Methods("PUT")
}

// Run starts the webserver
func (s *WebServer) Run() {
	slog.Debug("Server started on port:", port)
	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		slog.Error(err.Error())
	}
}

func (s *WebServer) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	item, err := s.StorageClient.GetItem(id)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if item == nil {
		slog.Debug("Item not found:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	slog.Debug("Item:", item)

	//return item as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (s *WebServer) upsertItemsHandler(w http.ResponseWriter, r *http.Request) {
	item, err := s.StorageClient.GetItem("1")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("Item:", item)

	//return item as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
