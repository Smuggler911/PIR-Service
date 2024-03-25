package repository

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CreateCommentary(c *gin.Context) {
	articleId := c.Param("article_id")

	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var article models.Article
	config.DB.Preload("ArticleCommentary").First(&article, articleId)

	var comments []models.ArticleCommentary
	var comment models.ArticleCommentary
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		return
	}
	commentBody := models.ArticleCommentary{
		Content:   comment.Content,
		UserID:    userId,
		ArticleID: uint(convertStringUint(articleId)),
	}
	comments = append(comments, commentBody)
	config.DB.Create(&commentBody)

	article.ArticleCommentary = comments

	config.DB.Model(&article).Updates(models.Article{
		ArticleCommentary: article.ArticleCommentary,
	})

	c.JSON(200, gin.H{
		"result": "ok",
	})

}
func GetCommentsByArticle(c *gin.Context) {
	articleId := c.Param("article_id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	offset := (page - 1) * limit

	var comments []models.ArticleCommentary
	var users []models.User
	config.DB.Where("is_banned", false).Find(&users)
	for _, user := range users {
		if !user.IsBanned {
			config.DB.Where("article_id = ?", articleId).Preload("User").Preload("Likes").Preload("Dislikes").Limit(limit).Offset(offset).Find(&comments)
			for i := 0; i < len(comments); i++ {
				comments[i].LikeCount = int64(likeCount(comments[i].Likes))
				comments[i].DislikeCount = int64(dislikeCount(comments[i].Dislikes))
			}
			c.JSON(200, gin.H{
				"comments": comments,
			})
		}
	}
}

func DeleteCommentary(c *gin.Context) {
	articleId := c.Param("article_id")

	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	commentaryId := c.Param("commentary_id")

	var article models.Article
	var comment models.ArticleCommentary
	config.DB.Preload("ArticleCommentary").First(&article, articleId)

	err := config.DB.Model(&article).Association("ArticleCommentary").Delete(models.ArticleCommentary{Id: uint(convertStringUint(commentaryId))})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "не удалось разорвать связь",
		})
		return
	}
	config.DB.Where("user_id = ?", userId).Delete(&comment, commentaryId)
	c.JSON(200, gin.H{"message": "deleted"})

}
func LikeCommentary(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	commenatryId := c.Param("commentary_id")

	var commnetary models.ArticleCommentary
	var likes []models.Like

	config.DB.Preload("Likes").Preload("Dislikes").First(&commnetary, commenatryId)
	for _, like := range commnetary.Likes {

		if userId == like.UserID {

			c.JSON(400, gin.H{
				"message": "вы  лайкали этот пост",
			})
			return

		}
	}
	if commnetary.Dislikes == nil {
		c.JSON(500, gin.H{
			"message": "что то не так",
		})
		return
	}
	if commnetary.Dislikes != nil && len(commnetary.Dislikes) != 0 {
		for _, dislike := range commnetary.Dislikes {
			if dislike.UserID == userId {
				err := config.DB.Model(&commnetary).Association("Dislikes").Delete(&dislike)
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
	commnetary.Likes = likes

	config.DB.Model(&commnetary).Updates(models.ArticleCommentary{
		Likes: commnetary.Likes,
	})

	c.JSON(200, gin.H{
		"message": "liked",
	})

}
func DislikeCommentary(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	commentaryId := c.Param("commentary_id")
	var commentary models.ArticleCommentary
	var dislikes []models.Dislike

	config.DB.Preload("Likes").Preload("Dislikes").First(&commentary, commentaryId)

	for _, dislike := range commentary.Dislikes {

		if userId == dislike.UserID {

			c.JSON(400, gin.H{
				"message": "настолько не понравился отзыв что вы решили влепить кучу дизлайков?",
			})
			return
		}
	}
	if commentary.Likes == nil {
		c.JSON(500, gin.H{
			"message": "что-то не так",
		})
		return
	}
	for _, like := range commentary.Likes {
		if like.UserID == userId && len(commentary.Likes) != 0 && commentary.Likes != nil {
			err := config.DB.Model(&commentary).Association("Likes").Delete(&like)
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

	commentary.Dislikes = dislikes

	config.DB.Model(&commentary).Updates(models.ArticleCommentary{
		Dislikes: commentary.Dislikes,
	})
	c.JSON(200, gin.H{
		"message": "disliked",
	})
}
