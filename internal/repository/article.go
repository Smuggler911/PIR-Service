package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"pir-serverSide/config"
	"pir-serverSide/internal/repository/models"
	"strconv"
	"time"
)

func CreateArticle(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

	//Главная картинка
	file, err := c.FormFile("mainPic")
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

	//Картинка 1ой главы
	firstFile, err := c.FormFile("chapteronePic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(firstFile.Filename)
	ChapteronePicName := uuid.New().String() + extension

	destinationChapteronePic := imgPath + ChapteronePicName

	fmt.Println("Destination Path:", destinationChapteronePic)

	if err := c.SaveUploadedFile(firstFile, destinationChapteronePic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	// Картинка второй главы
	secondFile, err := c.FormFile("chaptertwoPic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(secondFile.Filename)
	chaptertwoPicName := uuid.New().String() + extension

	destinationChaptertwoPic := imgPath + chaptertwoPicName

	fmt.Println("Destination Path:", destinationChaptertwoPic)

	if err := c.SaveUploadedFile(secondFile, destinationChaptertwoPic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	// Картинка третьей главы

	thridFile, err := c.FormFile("chapterthreePic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(thridFile.Filename)
	chapterthreePicName := uuid.New().String() + extension

	destinationChapterthreePic := imgPath + chapterthreePicName

	fmt.Println("Destination Path:", destinationChapterthreePic)

	if err := c.SaveUploadedFile(thridFile, destinationChapterthreePic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	articleBody := models.Article{
		Title:           c.PostForm("title"),
		MainPic:         newFileName,
		ChapterOne:      c.PostForm("chapterOne"),
		ChapterOnePic:   ChapteronePicName,
		ChapterTwo:      c.PostForm("chapterTwo"),
		ChapterThree:    c.PostForm("chapterThree"),
		ChapterTwoPic:   chaptertwoPicName,
		ChapterThreePic: chaptertwoPicName,
		UserId:          int(userId),
		CreatedAt:       time.Now().UTC(),
	}
	result := config.DB.Create(&articleBody)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при cоздании статьи ",
		})
		return
	}
	notification := models.Notification{
		UserID:       int(userId),
		Notification: "Выложил(а) новый пост",
	}
	result = config.DB.Create(&notification)

	c.JSON(200, gin.H{"result": "ok"})

}

func GetNewerArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var articles []models.Article
	config.DB.Order("created_at desc, title").Preload("User").Preload("Likes").Preload("ArticleCommentary").Preload("Views").Limit(limit).Offset(offset).Find(&articles)

	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].ViewsCount = int64(viewCount(articles[i].Views))
		articles[i].CommentaryCount = int64(len(articles[i].ArticleCommentary))
	}
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	config.DB.Preload("User").Preload("Views").Preload("Likes").Preload("ArticleCommentary").First(&article, id)

	article.LikeCount = int64(likeCount(article.Likes))
	article.CommentaryCount = int64(len(article.ArticleCommentary))

	c.JSON(200, gin.H{
		"article": article,
	})
}

func GetInterestingArticles(c *gin.Context) {

	var articles []models.Article
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit
	config.DB.Where("like_count > ?", 10).Preload("User").Preload("Views").Preload("Likes").Preload("ArticleCommentary").Limit(limit).Offset(offset).Find(&articles)
	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].ViewsCount = int64(viewCount(articles[i].Views))
		articles[i].CommentaryCount = int64(len(articles[i].ArticleCommentary))
	}
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func DeleteArticle(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id

	articleId := c.Param("article_id")

	var article models.Article
	if err := config.DB.Where("id = ? AND user_Id", articleId, userId).Preload("ArticlePictures").Find(&article).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "не найдена такая запись",
		})
		return
	}
	if err := config.DB.Unscoped().Delete(&article).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка удаления картинки",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "cтатья успешно удалена",
	})

}

