package db

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	createTables(db)
	seedUsers(db)
	seedFriends(db)
}

func createTables(db *sql.DB) {
	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id         SERIAL PRIMARY KEY,
        name       VARCHAR(100) NOT NULL,
        email      VARCHAR(100) UNIQUE NOT NULL,
        gender     VARCHAR(10)  NOT NULL,
        birth_date DATE         NOT NULL
    );`

	friendsTable := `
    CREATE TABLE IF NOT EXISTS user_friends (
        user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
        friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
        PRIMARY KEY (user_id, friend_id),
        CONSTRAINT no_self_friendship CHECK (user_id <> friend_id)
    );`

	if _, err := db.Exec(usersTable); err != nil {
		log.Fatalf("Ошибка создания таблицы users: %v", err)
	}
	if _, err := db.Exec(friendsTable); err != nil {
		log.Fatalf("Ошибка создания таблицы user_friends: %v", err)
	}
	log.Println("Таблицы созданы")
}

func seedUsers(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count >= 20 {
		log.Println("Пользователи уже добавлены")
		return
	}

	users := []struct{ name, email, gender, birth string }{
		{"Alice Johnson", "alice@mail.com", "female", "1995-03-15"},
		{"Bob Smith", "bob@mail.com", "male", "1990-07-22"},
		{"Carol White", "carol@mail.com", "female", "1998-11-05"},
		{"David Brown", "david@mail.com", "male", "1992-01-30"},
		{"Eva Martinez", "eva@mail.com", "female", "1997-06-18"},
		{"Frank Wilson", "frank@mail.com", "male", "1988-09-12"},
		{"Grace Lee", "grace@mail.com", "female", "2000-04-25"},
		{"Henry Taylor", "henry@mail.com", "male", "1993-12-08"},
		{"Iris Anderson", "iris@mail.com", "female", "1996-02-14"},
		{"Jack Thomas", "jack@mail.com", "male", "1991-08-03"},
		{"Kate Jackson", "kate@mail.com", "female", "1999-05-20"},
		{"Liam Harris", "liam@mail.com", "male", "1994-10-17"},
		{"Mia Garcia", "mia@mail.com", "female", "2001-07-09"},
		{"Noah Martinez", "noah@mail.com", "male", "1989-03-27"},
		{"Olivia Robinson", "olivia@mail.com", "female", "1997-11-11"},
		{"Peter Clark", "peter@mail.com", "male", "1986-06-04"},
		{"Quinn Lewis", "quinn@mail.com", "female", "2002-01-16"},
		{"Ryan Walker", "ryan@mail.com", "male", "1995-09-28"},
		{"Sophia Hall", "sophia@mail.com", "female", "1993-04-07"},
		{"Tom Allen", "tom@mail.com", "male", "1990-12-22"},
		{"Uma Young", "uma@mail.com", "female", "1998-08-14"},
		{"Victor Hernandez", "victor@mail.com", "male", "1987-05-31"},
	}

	query := `INSERT INTO users (name, email, gender, birth_date) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`
	for _, u := range users {
		if _, err := db.Exec(query, u.name, u.email, u.gender, u.birth); err != nil {
			log.Printf("Ошибка добавления %s: %v", u.name, err)
		}
	}
	log.Println("Добавлены 22 пользователя")
}

func seedFriends(db *sql.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM user_friends").Scan(&count)
	if count > 0 {
		log.Println("Дружеские связи уже добавлены")
		return
	}

	friendships := [][2]int{
		{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6},
		{2, 3}, {2, 4}, {2, 5}, {2, 7},
		{3, 4}, {3, 8},
		{4, 9}, {5, 10}, {5, 11},
		{6, 7}, {7, 8}, {8, 9}, {9, 10},
		{10, 11}, {11, 12}, {12, 13}, {13, 14},
		{14, 15}, {15, 16}, {16, 17}, {17, 18},
		{18, 19}, {19, 20},
	}

	query := `INSERT INTO user_friends (user_id, friend_id) VALUES ($1,$2),($2,$1) ON CONFLICT DO NOTHING`
	for _, f := range friendships {
		if _, err := db.Exec(query, f[0], f[1]); err != nil {
			log.Printf("Ошибка дружбы %d-%d: %v", f[0], f[1], err)
		}
	}
	log.Println("Дружеские связи добавлены")
}
