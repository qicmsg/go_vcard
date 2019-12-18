package user

import (
	"vcard/app/models/entity"
)

type User struct {
	// models.GormModelID
	ID            int          `gorm:"primary_key;auto_increment" json:"id"`
	UserName      string       `gorm:"type:varchar(32);column:username;not null" json:"username"`
	NickName      string       `gorm:"type:varchar(32);column:nickname" json:"nickname"`
	Password      string       `gorm:"type:varchar(255);column:password" json:"-"`
	Avatar        string       `gorm:"type:varchar(512);column:avatar" json:"avatar"`
	LastLoginTime *entity.Time `gorm:"type:datetime;column:lastlogintime" json:"lastlogintime,omitempty"`
	CreateTime    *entity.Time `gorm:"type:datetime;column:createtime" json:"createtime,omitempty"`
	UserCard      []UserCard   `gorm:"foreignkey:UserID" json:"user_card"`
}

func (User) TableName() string {
	return "users"
}
