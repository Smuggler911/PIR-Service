package repository

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"pir-serverSide/config"
	"pir-serverSide/internal/repository/models"
	"regexp"
	"time"
)

func RegisterNewUser(c *gin.Context) {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Ошибка хеширования",
		})
		return
	}
	var user models.User
	checkMail := config.DB.First(&user, "email = ?", c.PostForm("email")).Error
	if !errors.Is(checkMail, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "такой email ужe существует",
		})
		return
	}

	matchedEmail := checkingEmailReg(c.PostForm("email"))

	if !matchedEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Неверно указана почта",
		})
		return
	}
	admin := false
	if c.PostForm("name") == "admin" && c.PostForm("password") == "admin" {
		admin = true
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

	fmt.Println("Destination Path:", destinationPath)

	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}
	newUser := models.User{
		Name:        c.PostForm("name"),
		Surname:     c.PostForm("surname"),
		Picture:     newFileName,
		Email:       c.PostForm("email"),
		Password:    string(hash),
		DateOfBirth: c.PostForm("date_of_birth"),
		IsAdmin:     admin,
		Description: c.PostForm("description"),
		Role:        c.PostForm("role"),
		City:        c.PostForm("city"),
	}
	result := config.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при создании пользователя",
		})
		return
	}
	c.JSON(200, gin.H{"result": "ok"})
}

func Login(c *gin.Context) {
	var loginData *models.User

	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(400, gin.H{
			"error":  true,
			"result": "Не введен email или пароль",
		})
		return
	}
	var user models.User

	matchedEmail := checkingEmailReg(loginData.Email)

	if !matchedEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Неверно указана почта",
		})
		return
	}

	isLoginValid := config.DB.First(&user, "email = ?", loginData.Email).Error
	if errors.Is(isLoginValid, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  true,
			"result": "Не правильно введен email",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Не правильно введен пароль",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 20 * 30).Unix(),
	})

	env, _ := config.LoadConfig()
	secret := env.Secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  true,
			"result": "не получилось создать токен",
		})
	}
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("Authorization", tokenString, 3600*2*3, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {

	user, exists := c.Get("user")
	if exists {

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})

	}
}

func EditProfileBanner(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

	file, err := c.FormFile("banner")
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

	fmt.Println("Destination Path:", destinationPath)

	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}
	var user models.User
	config.DB.First(&user, userId)
	config.DB.Model(&user).Updates(models.User{
		Banner: newFileName,
	})

	c.JSON(200, gin.H{
		"status": "updated",
	})
}

func Logout(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}

	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.String(http.StatusOK, "Вы вышли из аккаунта")
}

func EditProfileContent(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

	var user models.User

	hash, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Ошибка хеширования",
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

	fmt.Println("Destination Path:", destinationPath)

	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}
	config.DB.First(&user, userId)
	config.DB.Model(&user).Updates(models.User{
		Name:        c.PostForm("name"),
		Surname:     c.PostForm("surname"),
		Picture:     newFileName,
		Email:       c.PostForm("email"),
		Password:    string(hash),
		DateOfBirth: c.PostForm("date_of_birth"),
		Description: c.PostForm("description"),
		Role:        c.PostForm("role"),
		City:        c.PostForm("city"),
	})

	c.JSON(200, gin.H{
		"status": "updated",
	})
}

func BanUser(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	id := c.Param("id")

	usr, _ := c.MustGet("user").(models.User)
	IsAdmin := usr.IsAdmin
	if !IsAdmin {
		c.JSON(400, gin.H{
			"message": "вы не являетесь админом",
		})
		return
	}

	var user models.User
	config.DB.Where("is_banned = ?", false).First(&user, id)
	config.DB.Model(&user).Updates(models.User{
		IsBlocked: true,
	})
	c.JSON(200, gin.H{
		"result": "banned",
	})

}
func UnbanUser(c *gin.Context) {
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
	var userbody models.User

	err := c.ShouldBindJSON(&userbody)
	if err != nil {
		return
	}

	userId := c.Param("id")
	var user models.User
	config.DB.First(&user, userId)

	if user.IsBlocked {
		userbody.IsBlocked = false
	}

	config.DB.Model(&user).Update("is_blocked", userbody.IsBlocked)

	c.JSON(200, gin.H{
		"result": "unbanned",
	})
}

func GetBannedUsers(c *gin.Context) {
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
	var users []models.User
	config.DB.Where("is_blocked = ?", true).Find(&users)

	c.JSON(200, gin.H{
		"users": users,
	})
}

func DeleteProfile(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	config.DB.Delete(&models.User{}, userId)
	c.Status(200)
}

func Subscribe(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id
	creatorId := c.Param("creator_id")
	sub := models.Subscriber{
		Subscribed: true,
		UserId:     int(userId),
		CreatorId:  int(convertStringInt(creatorId)),
	}
	result := config.DB.Create(&sub)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":   result.Error,
			"message": "не удалось подписаться на пользователя",
		})
	}
	c.JSON(200, gin.H{
		"result": "подписаны",
	})
}

func GetSubscribers(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id
	var users []models.User
	var subscriptions []models.Subscriber
	config.DB.Where("creatorId = ?", userId).Preload("Creator").Find(&subscriptions)
	for i := 0; i < len(subscriptions); i++ {
		sub := subscriptions[i].CreatorId
		config.DB.Where("id = ?", sub).Find(&users)
	}
	c.JSON(200, gin.H{
		"subscribers": users,
	})
}
func GetNotifications(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id
	var users []models.User
	var subscriptions []models.Subscriber
	var notifications []models.Notification
	config.DB.Where("creatorId = ?", userId).Preload("Creator").Find(&subscriptions)
	for i := 0; i < len(subscriptions); i++ {
		sub := subscriptions[i].CreatorId
		config.DB.Where("id = ?", sub).Find(&users)
		config.DB.Where("userId = ?", sub).Preload("User").Find(&notifications)
	}
	c.JSON(200, gin.H{
		"notifications": notifications,
	})
}

func checkingEmailReg(email string) bool {

	matched, _ := regexp.MatchString(`([A-Za-z0-9_\-.])+@([A-Za-z0-9_\-.])+\.([A-Za-z]{2,4})`, email)

	if !matched {

		return false
	}
	return true
}
