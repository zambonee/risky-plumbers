package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	"arcticwolf.com/cutler/models"
	"github.com/gorilla/mux"
)

const (
	port = ":8080"
)

var (
	writeTimeout = time.Second * 5
	readTimeout  = time.Second * 5
	IdleTimeout  = time.Second * 60
)

type Server struct {
	Backend models.DAOInterface
}

// New instantiates a new Server with the given DAO parameter, returning the http.Server for the caller to serve as it sees fit.
func New(backend models.DAOInterface) *http.Server {
	s := Server{
		Backend: backend,
	}
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/risks", s.GetAllRisks).Methods("GET")
	v1.HandleFunc("/risks", s.SaveRisk).Methods(("POST"))
	v1.HandleFunc("/risks/{id}", s.GetRiskByID).Methods("GET")
	return &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  IdleTimeout,
	}
}

// GetAllRisks handles the GET /risks API calls.
func (s Server) GetAllRisks(w http.ResponseWriter, r *http.Request) {
	// Consider bounding or paginating the results
	w.Header().Set("Content-Type", "application/json")
	risks := s.Backend.GetAllRisks(r.Context())
	if err := json.NewEncoder(w).Encode(risks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// SaveRisk handles the POST /risks API calls.
func (s Server) SaveRisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var risk models.CreateRiskRequest
	if err := json.NewDecoder(r.Body).Decode(&risk); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch risk.State {
	case models.StateAccepted, models.StateClosed, models.StateInvestigating, models.StateOpen:
		break
	default:
		http.Error(w, "invalid risk state", http.StatusBadRequest)
		return
	}
	result, err := s.Backend.SaveRisk(r.Context(), risk.State, risk.Title, risk.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}

// GetRiskByID handles the GET /risks/{id} API calls.
func (s Server) GetRiskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := mux.Vars(r)
	id := v["id"]
	if id == "" {
		http.Error(w, "no risk id in URI path", http.StatusBadRequest)
		return
	}
	risk, err := s.Backend.GetRiskByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if risk == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(risk)
	w.WriteHeader(http.StatusOK)
}
