package handlers

import (
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/IavilaGw/cat-api/internal/services"
)

type CatHandler struct {
	catService *services.CatService
}

func NewCatHandler(catService *services.CatService) *CatHandler {
	return &CatHandler{catService: catService}
}

func (h *CatHandler) GetRandomCat(c *gin.Context) {
	log.Println("GET /api/cat")

	catImage, imageData, err := h.catService.FetchAndSaveRandomCat()
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch cat image",
			"message": err.Error(),
		})
		return
	}

	c.Header("X-Image-ID", strconv.FormatUint(uint64(catImage.ID), 10))
	c.Header("X-Image-Hash", catImage.ImageHash)
	c.Data(http.StatusOK, catImage.ContentType, imageData)
}

func (h *CatHandler) GetCount(c *gin.Context) {
	log.Println("GET /api/count")

	count, err := h.catService.GetUniqueImageCount()
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get count",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

func (h *CatHandler) GetStats(c *gin.Context) {
	log.Println("GET /api/stats")

	stats, err := h.catService.GetStats()
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get stats",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *CatHandler) GetImageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	log.Printf("GET /api/image/%d", id)

	catImage, err := h.catService.GetCatImageByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
		return
	}

	c.Header("X-Image-Hash", catImage.ImageHash)
	c.Data(http.StatusOK, catImage.ContentType, catImage.ImageData)
}