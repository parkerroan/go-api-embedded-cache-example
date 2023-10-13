package webserver

import (
	"context"
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
func (s *WebServer) Run(ctx context.Context) error {
	slog.Debug("Server started on port:", port)
	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		return err
	}
	return nil
}

func (s *WebServer) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	item, err := s.StorageClient.GetItem(ctx, id)
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
	ctx := r.Context()

	//create item from request body
	var item storage.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if id exists in db
	if item.ID != "" {
		_, err := s.StorageClient.GetItem(ctx, item.ID)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	//upsert item in db
	result, err := s.StorageClient.UpsertItem(ctx, item)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//return item as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
