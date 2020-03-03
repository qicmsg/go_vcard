package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"time"
	"vcard/pkg/setting"
)

var DB *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key;auto_increment" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

type ModelID struct {
	ID int `gorm:"primary_key;auto_increment" json:"id"`
}

///region 日期转json代码
type Time struct {
	time.Time
}

// MarshalJSON 序列化为JSON
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("\"\""), nil
	}
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// UnmarshalJSON 反序列化为JSON
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	t.Time, err = time.Parse("2006-01-02 15:04:05", string(data)[1:20])
	return err
}

// String 重写String方法
func (t *Time) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// SetRaw 读取数据库值
func (t *Time) SetRaw(value interface{}) error {
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	}
	return nil
}

// RawValue 写入数据库
func (t *Time) RawValue() interface{} {
	str := t.Format("2006-01-02 15:04:05")
	if str == "0001-01-01 00:00:00" {
		return nil
	}
	return str
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

///endregion

///region setup与一些钩子方法
// Setup initializes the database instance
func Setup() {
	var err error
	DB, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	DB.SingularTable(true)
	DB.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DB.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	DB.Callback().Delete().Replace("gorm:delete", deleteCallback)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	DB.LogMode(setting.ServerSetting.RunMode == "debug")
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer DB.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

///endregion

///region 业务上扩展方法
// BuildWhere 构建where条件
func BuildWhere(db *gorm.DB, where interface{}) (*gorm.DB, error) {
	var err error
	t := reflect.TypeOf(where).Kind()
	if t == reflect.Struct || t == reflect.Map {
		db = db.Where(where)
	} else if t == reflect.Slice {
		for _, item := range where.([]interface{}) {
			item := item.([]interface{})
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("切片长度不能小于2")
				}
				columnstr := column.(string)
				// 拼接参数形式
				if strings.Index(columnstr, "?") > -1 {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and" //cond
					opt := "="
					_opt := " = "
					var val interface{}
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						_opt = " " + strings.ReplaceAll(opt, " ", "") + " "
						val = item[2]
					}

					if count == 4 {
						cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
					}

					/*
					   '=', '<', '>', '<=', '>=', '<>', '!=', '<=>',
					   'like', 'like binary', 'not like', 'ilike',
					   '&', '|', '^', '<<', '>>',
					   'rlike', 'regexp', 'not regexp',
					   '~', '~*', '!~', '!~*', 'similar to',
					   'not similar to', 'not ilike', '~~*', '!~~*',
					*/

					if strings.Index(" in notin ", _opt) > -1 {
						// val 是数组类型
						column = columnstr + " " + opt + " (?)"
					} else if strings.Index(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp", _opt) > -1 {
						column = columnstr + " " + opt + " ?"
					}

					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map /*Map*/ {
				db = db.Where(item)
			} else {
				/*
				// 解决and 与 or 混合查询，但这种写法有问题，会抛出 invalid query condition
				db = db.Where(func(db *gorm.DB) *gorm.DB {
					db, err = BuildWhere(db, item)
					if err != nil {
						panic(err)
					}
					return db
				})*/

				db, err = BuildWhere(db, item)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		return nil, errors.New("参数有误")
	}
	return db, nil
}

func BuildQueryList(db *gorm.DB, wheres interface{}, columns interface{}, orderBy interface{}, page, rows int) (*gorm.DB, error) {
	var err error
	db, err = BuildWhere(db, wheres)
	if err != nil {
		return nil, err
	}
	db = db.Select(columns)
	if orderBy != nil && orderBy != "" {
		db = db.Order(orderBy)
	}
	if page > 0 && rows > 0 {
		db = db.Limit(rows).Offset((page - 1) * rows)
	}
	return db, err
}

///endregion
