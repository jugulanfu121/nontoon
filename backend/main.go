package main

import (
	"fmt"

	"github.com/arifazola/nontoon/controllers"
	"github.com/arifazola/nontoon/repositories"
	"github.com/arifazola/nontoon/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
func main(){
	fmt.Println("Start");
	router := gin.Default()
	router.Use(cors.Default())

	localStorage := repositories.LocalStorage{
		BasePath: "./files",
	}

	videoService := services.VideoService {
		FileStorage: &localStorage,
	}

	videoController := controllers.VideoController{
		VideoService: &videoService,
	}
	
	router.GET("/videos", controllers.GetAllVideos)

	router.POST("/videos", videoController.UploadVideo)

	router.POST("/videos/chunks", videoController.UploadChunk)

	router.POST("/videos/merge", videoController.CompleteUpload)

	router.Run("localhost:8080")
}