package main

import (
	"pir-serverSide/config"
	"pir-serverSide/internal/repository/models"
)

func init() {
	config.ConnectDB()
}

func main() {

	err := config.DB.AutoMigrate(
		&models.User{}, &models.Subscriber{}, &models.Like{}, &models.Views{},
		&models.Article{}, &models.ArticleCommentary{}, &models.Notification{},
	)
	if err != nil {
		return
	}

}
