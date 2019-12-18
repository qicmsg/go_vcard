package user

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"strings"
	"vcard/app/models/entity"
)

type UserCard struct {
	ID      string `gorm:"type:varchar(32);primary_key;not null" json:"id"`
	UserID  int64  `gorm:"column:userid;not null" json:"userid"`
	Name    string `gorm:"type:varchar(32);column:name" json:"name"`
	Phone   string `gorm:"type:varchar(32);column:phone" json:"phone"`
	Sex     int    `gorm:"type:tinyint(4);column:sex" json:"sex"`
	Company string `gorm:"type:varchar(256);column:company" json:"company"`

	Province string `gorm:"type:varchar(32);column:province" json:"province"`
	City     string `gorm:"type:varchar(32);column:city" json:"city"`
	District string `gorm:"type:varchar(32);column:district" json:"district"`
	Height   int    `gorm:"type:int;column:height" json:"height"`
	Bust     int    `gorm:"type:int;column:bust" json:"bust"`
	Waist    int    `gorm:"type:int;column:waist" json:"waist"`
	Hipline  int    `gorm:"type:int;column:hipline" json:"hipline"`

	Cover   string `gorm:"type:varchar(256);column:cover" json:"cover"`
	Profile string `gorm:"type:varchar(512);column:profile" json:"profile"`
	Video   string `gorm:"type:varchar(128);column:video" json:"video"`

	Isavailable bool `gorm:"type:tinyint;column:isavailable" json:"isavailable,omitempty"`

	CreatedAt *entity.Time `json:"created_at,omitempty"`
	UpdatedAt *entity.Time `json:"updated_at,omitempty"`
}

func (UserCard) TableName() string {
	return "user_cards"
}

func (card *UserCard) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("id", strings.ReplaceAll(uuid.NewV5(uuid.NewV1(), card.Name).String(), "-", ""))
	return nil
}