func EditArticle(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	usr, _ := c.MustGet("user").(models.User)
	userId := usr.Id
	id := c.Param("id")

	var article models.Article

	//Главная картинка
	file, err := c.FormFile("mainPic")
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

	//Картинка 1ой главы
	firstFile, err := c.FormFile("chapteronePic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(firstFile.Filename)
	ChapteronePicName := uuid.New().String() + extension

	destinationChapteronePic := imgPath + ChapteronePicName

	fmt.Println("Destination Path:", destinationChapteronePic)

	if err := c.SaveUploadedFile(firstFile, destinationChapteronePic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	// Картинка второй главы
	secondFile, err := c.FormFile("chaptertwoPic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(secondFile.Filename)
	chaptertwoPicName := uuid.New().String() + extension

	destinationChaptertwoPic := imgPath + chaptertwoPicName

	fmt.Println("Destination Path:", destinationChaptertwoPic)

	if err := c.SaveUploadedFile(secondFile, destinationChaptertwoPic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	// Картинка третьей главы

	thridFile, err := c.FormFile("chapterthreePic")
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "нет картинки",
		})
		return
	}

	extension = filepath.Ext(thridFile.Filename)
	chapterthreePicName := uuid.New().String() + extension

	destinationChapterthreePic := imgPath + chapterthreePicName

	fmt.Println("Destination Path:", destinationChapterthreePic)

	if err := c.SaveUploadedFile(thridFile, destinationChapterthreePic); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно сохранить картинку",
			"error":   err.Error(),
		})
		return
	}

	config.DB.Preload("User").Preload("Likes").Preload("ArticleCommentary").Preload("Views").Where(" id = ? AND user_id", id, userId).Find(&article)

	config.DB.Model(&article).Updates(models.Article{
		Title:           c.PostForm("title"),
		MainPic:         newFileName,
		ChapterOne:      c.PostForm("chapterOne"),
		ChapterOnePic:   ChapteronePicName,
		ChapterTwo:      c.PostForm("chapterTwo"),
		ChapterThree:    c.PostForm("chapterThree"),
		ChapterTwoPic:   chaptertwoPicName,
		ChapterThreePic: chaptertwoPicName,
	})

	c.JSON(200, gin.H{
		"status": "обновлено",
	})
}

func LikeArticle(c *gin.Context) {
	articleId := c.Param("article_id")
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var artice models.Article
	var likes []models.Like

	config.DB.Preload("Likes").Preload("Dislikes").First(&artice, articleId)
	for _, like := range artice.Likes {

		if userId == like.UserID {

			c.JSON(400, gin.H{
				"message": "вы  лайкали этот пост",
			})
			return
		}
	}
	liked := 1

	likes = append(
		likes, models.Like{
			Liked:  int64(liked),
			UserID: userId,
		},
	)
	artice.Likes = likes

	config.DB.Model(&artice).Updates(models.Article{
		Likes: artice.Likes,
	})

	c.JSON(200, gin.H{
		"message": "liked",
	})
}

func ViewArticle(c *gin.Context) {
	articleId := c.Param("article_id")
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var article models.Article
	var views []models.Views
	config.DB.Preload("Views").First(&article, articleId)

	for _, view := range views {
		if userId == view.UserID {
			c.JSON(200, gin.H{
				"message": "вы заходите на этот пост второй раз , возможно он вам понравился не хотите поддержать автора лайком?",
			})
			return
		}
	}
	viewed := 1
	views = append(views, models.Views{
		Viewed: int64(viewed),
		UserID: userId,
	})
	article.Views = views
	config.DB.Model(&article).Updates(models.Article{
		Views: article.Views,
	})
	c.JSON(200, gin.H{
		"message": "просмотрен",
	})
}

func SearchForArticle(c *gin.Context) {
	city := c.Param("city")
	title := c.Param("title")
	var articles []*models.Article
	config.DB.Where("city = ? or title = ?", city, title).Preload("User").Preload("Likes").Preload("ArticleCommentary").Preload("Views").Find(&articles)
	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].ViewsCount = int64(viewCount(articles[i].Views))
		articles[i].CommentaryCount = int64(len(articles[i].ArticleCommentary))
	}
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func BanArticle(c *gin.Context) {
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
	var article models.Article
	config.DB.Where("blocked = ?", false).First(&article, id)
	config.DB.Model(&article).Updates(models.Article{
		Blocked: true,
	})
	c.JSON(200, gin.H{
		"result": "banned",
	})
}

func PublishArticles(c *gin.Context) {
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
	id := c.Param("city")
	var article models.Article
	config.DB.Where("published = ?", false).First(&article, id)
	config.DB.Model(&article).Updates(models.Article{
		Published: true,
	})
	c.JSON(200, gin.H{
		"result": "banned",
	})
}

func GetBannedArticles(c *gin.Context) {
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
	var articles []models.Article
	config.DB.Where("blocked = ?", true).Find(&articles)

	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func GetPublishedArticles(c *gin.Context) {
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
	var articles []models.Article
	config.DB.Where("published = ?", true).Find(&articles)

	c.JSON(200, gin.H{
		"articles": articles,
	})
}
