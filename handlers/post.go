package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/ptflp/go-chi-rest"

	"github.com/ptflp/go-chi-rest/database"

	"github.com/go-chi/chi"
	"github.com/ptflp/go-chi-rest/data"
	"github.com/ptflp/go-chi-rest/data/post"
)

// NewPostHandler ...
func NewPostHandler(db *database.DB) *PostsHandler {
	return &PostsHandler{
		service: post.NewPostService(db.Db),
	}
}

// PostsHandler ...
type PostsHandler struct {
	service data.PostRepo
}

// Fetch all post data
func (p *PostsHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, _ := p.service.Fetch(r.Context(), 5)

	respondwithJSON(w, http.StatusOK, payload)
}

// Create a new post
func (p *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	pst := models.Post{}
	json.NewDecoder(r.Body).Decode(&pst)

	newID, err := p.service.Create(r.Context(), &pst)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// Update a post by id
func (p *PostsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	pst := models.Post{ID: int64(id)}
	json.NewDecoder(r.Body).Decode(&pst)
	payload, err := p.service.Update(r.Context(), &pst)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a post details
func (p *PostsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.service.GetByID(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// Delete a post
func (p *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := p.service.Delete(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
