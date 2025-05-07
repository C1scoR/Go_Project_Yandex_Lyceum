package config

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Litedb struct {
	DB *sql.DB
}

func ConfigDB() *Litedb {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	DB, err := sql.Open("sqlite3", "./internal/database/app.db")
	if err != nil {
		log.Fatal(err)
	}
	// Ограничение по одновременным подключениям к базе данных
	DB.SetMaxOpenConns(25)
	// Оставляем 5 соединений, которые будут всегда готовы принять запрос
	// (Нужно, чтобы заново не создавать tcp-соединения для запроса в БД)
	DB.SetMaxIdleConns(5)
	// Обновляем переодично соединение с БД
	DB.SetConnMaxLifetime(5 * time.Minute)
	if err = DB.PingContext(ctx); err != nil {
		log.Fatalf("database/config/config.goПодключение к базе данных не было выполнено успешно: %q", err)
	}
	sqlStatement1 := `
	CREATE TABLE IF NOT EXISTS users (
  		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		email TEXT NOT NULL,
  		password_hash TEXT NOT NULL);`
	sqlStatement2 := `CREATE TABLE IF NOT EXISTS statements (
		user_id INTEGER,
		statement_id TEXT PRIMARY KEY,
		statement TEXT,
		result TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);
		`
	_, err = DB.ExecContext(ctx, sqlStatement1)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s", err, sqlStatement1)
	}

	_, err = DB.ExecContext(ctx, sqlStatement2)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s", err, sqlStatement2)
	}
	return &Litedb{DB: DB}
}
