package main

import (
	"MakeWish-serverSide/config"
	"MakeWish-serverSide/internal/repository/models"
)

func init() {
	config.ConnectDB()
}

func main() {

	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Pictures{},
		&models.Application{},
		&models.Makeup{},
		&models.Review{},
		&models.Stylization{},
		&models.Portfolio{},
		&models.ArticleCommentary{},
		&models.ArticlesCategory{},
		&models.OverallRating{},
		&models.Article{},
		&models.User{},
		&models.Review{},
		&models.Article{},
	)
	if err != nil {
		return
	}

}
