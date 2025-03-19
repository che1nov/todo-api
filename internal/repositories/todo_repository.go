package repositories

import (
	"errors"
	"fmt"
	"todo-api/internal/models"

	"gorm.io/gorm"
)

type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{DB: db}
}

func (r *TodoRepository) CreateTodo(todo *models.Todo) error {
	return r.DB.Create(todo).Error
}

func (r *TodoRepository) GetTodos(userID uint, page int, limit int) ([]models.Todo, error) {
	var todos []models.Todo
	offset := (page - 1) * limit
	if err := r.DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepository) GetTodoByID(userID uint, id uint) (*models.Todo, error) {
	var todo models.Todo
	if err := r.DB.Where("user_id = ? AND id = ?", userID, id).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo with ID %d not found", id)
		}
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) UpdateTodo(todo *models.Todo) error {
	return r.DB.Model(todo).Updates(todo).Error
}

func (r *TodoRepository) DeleteTodo(id uint) error {
	var todo models.Todo
	if err := r.DB.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("todo with ID %d not found", id)
		}
		return err
	}
	return r.DB.Delete(&todo).Error
}
