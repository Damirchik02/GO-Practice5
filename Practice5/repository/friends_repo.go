package repository

import (
	"Practice5/models"
	"database/sql"
	"fmt"
)

func GetCommonFriends(db *sql.DB, userID1, userID2 int) ([]models.User, error) {
	if userID1 == userID2 {
		return nil, fmt.Errorf("нельзя сравнивать пользователя с самим собой")
	}

	query := `
        SELECT u.id, u.name, u.email, u.gender, u.birth_date
        FROM user_friends uf1
        JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
        JOIN users u ON u.id = uf1.friend_id
        WHERE uf1.user_id = $1
          AND uf2.user_id = $2
        ORDER BY u.id
    `

	rows, err := db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}

	if friends == nil {
		friends = []models.User{}
	}

	return friends, nil
}
