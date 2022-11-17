package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"pikpo2/config"
	"pikpo2/internal/constant"
	"pikpo2/internal/membership"
	"pikpo2/internal/tier"
	"pikpo2/internal/user"
)

// var DB *gorm.DB

func main() {
	cfg := config.New()

	// db, err := gorm.Open(mysql.Open(cfg.Database.DSN))
	// if err != nil {
	// 	fmt.Println("ERROR DATABASE CONNECTION")
	// }

	// db.AutoMigrate()

	// DB = db

	db, err := sql.Open("mysql", cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	tierRepo := tier.NewTierRepository(db, constant.TableTier)

	userRepo := user.NewUserRepository(db, constant.TableUser, constant.TableTier)
	userUseCase := user.NewUserUseCase(userRepo, tierRepo)

	membershipUsecase := membership.NewMembershipUseCase(userRepo, tierRepo)

	user.NewUserHandler(router, userUseCase)
	membership.NewMembershipHanlder(router, membershipUsecase)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	fmt.Println("SERVER ON")
	log.Fatal(server.ListenAndServe())
}
