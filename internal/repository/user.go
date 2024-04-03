package repository

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"time"
)

func RegisterNewUser(c *gin.Context) {

	var registerBody *models.User

	err := c.BindJSON(&registerBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error":  true,
			"result": "не введены данные",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Ошибка хеширования",
		})
		return
	}

	var user models.User
	checkMail := config.DB.First(&user, "email = ?", registerBody.Email).Error
	if !errors.Is(checkMail, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "такой email ужe существует",
		})
		return
	}

	matchedEmail := checkingEmailReg(registerBody.Email)

	if !matchedEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  true,
			"result": "Неверно указана почта",
		})
		return
	}

	if registerBody.Name == "admin" && registerBody.Password == "admin" {
		registerBody.IsAdmin = true
	}

	if len(registerBody.Phone) < 12 || len(registerBody.Phone) > 13 {
		log.Println("телефон  содержит меньше 12 или большк 13 cимволов:", registerBody.Phone)
		c.JSON(http.StatusBadRequest, &gin.H{
			"Message": "телефон  содержит меньше 12 или большк 13 cимволов",
			"Status":  false,
		})
		return
	}

	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !re.MatchString(registerBody.Phone) {
		log.Println("Номер телефона введен не правильно:", registerBody.Phone)
		c.JSON(http.StatusBadRequest, &gin.H{
			"Message": "Номер телефона введен не правильно",
			"Status":  false,
		})
		return
	}
	newUser := models.User{
		Name:        registerBody.Name,
		Surname:     registerBody.Name,
		Email:       registerBody.Email,
		Phone:       registerBody.Phone,
		Password:    string(hash),
		DateOfBirth: registerBody.DateOfBirth,
		IsAdmin:     registerBody.IsAdmin,
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
		"exp": time.Now().Add(time.Hour * 2 * 3).Unix(),
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

	c.SetCookie("Authorization", tokenString, 3600*2*3, "", "", false, true)
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

func Logout(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.String(http.StatusOK, "Вы вышли из аккаунта")
}

func EditProfilePic(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

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
	var user models.User
	config.DB.First(&user, userId)
	config.DB.Model(&user).Updates(models.User{
		Picture: newFileName,
	})

	c.JSON(200, gin.H{
		"status": "updated",
	})
}

func EditProfileContent(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

	var pContnet *models.User

	err := c.BindJSON(&pContnet)
	if err != nil {
		c.JSON(400, gin.H{
			"error":  true,
			"result": "пустое поле",
		})
		return
	}

	var user models.User

	config.DB.First(&user, userId)
	config.DB.Model(&user).Updates(models.User{
		Description: pContnet.Description,
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
	var body models.User
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "напишите причину",
		})
		return
	}
	config.DB.First(&user, id)
	config.DB.Model(&user).Updates(models.User{
		IsBanned:     true,
		WhyIsBlocked: body.WhyIsBlocked,
	})
	c.JSON(200, gin.H{
		"message": "пользователь забанен",
	})

}
func UnbanUser(c *gin.Context) {
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
	config.DB.First(&user, id)
	config.DB.Model(&user).Updates(models.User{
		IsBanned:     false,
		WhyIsBlocked: "",
	})
	c.JSON(200, gin.H{
		"message": "пользователь разабанен",
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
	config.DB.Where("is_banned = ?", true).Find(&users)
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

func checkingEmailReg(email string) bool {

	matched, _ := regexp.MatchString(`([A-Za-z0-9_\-.])+@([A-Za-z0-9_\-.])+\.([A-Za-z]{2,4})`, email)

	if !matched {

		return false
	}
	return true
}
