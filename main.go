package main

import (
	"finalCode/database"
	"finalCode/handlers"
	"finalCode/repositories"
	"finalCode/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// เปิดการเชื่อมต่อกับฐานข้อมูล
	db := database.Postgresql()
	defer db.Close()

	// ตรวจสอบการเชื่อมต่อกับฐานข้อมูล
	err := db.Ping()
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	// สร้าง instances ของ repositories, services, และ handlers
	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerAdapter(s)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"}
	router.Use(cors.New(config))

	router.POST("/api/requestDoc", h.AddActivityHandlers)
	router.GET("/api/requestFile/:scdact_id", h.GETuserFileHandlers)

	err = router.Run(":8062")
	if err != nil {
		log.Fatal(err.Error())
	}
}
