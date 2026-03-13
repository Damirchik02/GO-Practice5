package repository

import (
	"Practice5/models"
	"database/sql"
	"fmt"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var allowedOrderFields = map[string]string{
	"id":         "id",
	"name":       "name",
	"email":      "email",
	"gender":     "gender",
	"birth_date": "birth_date",
}

func (r *UserRepository) GetPaginatedUsers(filter models.UserFilter) (models.PaginatedResponse, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *filter.ID)
		argIdx++
	}
	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Name+"%")
		argIdx++
	}
	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+*filter.Email+"%")
		argIdx++
	}
	if filter.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, *filter.Gender)
		argIdx++
	}
	if filter.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", argIdx))
		args = append(args, *filter.BirthDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "id"
	if col, ok := allowedOrderFields[filter.OrderBy]; ok {
		orderBy = col
	}

	orderDir := "ASC"
	if strings.ToUpper(filter.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("ошибка подсчёта: %w", err)
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}
	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(
		`SELECT id, name, email, gender, birth_date FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIdx, argIdx+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	if users == nil {
		users = []models.User{}
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}, nil
}
