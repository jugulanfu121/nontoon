package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arifazola/nontoon/services"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	VideoService *services.VideoService
}

func GetAllVideos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, services.GetAllVideos())
}

func (videoController *VideoController) UploadVideo(c *gin.Context){
	fileHeader, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot open file"})
		return
	}

	log.Println("received file: ", fileHeader.Filename)

	uploadVideo, errUpload := videoController.VideoService.SaveVideo(file, fileHeader.Filename, fileHeader.Size)

	if errUpload != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed upload video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": uploadVideo})
}

func (c *VideoController) UploadChunk(ctx *gin.Context) {
    uploadID := ctx.PostForm("uploadId")
    chunkIndexStr := ctx.PostForm("chunkIndex")

    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(400, gin.H{"error": "file required"})
        return
    }

    f, _ := file.Open()
    defer f.Close()

    index, _ := strconv.Atoi(chunkIndexStr)

    err = c.VideoService.SaveChunk(uploadID, index, f, ctx)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(200, gin.H{"status": "chunk uploaded"})
}

func (c *VideoController) CompleteUpload(ctx *gin.Context) {
    uploadID := ctx.PostForm("uploadId")
    filename := ctx.PostForm("filename")
    totalChunksStr := ctx.PostForm("totalChunks")

    totalChunks, _ := strconv.Atoi(totalChunksStr)

    path, err := c.VideoService.MergeChunks(uploadID, filename, totalChunks)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(200, gin.H{
        "message": "upload complete",
        "path": path,
    })
}