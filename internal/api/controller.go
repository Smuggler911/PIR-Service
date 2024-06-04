package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"pir-serverSide/internal/api/middleware"
	"pir-serverSide/internal/repository"
	"time"
)

type User interface {
}
type Review interface {
}
type Article interface {
}
type Commentary interface {
}
type Handler struct {
	User
	Review
	Commentary
	Article
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))
	router.Static("/picture/", "/app/pkg/files")

	fmt.Println("......auth......")
	auth := router.Group("/auth")
	{

		auth.POST("/register", repository.RegisterNewUser)
		auth.POST("/login", repository.Login)
		auth.GET("/validate", middleware.RequireAuth, repository.Validate)
		auth.POST("/logout", middleware.RequireAuth, repository.Logout)

	}
	fmt.Println("......profile......")
	profile := router.Group("/profile")
	{
		profile.PUT("/profile", middleware.RequireAuth, repository.EditProfileBanner)
		profile.PUT("/content-update", middleware.RequireAuth, repository.EditProfileContent)
		profile.DELETE("/delete", middleware.RequireAuth, repository.DeleteProfile)
		profile.POST("/subscribe/:creatorId", middleware.RequireAuth, repository.Subscribe)
		profile.GET("/subscribers", middleware.RequireAuth, repository.GetSubscribers)
		profile.GET("/notifications", middleware.RequireAuth, repository.GetNotifications)

	}

	fmt.Println("...article....")
	article := router.Group("/article")
	{

		article.GET("/:id", repository.GetArticle)
		article.GET("/interesting", repository.GetInterestingArticles)
		article.PUT("/like-article/:article_id", middleware.RequireAuth, repository.LikeArticle)
		article.POST("/create-commentary/:article_id", middleware.RequireAuth, repository.CreateCommentary)
		article.GET("/comments/:article_id", repository.GetCommentsByArticle)
		article.DELETE("/commentary-delete/:article_id/:commentary_id", middleware.RequireAuth, repository.DeleteCommentary)
		article.PUT("/like-commentary/:commentary_id", middleware.RequireAuth, repository.LikeCommentary)
		article.GET("/new-articles/", repository.GetNewerArticles)
		article.DELETE("/delete-article/:article_id", middleware.RequireAuth, repository.DeleteArticle)
		article.PUT("/edit-article/:id", middleware.RequireAuth, repository.EditArticle)
		article.PUT("/view-article/:article_id", middleware.RequireAuth, repository.ViewArticle)
		article.GET("/search-article/:city/:title", repository.SearchForArticle)

	}

	fmt.Println("....admin-panel...")
	admin := router.Group("/admin-panel")
	{

		admin.PUT("/ban-user/:id", middleware.RequireAuth, repository.BanUser)
		admin.PUT("/unban-user/:id", middleware.RequireAuth, repository.UnbanUser)
		admin.GET("/banned", middleware.RequireAuth, repository.GetBannedUsers)
		admin.POST("/create-article", middleware.RequireAuth, repository.CreateArticle)
		admin.PUT("/ban-comment/:id", middleware.RequireAuth, repository.BanCommentary)
		admin.PUT("/publish-comment/:id", middleware.RequireAuth, repository.PublishCommentary)
		admin.GET("/get-banned-comment", middleware.RequireAuth, repository.GetBanned)
		admin.GET("/get-published-comment", middleware.RequireAuth, repository.GetPublished)
		article.PUT("/ban-article/:article_id", middleware.RequireAuth, repository.BanArticle)
		article.PUT("/publish-article/:article_id", middleware.RequireAuth, repository.PublishArticles)
		article.GET("/get-banned-article", middleware.RequireAuth, repository.GetBannedArticles)
		article.GET("/get-published-article", middleware.RequireAuth, repository.GetPublishedArticles)
	}

	return router

}
