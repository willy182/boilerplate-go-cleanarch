package model

import "time"

// GormArticle data of struct
type GormArticle struct {
	ID          int        `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Title       string     `gorm:"type:varchar(100);NOT NULL"`
	Summary     string     `gorm:"type:varchar(250);NOT NULL"`
	Description string     `gorm:"type:text"`
	Image       string     `gorm:"type:varchar(150)"`
	Created     *time.Time `gorm:"type:timestamp(6) with time zone;NOT NULL"`
	Modified    *time.Time `gorm:"type:timestamp(6) with time zone"`
}

// Article data of struct
type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	Created     time.Time `json:"created"`
	Modified    string    `json:"modified,omitempty"`
}
