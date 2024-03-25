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

func CreateReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
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

	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	reviewBody := models.Review{
		ReviewPictures: picture,
		UserID:         userId,
		Rate:           convertStringInt(c.PostForm("rate")),
		ReviewText:     c.PostForm("review_text"),
	}
	result := config.DB.Create(&reviewBody)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":  result.Error,
			"result": "ошибка при оставлении отзыва",
		})
		return
	}

	c.JSON(200, gin.H{"result": "ok"})
}

func GetReviews(c *gin.Context) {
	var reviews []models.Review
	config.DB.Preload("User").Preload("Likes").Preload("Dislikes").Preload("ReviewPictures").Find(&reviews)
	for i := 0; i < len(reviews); i++ {
		reviews[i].LikeCount = int64(likeCount(reviews[i].Likes))
		reviews[i].DislikeCount = int64(dislikeCount(reviews[i].Dislikes))
	}
	c.JSON(200, gin.H{
		"reviews": reviews,
	})
}

func DeleteReviewPicture(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	reviewId := c.Param("review_id")
	pictureId := c.Param("picture_id")

	var review models.Review

	if err := config.DB.Where("user_id = ? AND id = ?", userId, reviewId).Preload("ReviewPictures").Find(&review).Error; err != nil {
		c.JSON(400, gin.H{
			"message": "не найдена такая запись",
		})
		return
	}

	var foundPicture bool
	for i, picture := range review.ReviewPictures {
		if picture.Id == convertStringUint(pictureId) {

			review.ReviewPictures = append(review.ReviewPictures[:i], review.ReviewPictures[i+1:]...)
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

	if err := config.DB.Save(&review).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка обновления ревью модели",
		})
		return
	}

	err := config.DB.Model(&review).Association("ReviewPictures").Delete(&models.Pictures{Id: convertStringUint(pictureId)})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка удаления связи между ревью и картинкой",
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

func EditReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id
	reviewId := c.Param("review_id")

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
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

	var review models.Review

	config.DB.Preload("ReviewPictures").Where("user_id = ? AND id = ?", userId, reviewId).Find(&review)

	review.ReviewPictures = picture

	config.DB.Model(&review).Updates(models.Review{
		ReviewPictures: review.ReviewPictures,
		Rate:           convertStringInt(c.PostForm("rate")),
		ReviewText:     c.PostForm("review_text"),
	})

	c.JSON(200, gin.H{
		"status": "обновлено",
	})

}

func DeleteReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	id := c.Param("id")
	user, _ := c.MustGet("user").(models.User)
	isAdmin := user.IsAdmin

	if !isAdmin {
		c.JSON(400, gin.H{
			"message": "вы не являтесь админом",
		})
		return
	}

	var review models.Review
	err := config.DB.Preload("ReviewPictures").First(&review, id).Error

	if err != nil {
		c.JSON(400, gin.H{
			"message": " запись не найдена ",
		})
		return
	}
	var pictures []models.Pictures
	err = config.DB.Model(&review).Association("ReviewPictures").Find(&pictures)
	env, _ := config.LoadConfig()
	imgPath := env.ImgPath

	for _, picture := range pictures {
		fmt.Println(picture.Picture)
		err = os.Remove(imgPath + picture.Picture)
		if err != nil {
			fmt.Println("oшибка удаление файла")
			return
		}
	}
	err = config.DB.Model(&review).Association("ReviewPictures").Clear()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "не удалось удалить взаимоотношения таблицы",
		})
		return
	}
	config.DB.Delete(&pictures)

	config.DB.Delete(&review)

	c.JSON(200, gin.H{
		"result": "удалено",
	})

}

func DeleteYourReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	reviewId := c.Param("review_id")
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var review models.Review

	err := config.DB.Preload("ReviewPictures").Where("user_id = ? AND id = ?", userId, reviewId).First(&review).Error
	if err != nil {
		c.JSON(400, gin.H{
			"message": " запись не найдена ",
		})
		return
	}

	var pictures []models.Pictures
	err = config.DB.Model(&review).Association("ReviewPictures").Find(&pictures)
	env, _ := config.LoadConfig()
	imgPath := env.ImgPath

	for _, picture := range pictures {
		fmt.Println(picture.Picture)
		err = os.Remove(imgPath + picture.Picture)
		if err != nil {
			fmt.Println("oшибка удаление файла")
			return
		}
	}
	err = config.DB.Model(&review).Association("ReviewPictures").Clear()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "не удалось удалить взаимоотношения таблицы",
		})
		return
	}

	config.DB.Delete(&pictures)

	config.DB.Delete(&review)

	c.JSON(200, gin.H{
		"result": "удалено",
	})
}

