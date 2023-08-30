package httpapi

import (
	"avito_task/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const port = "8080"

func NewHTTPHandler(manager service.Manager) *HTTPHandler {
	return &HTTPHandler{
		manager: manager,
	}
}

type HTTPHandler struct {
	manager service.Manager
}

type AddSeg struct {
	Segments []string `json:"segments"`
}

func (h *HTTPHandler) CreateSegments(rw http.ResponseWriter, r *http.Request) {
	var data AddSeg
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("ADD: ", data.Segments)
	err = h.manager.PostS(data.Segments)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}
	resp, _ := json.Marshal(string("OK!"))
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

type ResponseSegments struct {
	Segments pq.StringArray `json:"segments"`
}

func (h *HTTPHandler) GetSegments(rw http.ResponseWriter, r *http.Request) {
	personID, err := strconv.Atoi(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:])
	if err != nil {
		http.Error(rw, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	log.Println("GET SEGMENTS: ", personID)
	var seg ResponseSegments
	seg.Segments, err = h.manager.GetSegments(personID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}

	resp, _ := json.Marshal(seg)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

type ResponseIDs struct {
	IDs pq.Int64Array `json:"ids"`
}

func (h *HTTPHandler) GetIDs(rw http.ResponseWriter, r *http.Request) {
	segment := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

	log.Println("GET IDS: ", segment)
	var seg ResponseIDs
	var err error
	seg.IDs, err = h.manager.GetIDs(segment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}

	resp, _ := json.Marshal(seg)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

type AddRequest struct {
	PersonID int      `json:"personID"`
	Segments []string `json:"segments"`
}

func (h *HTTPHandler) AddPersonAndSegments(rw http.ResponseWriter, r *http.Request) {
	var data AddRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("ADD: ", data.PersonID, data.Segments)
	err = h.manager.PostPaS(data.PersonID, data.Segments)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}
	resp, _ := json.Marshal(string("OK!"))
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

type DeleteSeg struct {
	Segment string `json:"segment"`
}

func (h *HTTPHandler) DeleteSegment(rw http.ResponseWriter, r *http.Request) {
	var data DeleteSeg
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("DELETE: ", data.Segment)
	err = h.manager.DeleteSegment(data.Segment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}
	resp, _ := json.Marshal(string("OK!"))
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

type DeleteSegs struct {
	Segments []string `json:"segments"`
}

func (h *HTTPHandler) DeleteSegments(rw http.ResponseWriter, r *http.Request) {
	personID, err := strconv.Atoi(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:])
	if err != nil {
		http.Error(rw, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	var data DeleteSegs
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("DELETE: ", personID, data.Segments)
	err = h.manager.DeleteSegments(personID, data.Segments)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}

	resp, _ := json.Marshal(string("OK!"))
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(resp)
}

func NewServer(manager service.Manager) *http.Server {

	handler := NewHTTPHandler(manager)
	r := mux.NewRouter()
	r.HandleFunc("/api/seg", handler.CreateSegments).Methods(http.MethodPost)
	r.HandleFunc("/api/add", handler.AddPersonAndSegments).Methods(http.MethodPost)
	r.HandleFunc("/api/person/{personID:[a-z0-9]+}", handler.GetSegments).Methods(http.MethodGet)
	r.HandleFunc("/api/segment/{segment:[A-Za-z0-9_\\-]+}", handler.GetIDs).Methods(http.MethodGet)
	r.HandleFunc("/api/deleteSegment", handler.DeleteSegment).Methods(http.MethodPost)
	r.HandleFunc("/api/deleteSegments/{personID:[a-z0-9]+}", handler.DeleteSegments).Methods(http.MethodPost)

	return &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
