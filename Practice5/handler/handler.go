package handler

import (
	"Practice5/models"
	"Practice5/repository"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	userRepo *repository.UserRepository
	db       *sql.DB
}

func NewHandler(userRepo *repository.UserRepository, db *sql.DB) *Handler {
	return &Handler{userRepo: userRepo, db: db}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.GetUsers)
	mux.HandleFunc("/users/common-friends", h.GetCommonFriends)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	filter := models.UserFilter{}

	filter.Page = parseIntDefault(q.Get("page"), 1)
	filter.PageSize = parseIntDefault(q.Get("page_size"), 10)
	filter.OrderBy = q.Get("order_by")
	filter.OrderDir = q.Get("order_dir")

	if v := q.Get("id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, "id должен быть числом", http.StatusBadRequest)
			return
		}
		filter.ID = &id
	}
	if v := q.Get("name"); v != "" {
		filter.Name = &v
	}
	if v := q.Get("email"); v != "" {
		filter.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		lower := strings.ToLower(v)
		filter.Gender = &lower
	}
	if v := q.Get("birth_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			writeError(w, "birth_date должен быть в формате YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		filter.BirthDate = &t
	}

	result, err := h.userRepo.GetPaginatedUsers(filter)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, result, http.StatusOK)
}

func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	user1Str := q.Get("user1")
	user2Str := q.Get("user2")

	if user1Str == "" || user2Str == "" {
		writeError(w, "Укажи параметры user1 и user2", http.StatusBadRequest)
		return
	}

	user1, err := strconv.Atoi(user1Str)
	if err != nil {
		writeError(w, "user1 должен быть числом", http.StatusBadRequest)
		return
	}
	user2, err := strconv.Atoi(user2Str)
	if err != nil {
		writeError(w, "user2 должен быть числом", http.StatusBadRequest)
		return
	}

	friends, err := repository.GetCommonFriends(h.db, user1, user2)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]interface{}{
		"user1_id":       user1,
		"user2_id":       user2,
		"common_friends": friends,
		"count":          len(friends),
	}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return def
	}
	return v
}
