package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/mlc-project-template/internal/firebase"
)

const PORT string = "3000"

func main() {

	firebase.InitFirebase(context.Background())

	router := gin.Default()
	router.Use(cors.Default())

	router.Run(":" + PORT)
}
