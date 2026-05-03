package controllers

import (
	"net/http"

	"github.com/arifazola/nontoon/services"
	"github.com/gin-gonic/gin"
)

type ChunkController struct {
	VideoService *services.VideoService
}

func (ctrl *ChunkController) GetLatestUploadedChunk(c *gin.Context) {
	uploadId := c.Param("uploadId")
	latestUploadedChunk, err := ctrl.VideoService.GetLatestUploadedChunk(c, uploadId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"index": latestUploadedChunk.Index})
}