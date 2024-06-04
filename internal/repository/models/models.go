package models

import "time"

// Общее

type User struct {
	Id               uint   `gorm:"primaryKey" json:"id"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	Email            string `gorm:"unique;" json:"email"`
	Password         string `json:"password"`
	DateOfBirth      string `json:"date_of_birth"`
	Picture          string `json:"picture,omitempty"`
	Banner           string `json:"banner ,omitempty"`
	Description      string `json:"description"`
	Role             string `json:"role"`
	City             string `json:"city"`
	IsAdmin          bool   `json:"is_admin"`
	IsBlocked        bool   `json:"is_blocked"`
	SubscribersScore int    `json:"subscribers,omitempty"`
}

type Subscriber struct {
	ID         uint  `gorm:"primaryKey" json:"id"`
	UserId     int   `json:"subscriber"`
	User       *User `gorm:"foreignKey:UserId " json:"user ,omitempty"`
	Subscribed bool  `json:"subscribed"`
	CreatorId  int   `json:"creatorId"`
	Creator    *User `gorm:"foreignKey:CreatorId " json:"creator ,omitempty"`
}

type Like struct {
	Id     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `json:"user_id" gorm:"index"`
	User   User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Liked  int64 `json:"liked,omitempty"`
}
type Views struct {
	Id     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `json:"user_id" gorm:"index"`
	User   User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Viewed int64 `json:"viewed"`
}

// Статья

type Article struct {
	Id                uint                `gorm:"primaryKey" json:"id"`
	Title             string              `json:"title"`
	MainPic           string              `json:"main_pic"`
	ChapterOne        string              `json:"chapterOne"`
	ChapterOnePic     string              `json:"chapteronePic"`
	ChapterTwo        string              `json:"chapterTwo"`
	ChapterTwoPic     string              `json:"chaptertwoPic"`
	ChapterThree      string              `json:"chapterThree"`
	ChapterThreePic   string              `json:"chapterthreePic"`
	ArticleCommentary []ArticleCommentary `gorm:"many2many:article_comments;" json:"comments,omitempty"`
	CommentaryCount   int64               `json:"commentCount,omitempty"`
	Likes             []Like              `gorm:"many2many:articles_likes;" json:"article_likes,omitempty"`
	Views             []Views             `gorm:"many2many:articles_views;" json:"articles_views,omitempty"`
	LikeCount         int64               `json:"like_count,omitempty"`
	ViewsCount        int64               `json:"views_count"`
	Blocked           bool                `json:"blocked,omitempty"`
	Published         bool                `json:"published,omitempty"`
	UserId            int                 `json:"creator_id"`
	User              *User               `gorm:"foreignKey:UserId " json:"creator ,omitempty"`
	CreatedAt         time.Time           `json:"created_at,omitempty"`
}

type ArticleCommentary struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id" gorm:"index"`
	User      User      `gorm:"foreignKey:UserID" json:"users,omitempty"`
	Content   string    `json:"content"`
	ArticleID uint      `json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Likes     []Like    `gorm:"many2many:commentaries_likes;" json:"commentary_likes,omitempty"`
	LikeCount int64     `json:"like_count,omitempty"`
	Blocked   bool      `json:"blocked,omitempty"`
	Published bool      `json:"published,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

//Уведомления

type Notification struct {
	Id           uint   `gorm:"primaryKey" json:"id"`
	Notification string `json:"notification"`
	UserID       int    `json:"userid"`
	User         User   `gorm:"foreignKey:UserID" json:"users,omitempty"`
}