func OverallRating(c *gin.Context) {

	var ovrl models.OverallRating

	var reviews []models.Review
	config.DB.Find(&ovrl)
	config.DB.Find(&reviews)

	ovrl.Overall = calcAveragePoints(reviews)
	config.DB.Where("rate = ?", 5).Find(&reviews)
	if len(reviews) == 0 {
		ovrl.FiveStar = 0
	}
	ovrl.FiveStar = float64(len(reviews))
	config.DB.Where("rate = ?", 4).Find(&reviews)
	if len(reviews) == 0 {
		ovrl.FourStar = 0
	}
	ovrl.FourStar = float64(len(reviews))
	config.DB.Where("rate = ?", 3).Find(&reviews)
	if len(reviews) == 0 {
		ovrl.ThirdStar = 0
	}
	ovrl.ThirdStar = float64(len(reviews))
	config.DB.Where("rate = ?", 2).Find(&reviews)
	if len(reviews) == 0 {
		ovrl.TwoStar = 0
	}
	ovrl.TwoStar = float64(len(reviews))
	config.DB.Where("rate = ?", 1).Find(&reviews)
	if len(reviews) == 0 {
		ovrl.OneStar = 0
	}
	ovrl.OneStar = float64(len(reviews))

	if err := config.DB.Save(&ovrl).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "ошибка обновления  модели",
		})
		return
	}

	c.JSON(200, gin.H{
		"overall": ovrl,
	})
}

func LikeReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	reviewId := c.Param("review_id")
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var review models.Review
	var likes []models.Like

	config.DB.Preload("Likes").Preload("Dislikes").First(&review, reviewId)

	for _, like := range review.Likes {

		if userId == like.UserID {

			c.JSON(400, gin.H{
				"message": "вы  лайкали этот пост",
			})
			return

		}
	}
	if review.Dislikes == nil {
		c.JSON(500, gin.H{
			"message": "что то не так",
		})
		return
	}

	if review.Dislikes != nil && len(review.Dislikes) != 0 {
		for _, dislike := range review.Dislikes {
			if dislike.UserID == userId {
				err := config.DB.Model(&review).Association("Dislikes").Delete(&dislike)
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

	review.Likes = likes

	config.DB.Model(&review).Updates(models.Review{
		Likes: review.Likes,
	})

	c.JSON(200, gin.H{
		"message": "liked",
	})
}

func DislikeReview(c *gin.Context) {
	exists, _ := c.Get("user")
	if exists == nil {
		return
	}
	reviewId := c.Param("review_id")
	user, _ := c.MustGet("user").(models.User)
	userId := user.Id

	var review models.Review
	var dislikes []models.Dislike

	config.DB.Preload("Likes").Preload("Dislikes").First(&review, reviewId)

	for _, dislike := range review.Dislikes {

		if userId == dislike.UserID {

			c.JSON(400, gin.H{
				"message": "настолько не понравился отзыв что вы решили влепить кучу дизлайков?",
			})
			return
		}
	}
	if review.Likes == nil {
		c.JSON(500, gin.H{
			"message": "что-то не так",
		})
		return
	}

	for _, like := range review.Likes {
		if like.UserID == userId && len(review.Likes) != 0 && review.Likes != nil {
			err := config.DB.Model(&review).Association("Likes").Delete(&like)
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

	review.Dislikes = dislikes

	config.DB.Model(&review).Updates(models.Review{
		Dislikes: review.Dislikes,
	})
	c.JSON(200, gin.H{
		"message": "disliked",
	})
}

func convertStringUint(strValue string) uint64 {
	Value, err := strconv.ParseUint(strValue, 10, 64)
	if err != nil {
		return 0
	}
	return Value
}

func convertStringInt(strValue string) int32 {
	Value, err := strconv.ParseInt(strValue, 10, 32)
	if err != nil {
		return 0
	}
	return int32(Value)

}
func calcAveragePoints(reviews []models.Review) float64 {
	if len(reviews) != 0 {
		result := float64(0)
		for i := 0; i < len(reviews); i++ {
			rate := reviews[i].Rate
			result = result + float64(rate)
		}
		return round(result/float64(len(reviews)), 1)
	} else {
		return 0
	}
}

func likeCount(likes []models.Like) float64 {
	if len(likes) != 0 {
		result := float64(0)
		for i := 0; i < len(likes); i++ {
			like := float64(likes[i].Liked)
			result = result + like
		}
		return result

	} else {
		return 0
	}
}
func dislikeCount(dislikes []models.Dislike) float64 {
	if len(dislikes) != 0 {
		result := float64(0)
		for i := 0; i < len(dislikes); i++ {
			dislike := float64(dislikes[i].Disliked)
			result = result + dislike
		}
		return result

	} else {
		return 0
	}
}

func round(number float64, decimals int) float64 {
	output := strconv.FormatFloat(number, 'f', decimals, 64)
	result, _ := strconv.ParseFloat(output, 64)
	return result
}
