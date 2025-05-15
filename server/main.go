package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"server/handlers"
	"server/middlewares"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB

func main() {
	connStr := "postgres://user:password@db:5432/booking?sslmode=disable"
	var err error

	time.Sleep(5 * time.Second)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	// Настройка миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Ошибка создания драйвера миграции: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к миграциям внутри контейнера
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalf("Ошибка при создании миграции: %v\n", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграции: %v\n", err)
	} else if err == migrate.ErrNoChange {
		log.Println("Миграции уже применены. Нет изменений.")
	} else {
		log.Println("Миграции успешно применены!")
	}

	// Настроим маршруты
	http.Handle("/api/login", handlers.LoginHandler(db))
	http.Handle("/api/register", handlers.RegisterHandler(db))
	http.Handle("/api/booking", middlewares.CheckAuth(handlers.BookParkingSpot(db)))
	http.Handle("/api/bookings", middlewares.CheckAuth(handlers.GetOccupiedSpots(db)))

	// Создаем и настраиваем CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                             // Разрешаем все источники, замените на нужные для вашего приложения
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Разрешаем методы
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Разрешаем заголовки
		AllowCredentials: true,                                      // Разрешаем передачу cookies, если нужно
	})

	// Оборачиваем маршруты в CORS middleware
	handler := corsMiddleware.Handler(http.DefaultServeMux)

	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
