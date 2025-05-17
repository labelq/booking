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
	"github.com/gorilla/mux"
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
		"file://migrations",
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

	// Создаем новый роутер
	router := mux.NewRouter()

	// Настраиваем маршруты
	router.HandleFunc("/api/login", handlers.LoginHandler(db)).Methods("POST")
	router.HandleFunc("/api/register", handlers.RegisterHandler(db)).Methods("POST")
	router.Handle("/api/booking", middlewares.CheckAuth(handlers.BookParkingSpot(db))).Methods("POST")
	router.Handle("/api/bookings", middlewares.CheckAuth(handlers.GetOccupiedSpots(db))).Methods("GET")

	// Административные маршруты
	router.HandleFunc("/api/admin/bookings", handlers.GetAllBookings(db)).Methods("GET")
	router.HandleFunc("/api/admin/bookings/{id}", handlers.CancelBooking(db)).Methods("DELETE")
	router.HandleFunc("/api/admin/blocked-spots", handlers.GetBlockedSpots(db)).Methods("GET")
	router.HandleFunc("/api/admin/spots/toggle-block", handlers.ToggleSpotBlock(db)).Methods("POST")
	router.HandleFunc("/api/admin/users", handlers.GetUsersHandler(db)).Methods("GET")
	router.HandleFunc("/api/admin/users/{id}/role", handlers.UpdateUserRoleHandler(db)).Methods("PUT")

	// Создаем и настраиваем CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:80", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Применяем CORS middleware к роутеру
	handler := corsMiddleware.Handler(router)

	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
