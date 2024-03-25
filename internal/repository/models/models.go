package models

import (
	"time"
)

// Общее

type User struct {
	Id           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `json:"name"`
	Surname      string        `json:"surname"`
	Email        string        `gorm:"unique;" json:"email"`
	Phone        string        `gorm:"unique;" json:"phone" validate:"required"`
	Password     string        `json:"password"`
	DateOfBirth  string        `json:"date_of_birth"`
	IsBanned     bool          `json:"is_banned" `
	IsAdmin      bool          `json:"is_admin"`
	Picture      string        `json:"picture,omitempty"`
	Description  string        `json:"description"`
	WhyIsBlocked string        `json:"whyIsBlocked,omitempty"`
	Applications []Application `gorm:"many2many:application_users;" json:"applications,omitempty"`
}

type OverallRating struct {
	Id        uint    `gorm:"primaryKey" json:"id"`
	Overall   float64 `json:"overall,omitempty"`
	FiveStar  float64 `json:"five_star,omitempty"`
	FourStar  float64 `json:"four_star,omitempty"`
	ThirdStar float64 `json:"third_star,omitempty"`
	TwoStar   float64 `json:"two_star,omitempty"`
	OneStar   float64 `json:"one_star,omitempty"`
}

type Review struct {
	Id             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `json:"user_id" gorm:"index"`
	User           User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Rate           int32      `json:"rate"`
	ReviewText     string     `json:"review_text"`
	Likes          []Like     `gorm:"many2many:reviews_likes;" json:"review_likes,omitempty"`
	Dislikes       []Dislike  `gorm:"many2many:reviews_dislikes;" json:"review_dislikes,omitempty"`
	LikeCount      int64      `json:"like_count,omitempty"`
	DislikeCount   int64      `json:"dislike_count,omitempty"`
	ReviewPictures []Pictures `gorm:"many2many:reviews_pictures;" json:"review_pictures,omitempty"`
}

type Pictures struct {
	Id      uint64 `gorm:"primaryKey" json:"id"`
	Picture string `json:"picture"`
}

type Portfolio struct {
	Id      uint   `gorm:"primaryKey" json:"id"`
	Picture string `json:"picture,omitempty"`
}

type Like struct {
	Id     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `json:"user_id" gorm:"index"`
	User   User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Liked  int64 `json:"liked,omitempty"`
}
type Dislike struct {
	Id       uint  `gorm:"primaryKey" json:"id"`
	UserID   uint  `json:"user_id" gorm:"index"`
	User     User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Disliked int64 `json:"disliked,omitempty"`
}

// Статья

type ArticlesCategory struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Id                uint                `gorm:"primaryKey" json:"id"`
	Title             string              `json:"title"`
	ACategoryID       uint                `json:"a_category_id"`
	ArticlesCategory  ArticlesCategory    `gorm:"foreignKey:ACategoryID" json:"article_category,omitempty"`
	Description       string              `json:"description"`
	Content           string              `json:"content"`
	ArticleCommentary []ArticleCommentary `gorm:"many2many:article_comments;" json:"comments,omitempty"`
	Likes             []Like              `gorm:"many2many:articles_likes;" json:"article_likes,omitempty"`
	Dislikes          []Dislike           `gorm:"many2many:articles_dislikes;" json:"article_dislikes,omitempty"`
	LikeCount         int64               `json:"like_count,omitempty"`
	DislikeCount      int64               `json:"dislike_count,omitempty"`
	ArticlePictures   []Pictures          `gorm:"many2many:articles_pictures;" json:"article_pictures,omitempty"`
}

type ArticleCommentary struct {
	Id           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id" gorm:"index"`
	User         User      `gorm:"foreignKey:UserID" json:"users,omitempty"`
	Content      string    `json:"content"`
	ArticleID    uint      `json:"article_id"`
	Article      Article   `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Likes        []Like    `gorm:"many2many:commentaries_likes;" json:"commentary_likes,omitempty"`
	Dislikes     []Dislike `gorm:"many2many:commentaries_dislikes;" json:"commentary_dislikes,omitempty"`
	LikeCount    int64     `json:"like_count,omitempty"`
	DislikeCount int64     `json:"dislike_count,omitempty"`
}

// Заявка

type Makeup struct {
	Id   uint   `gorm:"PrimaryKey" json:"id"`
	Name string `json:"name"`
}

type Stylization struct {
	Id   uint   `gorm:"PrimaryKey" json:"id"`
	Name string `json:"name"`
}

type Application struct {
	Id                uint        `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time   `json:"created_at"`
	UserID            uint        `json:"user_id" gorm:"index"`
	User              User        `gorm:"foreignKey:UserID" json:"users,omitempty"`
	StylizationID     uint        `json:"stylization_id" gorm:"index"`
	Stylization       Stylization `gorm:"foreignKey:StylizationID " json:"stylization,omitempty"`
	MakeupID          uint        `json:"makeup_id" gorm:"index"`
	Makeup            Makeup      `gorm:"foreignKey:MakeupID " json:"makeup,omitempty"`
	Day               string      `json:"day"`
	Time              string      `json:"time"`
	Preferences       string      `json:"preferences"`
	ReferencePictures []Pictures  `gorm:"many2many:references_pictures;" json:"reference_pictures,omitempty"`
	FacePictures      []Pictures  `gorm:"many2many:face_pictures;" json:"face_pictures,omitempty"`
	IsDeclined        bool        `json:"is_declined"`
	IsInProgress      bool        `json:"is_in_progress"`
	Done              bool        `json:"done"`
}
