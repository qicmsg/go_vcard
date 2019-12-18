package user

import (
	"github.com/jinzhu/gorm"
	"vcard/app/models/entity"
	userModel "vcard/app/models/user"
)

type UserStruct struct {
	Cache bool
}

//type User UserStruct = UserStruct{Cache: true}

// 根据用户id查找用户信息
func (service UserStruct) GetById(id int) (*userModel.User, error) {
	var model userModel.User
	var where userModel.User
	db := entity.DB

	where.ID = id

	// model.LastLoginTime = &models.Time{Time: time.Now()}

	// err := db.Where("id = ?", id).First(&model).Error
	// err := db.First(&model, id).Error
	// err := db.Where(where).Preload("UserCard").First(&model, id).Error

	err := db.Where(where).Preload("UserCard", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).First(&model, id).Error

	/*if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}*/

	//panic("出错了！！")

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (service UserStruct) List(wheres interface{}, columns interface{}, orderBy interface{}, page, rows int, total *int) ([]*userModel.User, error) {
	var err error
	db := entity.DB

	var model []*userModel.User
	var mod userModel.User


	db, err = entity.BuildQueryList(db, wheres, columns, orderBy, page, rows)
	if err != nil {
		return nil, err
	}
	err = db.Preload("UserCard", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).Find(&model).Error

	db, err = entity.BuildWhere(db, wheres)

	db = entity.DB
	db.Model(&mod).Count(total)

	return model, nil
}
