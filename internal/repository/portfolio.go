package repository

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strconv"
)

func AddPortfolioWork(c *gin.Context) {

	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	IsAdmin := usr.IsAdmin
	if !IsAdmin {
		c.JSON(400, gin.H{
			"message": "вы не являетесь админом",
		})
		return
	}

	file, err := c.FormFile("picture")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	env, _ := config.LoadConfig()
	imgPath := env.ImgPath

	destinationPath := imgPath + newFileName
	if err = c.SaveUploadedFile(file, destinationPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}
	protfolio := models.Portfolio{
		Picture: newFileName,
	}
	result := config.DB.Create(&protfolio)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка при создании пользователя",
			"error":   result.Error,
		})
	}
	c.JSON(200, gin.H{
		"message": "добавлено в портфолио",
	})

}

func DeletePortfolioWork(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	IsAdmin := usr.IsAdmin
	if !IsAdmin {
		c.JSON(400, gin.H{
			"message": "вы не являетесь админом",
		})
		return
	}
	id := c.Param("id")

	config.DB.Delete(&models.Portfolio{}, id)
	c.Status(200)
}

func GetAllPortfolios(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit
	var portfolio []models.Portfolio
	config.DB.Limit(limit).Offset(offset).Find(&portfolio)
	c.JSON(200, gin.H{
		"portfolio": portfolio,
	})
}
