package repository

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func CreateArticle(c *gin.Context) {
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
	form, _ := c.MultipartForm()
	files := form.File["pictures"]

	var picture []models.Pictures

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
		picture = append(picture, models.Pictures{
			Picture: newFileName,
		})
	}

	if len(picture) > 6 {
		c.JSON(400, gin.H{
			"message": "нельзя добавлять больше 6 файлов",
		})
		return
	}

	articleBody := models.Article{
		ArticlePictures: picture,
		Title:           c.PostForm("title"),
		ACategoryID:     uint(convertStringUint(c.PostForm("a_category_id"))),
		Description:     c.PostForm("description"),
		Content:         c.PostForm("content"),
	}
	result := config.DB.Create(&articleBody)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при cоздании статьи ",
		})
		return
	}
	c.JSON(200, gin.H{"result": "ok"})

}

func GetAllArticles(c *gin.Context) {
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
	config.DB.Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Limit(limit).Offset(offset).Find(&articles)
	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].DislikeCount = int64(dislikeCount(articles[i].Dislikes))
	}
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func GetArticlesByCategory(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit
	categoryId := c.Param("category_id")

	var articles []models.Article
	config.DB.Where("a_category_id = ?", categoryId).Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Limit(limit).Offset(offset).Find(&articles)
	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].DislikeCount = int64(dislikeCount(articles[i].Dislikes))
	}
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func GetArticle(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit
	id := c.Param("id")
	var article models.Article
	config.DB.Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").First(&article, id)
	categoryId := article.ACategoryID

	article.LikeCount = int64(likeCount(article.Likes))
	article.DislikeCount = int64(dislikeCount(article.Dislikes))

	var articles []models.Article
	config.DB.Where("a_category_id = ?", categoryId).Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Limit(limit).Offset(offset).Find(&articles)
	for i := 0; i < len(articles); i++ {
		articles[i].LikeCount = int64(likeCount(articles[i].Likes))
		articles[i].DislikeCount = int64(dislikeCount(articles[i].Dislikes))
	}
	c.JSON(200, gin.H{
		"article":  article,
		"articles": articles,
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
	config.DB.Where("like_count > ?", 1).Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Limit(limit).Offset(offset).Find(&articles)
	c.JSON(200, gin.H{
		"articles": articles,
	})
}

func DeletePictureFromArticle(c *gin.Context) {
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
	articleId := c.Param("article_id")
	pictureId := c.Param("picture_id")

	var article models.Article
	if err := config.DB.Where("id = ?", articleId).Preload("ArticlePictures").Find(&article).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "не найдена такая запись",
		})
		return
	}
	var foundPicture bool
	for i, picture := range article.ArticlePictures {
		if picture.Id == convertStringUint(pictureId) {

			article.ArticlePictures = append(article.ArticlePictures[:i], article.ArticlePictures[i+1:]...)
			foundPicture = true
			break
		}
	}
	if !foundPicture {
		c.JSON(400, gin.H{
			"message": "картинка не найдена",
		})
		return
	}
	if err := config.DB.Save(&article).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка обновления  модели статьи ",
		})
		return
	}
	err := config.DB.Model(&article).Association("ArticlePictures").Delete(&models.Pictures{Id: convertStringUint(pictureId)})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка удаления связи между cтатьей и картинкой",
		})
		return
	}
	var picture models.Pictures
	config.DB.Where("id = ? ", pictureId).Find(&picture)

	env, _ := config.LoadConfig()
	imgPath := env.ImgPath

	err = os.Remove(imgPath + picture.Picture)
	if err != nil {
		fmt.Println("oшибка удаление файла")
		return
	}

	if err := config.DB.Delete(&models.Pictures{}, pictureId).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка удаления картинки",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Картинка успешно удалена",
	})

}

func EditArticle(c *gin.Context) {
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

	form, _ := c.MultipartForm()
	files := form.File["pictures"]
	var article models.Article

	var picture []models.Pictures
	config.DB.Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Where(" id = ?", id).Find(&article)

	article.ArticlePictures = picture
	err := config.DB.Model(&article).Association("ArticlePictures").Clear()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "не удалось очистить связи",
		})
		return
	}

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
		picture = append(picture, models.Pictures{
			Picture: newFileName,
		})
	}

	config.DB.Preload("ArticlesCategory").Preload("Likes").Preload("Dislikes").Preload("ArticleCommentary").Preload("ArticlePictures").Where(" id = ?", id).Find(&article)

	article.ArticlePictures = picture

	config.DB.Model(&article).Updates(models.Article{
		ArticlePictures: article.ArticlePictures,
		Title:           c.PostForm("title"),
		ACategoryID:     uint(convertStringUint(c.PostForm("a_category_id"))),
		Description:     c.PostForm("description"),
		Content:         c.PostForm("content"),
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
	if artice.Dislikes == nil {
		c.JSON(500, gin.H{
			"message": "что то не так",
		})
		return
	}
	if artice.Dislikes != nil && len(artice.Dislikes) != 0 {
		for _, dislike := range artice.Dislikes {
			if dislike.UserID == userId {
				err := config.DB.Model(&artice).Association("Dislikes").Delete(&dislike)
				if err != nil {
					c.JSON(500, gin.H{
						"message": err,
					})
					return
				}
				break
			}
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

func DislikeArticle(c *gin.Context) {
	articleID := c.Param("article_id")
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var article models.Article
	var dislikes []models.Dislike

	config.DB.Preload("Likes").Preload("Dislikes").First(&article, articleID)

	for _, dislike := range article.Dislikes {

		if userId == dislike.UserID {

			c.JSON(400, gin.H{
				"message": "настолько не понравился отзыв что вы решили влепить кучу дизлайков?",
			})
			return
		}
	}
	if article.Likes == nil {
		c.JSON(500, gin.H{
			"message": "что-то не так",
		})
		return
	}
	for _, like := range article.Likes {
		if like.UserID == userId && len(article.Likes) != 0 && article.Likes != nil {
			err := config.DB.Model(&article).Association("Likes").Delete(&like)
			if err != nil {
				c.JSON(500, gin.H{
					"message": err,
				})
				return
			}
		}
	}
	disliked := 1

	dislikes = append(
		dislikes, models.Dislike{
			Disliked: int64(disliked),
			UserID:   userId,
		},
	)

	article.Dislikes = dislikes

	config.DB.Model(&article).Updates(models.Article{
		Dislikes: article.Dislikes,
	})
	c.JSON(200, gin.H{
		"message": "disliked",
	})
}

func CreateArticlesCategory(c *gin.Context) {
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
	var categoryBody models.ArticlesCategory
	err := c.ShouldBindJSON(&categoryBody)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "должно быть название категории",
		})
		return
	}
	request := models.ArticlesCategory{
		Name: categoryBody.Name,
	}
	result := config.DB.Create(&request)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при cоздании категории",
		})
		return
	}
	c.JSON(200, gin.H{"result": "ok"})
}

func GetArticlesCategory(c *gin.Context) {
	var aCategories []models.ArticlesCategory
	config.DB.Find(&aCategories)
	c.JSON(200, gin.H{
		"article_categories": aCategories,
	})
}

func DeleteArticlesCategory(c *gin.Context) {
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
	err := config.DB.Where("a_category_id = ?", id).Model(&article).Association("ArticlesCategory").Clear()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "не удалось разорвать связь ",
		})
		return
	}
	config.DB.Delete(&models.ArticlesCategory{}, id)
	c.JSON(200, gin.H{
		"message": "deleted",
	})

}
