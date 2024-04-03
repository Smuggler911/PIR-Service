package repository

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func CreateMakeUp(c *gin.Context) {
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
	var makeupBody models.Makeup
	err := c.ShouldBindJSON(&makeupBody)
	if err != nil {
		return
	}
	makeup := models.Makeup{
		Name: makeupBody.Name,
	}
	result := config.DB.Create(&makeup)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
			"error":   result.Error,
		})
	}
	c.Status(200)

}
func CreateStylization(c *gin.Context) {
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
	var stylizationBody models.Stylization
	err := c.ShouldBindJSON(&stylizationBody)
	if err != nil {
		return
	}
	stylization := models.Stylization{
		Name: stylizationBody.Name,
	}
	result := config.DB.Create(&stylization)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
			"error":   result.Error,
		})
	}
	c.Status(200)

}
func CrateApplication(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	form, _ := c.MultipartForm()
	files := form.File["ref_pictures"]

	var refpicture []models.Pictures

	for _, file := range files {

		extension := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + extension
		env, _ := config.LoadConfig()
		imgPath := env.ImgPath
		destinationPath := imgPath + newFileName
		err := c.SaveUploadedFile(file, destinationPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
				"error":   err.Error(),
			})
			return
		}
		refpicture = append(refpicture, models.Pictures{
			Picture: newFileName,
		})
	}
	var facepicture []models.Pictures
	facepic := form.File["face_pictures"]
	for _, file := range facepic {

		extension := filepath.Ext(file.Filename)
		FaceFileName := uuid.New().String() + extension
		env, _ := config.LoadConfig()
		imgPath := env.ImgPath
		destinationPath := imgPath + FaceFileName
		err := c.SaveUploadedFile(file, destinationPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
				"error":   err.Error(),
			})
			return
		}
		facepicture = append(facepicture, models.Pictures{
			Picture: FaceFileName,
		})
	}

	application := models.Application{
		Day:               c.PostForm("day"),
		Time:              c.PostForm("time"),
		CreatedAt:         time.Now().UTC(),
		UserID:            userId,
		StylizationID:     uint(convertStringUint(c.PostForm("stylization_id"))),
		MakeupID:          uint(convertStringUint(c.PostForm("makeup_id"))),
		Preferences:       c.PostForm("preferences"),
		ReferencePictures: refpicture,
		FacePictures:      facepicture,
	}

	result := config.DB.Create(&application)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при cоздании заявки ",
		})
		return
	}
	c.Status(200)
}

func GetAllUserApplications(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application
	config.DB.Where("user_id = ?", userId).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}

func GetAllApplications(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application
	config.DB.Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}

func GetDoneApplications(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application
	config.DB.Where("done = ?", true).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}
func GetInProgressApplications(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application
	config.DB.Where("is_in_progress = ?", true).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}
func GetDeclinedApplications(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application
	config.DB.Where("is_declined = ?", true).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}
func GetNewApplications(c *gin.Context) {
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

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var applications []models.Application

	config.DB.Order("created_at DESC").Where("created_at <= ?", time.Now().UTC()).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").Limit(limit).Offset(offset).Find(&applications)
	c.JSON(200, gin.H{
		"applications": applications,
	})
}

func CancelApplication(c *gin.Context) {

	id := c.Param("id")
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var application models.Application
	config.DB.Where("is_in_progress = ? AND user_id = ?", false, userId).Preload("User").Preload("Stylization").Preload("Makeup").Preload("ReferencePictures").Preload("FacePictures").First(&application, id)
	config.DB.Model(&application).Association("User").Clear()
	config.DB.Model(&application).Association("Stylization").Clear()
	config.DB.Model(&application).Association("FacePictures").Clear()
	config.DB.Model(&application).Association("ReferencePictures").Clear()
	config.DB.Model(&application).Association("Makeup").Clear()
	config.DB.Delete(&application)

	c.Status(200)

}

func SetIsInProgress(c *gin.Context) {

	id := c.Param("id")

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
	var application models.Application

	config.DB.First(&application, id)
	config.DB.Model(&application).Updates(&models.Application{
		IsInProgress: true,
	})
	c.JSON(200, gin.H{
		"message": "работе в прогрессе",
	})

}
func SetIsDeclined(c *gin.Context) {

	id := c.Param("id")

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
	var application models.Application

	config.DB.First(&application, id)
	config.DB.Model(&application).Updates(&models.Application{
		IsDeclined: true,
	})
	c.JSON(200, gin.H{
		"message": "работе в прогрессе",
	})

}
func SetIsDone(c *gin.Context) {

	id := c.Param("id")

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
	var application models.Application

	config.DB.First(&application, id)
	config.DB.Model(&application).Updates(&models.Application{
		Done: true,
	})
	c.JSON(200, gin.H{
		"message": "работе в прогрессе",
	})

}
func DeleteStylizationCategory(c *gin.Context) {
	stylizaytionId := c.Param("id")

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
	var stylization models.Stylization
	config.DB.First(&stylization, stylizaytionId)
	var applications []models.Application
	config.DB.Model(&applications).Association("Stylization").Delete(&stylization)
	config.DB.Delete(&stylization)
}
func DeleteMakeupCategory(c *gin.Context) {
	makeupId := c.Param("id")

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
	var makeup models.Makeup
	config.DB.First(&makeup, makeupId)
	var applications []models.Application
	config.DB.Model(&applications).Association("Stylization").Delete(&makeup)
	config.DB.Delete(&makeup)
}
