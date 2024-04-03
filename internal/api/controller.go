package api

import (
	"MakeWish-serverSide/internal/api/middleware"
	"MakeWish-serverSide/internal/repository"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type User interface {
	RegisterNewUser()
	Login()
	Validate()
	EditProfilePic()
	EditProfileContent()
	DeleteProfile()
	Logout()
	BanUser()
	UnbanUser()
	GetBannedUsers()
}
type Review interface {
	CreateReview()
	GetReviews()
	DeleteReviewPicture()
	EditReview()
	OverallRating()
	DeleteReview()
	DeleteYourReview()
	LikeReview()
	DislikeReview()
}
type Article interface {
	CreateArticle()
	CreateArticlesCategory()
	GetAllArticles()
	GetArticlesCategory()
	DeleteArticlesCategory()
	GetArticlesByCategory()
	GetArticle()
	EditArticle()
	DeletePictureFromArticle()
	LikeArticle()
	DislikeArticle()
	GetInterestingArticles()
}
type Commentary interface {
	CreateCommentary()
	GetCommentsByArticle()
	DeleteCommentary()
	LikeCommentary()
	DislikeCommentary()
}
type Handler struct {
	User
	Review
	Commentary
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
		profile.PUT("/picture-update", middleware.RequireAuth, repository.EditProfilePic)
		profile.PUT("/content-update", middleware.RequireAuth, repository.EditProfileContent)
		profile.DELETE("/delete", middleware.RequireAuth, repository.DeleteProfile)

	}
	fmt.Println("....review.....")
	review := router.Group("/review")
	{
		review.POST("/create", middleware.RequireAuth, repository.CreateReview)
		review.GET("", repository.GetReviews)
		review.DELETE("/delete-review-pic/:review_id/:picture_id", middleware.RequireAuth, repository.DeleteReviewPicture)
		review.PUT("/update-review/:review_id", middleware.RequireAuth, repository.EditReview)
		review.DELETE("/delete-review/:review_id", middleware.RequireAuth, repository.DeleteYourReview)
		review.PUT("/like/:review_id", middleware.RequireAuth, repository.LikeReview)
		review.PUT("/dislike/:review_id", middleware.RequireAuth, repository.DislikeReview)

	}
	fmt.Println("......overall-rate.....")
	ovrl := router.Group("/overall-rate")
	{
		ovrl.GET("", repository.OverallRating)
	}

	fmt.Println("...article....")
	article := router.Group("/article")
	{

		article.GET("", repository.GetAllArticles)
		article.GET("/categories", repository.GetArticlesCategory)
		article.GET("/article-by-category/:category_id", repository.GetArticlesByCategory)
		article.GET("/:id", repository.GetArticle)
		article.GET("/interesting", repository.GetInterestingArticles)
		article.PUT("/like-article/:article_id", middleware.RequireAuth, repository.LikeArticle)
		article.PUT("/dislike-article/:article_id", middleware.RequireAuth, repository.DislikeArticle)
		article.POST("/create-commentary/:article_id", middleware.RequireAuth, repository.CreateCommentary)
		article.GET("/comments/:article_id", repository.GetCommentsByArticle)
		article.DELETE("/commentary-delete/:article_id/:commentary_id", middleware.RequireAuth, repository.DeleteCommentary)
		article.PUT("/like-commentary/:commentary_id", middleware.RequireAuth, repository.LikeCommentary)
		article.PUT("/dislike-commentary/:commentary_id", middleware.RequireAuth, repository.DislikeCommentary)

	}

	fmt.Println("...portfolio...")
	portfolio := router.Group("/portfolio")
	{
		portfolio.GET("", repository.GetAllPortfolios)
	}

	fmt.Println("...application...")
	application := router.Group("/application")
	{
		application.POST("/create", middleware.RequireAuth, repository.CrateApplication)
		application.GET("/user-applications", middleware.RequireAuth, repository.GetAllUserApplications)
		application.DELETE("/cancel-application/:id", middleware.RequireAuth, repository.CancelApplication)

	}

	fmt.Println("....admin-panel...")
	admin := router.Group("/admin-panel")
	{
		admin.DELETE("/delete-review/:id", middleware.RequireAuth, repository.DeleteReview)
		admin.PUT("/ban-user/:id", middleware.RequireAuth, repository.BanUser)
		admin.PUT("/unban-user/:id", middleware.RequireAuth, repository.UnbanUser)
		admin.GET("/banned", middleware.RequireAuth, repository.GetBannedUsers)
		admin.POST("/create-article", middleware.RequireAuth, repository.CreateArticle)
		admin.POST("/create-article-category", middleware.RequireAuth, repository.CreateArticlesCategory)
		admin.DELETE("/delete-article-category/:id", middleware.RequireAuth, repository.DeleteArticlesCategory)
		admin.PUT("/edit-article/:id", middleware.RequireAuth, repository.EditArticle)
		admin.DELETE("/delete-article-picture/:article_id/:picture_id", middleware.RequireAuth, repository.DeletePictureFromArticle)
		admin.POST("/add-portfolio", middleware.RequireAuth, repository.AddPortfolioWork)
		admin.DELETE("/delete-portfolio/:id", middleware.RequireAuth, repository.DeletePortfolioWork)
		admin.POST("/create-stylization", middleware.RequireAuth, repository.CreateStylization)
		admin.POST("/create-makeup", middleware.RequireAuth, repository.CreateMakeUp)
		admin.GET("/all-applications", middleware.RequireAuth, repository.GetAllApplications)
		admin.GET("/done-applications", middleware.RequireAuth, repository.GetDoneApplications)
		admin.GET("/progress-applications", middleware.RequireAuth, repository.GetInProgressApplications)
		admin.GET("/declined-applications", middleware.RequireAuth, repository.GetDeclinedApplications)
		admin.GET("/new-applications", middleware.RequireAuth, repository.GetNewApplications)
		admin.PUT("/set-application-progress/:id", middleware.RequireAuth, repository.SetIsInProgress)
		admin.PUT("/set-application-decline/:id", middleware.RequireAuth, repository.SetIsDeclined)
		admin.PUT("/set-application-done/:id", middleware.RequireAuth, repository.SetIsDone)
		admin.DELETE("delete-stylization/:id", middleware.RequireAuth, repository.DeleteStylizationCategory)
		admin.DELETE("delete-makeup/:id", middleware.RequireAuth, repository.DeleteMakeupCategory)
	}

	return router

}
