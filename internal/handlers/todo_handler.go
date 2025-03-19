package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/internal/models"
	"todo-api/internal/repositories"
	"todo-api/internal/utils"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(float64)
	todo.UserID = uint(userID)

	db := r.Context().Value("db").(*gorm.DB)
	repo := repositories.NewTodoRepository(db)
	if err := repo.CreateTodo(&todo); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	userID := r.Context().Value("user_id").(float64)
	db := r.Context().Value("db").(*gorm.DB)
	repo := repositories.NewTodoRepository(db)
	todos, err := repo.GetTodos(uint(userID), page, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}

	response := map[string]interface{}{
		"data":  todos,
		"page":  page,
		"limit": limit,
		"total": len(todos),
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	userID := r.Context().Value("user_id").(float64)
	db := r.Context().Value("db").(*gorm.DB)
	repo := repositories.NewTodoRepository(db)
	todo, err := repo.GetTodoByID(uint(userID), uint(id))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var updatedTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(float64)
	db := r.Context().Value("db").(*gorm.DB)
	repo := repositories.NewTodoRepository(db)
	if err := repo.UpdateTodo(&models.Todo{
		ID:          uint(id),
		Title:       updatedTodo.Title,
		Description: updatedTodo.Description,
		UserID:      uint(userID),
	}); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedTodo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	db := r.Context().Value("db").(*gorm.DB)
	repo := repositories.NewTodoRepository(db)
	if err := repo.DeleteTodo(uint(id)); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
