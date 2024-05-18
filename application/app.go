package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	router http.Handler
	db     *sql.DB
}

func New() *App {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	app := &App{
		db: db,
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// Выполнение запросов создания таблиц
	_, err := a.db.Exec(`
        CREATE TABLE IF NOT EXISTS ` + "`users`" + ` (
            user_id INTEGER PRIMARY KEY,
            image TEXT,
            email TEXT,
            password TEXT,
            nickname TEXT);
        
        CREATE TABLE IF NOT EXISTS ` + "`recipes`" + ` (
            recipe_id INTEGER PRIMARY KEY,
            user_id INTEGER,
            name TEXT,
            image TEXT,
            description TEXT,
            products TEXT,
            category TEXT,
            tag TEXT,
			video TEXT);

			
         CREATE TABLE IF NOT EXISTS ` + "`step`" + ` (
            step_number INTEGER PRIMARY KEY,
            step TEXT);

         CREATE TABLE IF NOT EXISTS ` + "`product`" + ` (
            product_id INTEGER PRIMARY KEY,
            name TEXT,
            weight INTEGER);
		
			CREATE TABLE IF NOT EXISTS ` + "`comments`" + ` (
			number INTEGER PRIMARY KEY,
			id_recipe INTEGER,
			image_user TEXT,
			name_user TEXT,
			comment TEXT,
			date DATETIME,
			image TEXT);
	
    `)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Проверка подключения к базе данных
	err = a.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Закрытие базы данных при завершении работы
	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("failed to close ", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
