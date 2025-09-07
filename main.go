package main

import (
	"Personal-expense-tracking-system/config"
	"Personal-expense-tracking-system/database"
	"Personal-expense-tracking-system/rep"
	"Personal-expense-tracking-system/service"
	"Personal-expense-tracking-system/transport"
	"log"
	"net/http"
)

func main() {
	//загрузка конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cant load conf. :", err)
	}

	db, err := database.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// --- Инициализация зависимостей для Пользователей ---
	userRepo := rep.NewUserRepo(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)

	// --- Инициализация зависимостей для Расходов ---
	expenseRepo := rep.NewExpenseRepo(db)
	expenseService := service.NewExpenseService(expenseRepo)

	// Теперь передаем оба сервиса в роутер
	router := transport.NewRouter(userService, expenseService, cfg.JWTSecret)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
