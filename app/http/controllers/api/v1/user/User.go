package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vcard/app/models/entity"
	"vcard/app/models/user"
	userService "vcard/app/services/user"
	"vcard/pkg/utils"
	"vcard/pkg/utils/result"
)

// info 用户信息
func Info(ctx *gin.Context) {
	var userinfo user.User

	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		id = 0
	}

	//userService.User = userService.UserStruct{Cache: false}

	bll := userService.UserStruct{Cache: false}

	fmt.Printf("--------:%v", &bll)
	var model, _ = bll.GetById(id)

	// 主要测试json序列化与反序列化
	keyBytes, err := json.Marshal(model)
	fmt.Println(string(keyBytes))
	json.Unmarshal(keyBytes, &userinfo)

	ctx.JSON(http.StatusOK, utils.Result(result.OK, userinfo, ""))
}

func List(ctx *gin.Context) {
	total := 0

	bll := userService.UserStruct{Cache: false}

	where := []interface{}{
		[]interface{}{"id", "in", []int{1, 2}},
	}

	list, _ := bll.List(where, []string{"*"}, "id desc", 2, 1, &total)

	ctx.JSON(http.StatusOK, utils.Result(result.OK, map[string]interface{}{
		"list": list, "total": total,
	}, ""))
}

func Save(ctx *gin.Context) {
	var card user.UserCard
	card.UserID = 2
	card.Name = "yonds"
	card.Phone = "17611168388"
	card.Sex = 1

	card.Isavailable = true
	//card.CreatedAt = &models.Time{Time: time.Now()}

	entity.DB.Create(&card)

	ctx.JSON(http.StatusOK, card)
}

func Test(ctx *gin.Context) {
	db := entity.DB
	var err error
	var list []*user.User

	/*
		// 测试1 等于
		where := []interface{}{
			[]interface{}{"id", "=", 1},
			[]interface{}{"username", "chen"},
		}
		// SELECT * FROM `users`  WHERE (id = 1) and (username = 'chen')
	*/

	/*
		// 测试2 结构体
		where := user.User{ID: 1, UserName: "chen"}
		// SELECT * FROM `users`  WHERE (id = 1) and (username = 'chen')
	*/

	/*
		// 测试3 in 与 Or
		where := []interface{}{
			[]interface{}{"id", "in", []int{1, 2}},
			[]interface{}{"username", "=", "chen", "or"},
		}
		// SELECT * FROM `users`  WHERE (id in ('1','2')) OR (username = 'chen')
	*/

	/*
		// 测试4 map
		where := map[string]interface{}{"id": 1, "username": "chen"}
		// SELECT * FROM `users`  WHERE (`users`.`id` = '1') AND (`users`.`username` = 'chen')
	*/
	/*
		// 测试5 where (a=1) and (b=1 or c=1) 形式，之种写法不正确
		where := []interface{}{
			[]interface{}{"id", "in", []int{1, 2}},
			[]interface{}{
				[]interface{}{"username", "=", "chen"},
				[]interface{}{"username", "=", "yond", "or"},
			},
		}
		// !!! SELECT * FROM `users`  WHERE (id in ('1','2')) AND (username = 'chen') OR (username = 'yond')
		// !!! 这个与我们设想的不一样，gorm底层不会对第二个组合元素分组
		// !!! 设想是 SELECT * FROM `users`  WHERE (id in ('1','2')) AND ((username = 'chen') OR (username = 'yond'))
	*/
	/*
		// 测试5.1 where (a=1) and (b=1 or c=1) 形式
		where := []interface{}{
			[]interface{}{"id", "in", []int{1, 2}},
			[]interface{}{"username = ? or nickname = ?", "chen", "yond"},
		}
		// SELECT * FROM `users`  WHERE (id in ('1','2')) AND (username = 'chen' or nickname = 'yond')
	*/
	/*
		// 测试6 not in 与 Or
		where := []interface{}{
			[]interface{}{"id", "not in", []int{1}},
			[]interface{}{"username", "=", "chen", "or"},
		}
		// SELECT * FROM `users`  WHERE (id not in ('1')) OR (username = 'chen')
	*/
	where := []interface{}{
		[]interface{}{"id", "in", []int{1, 2}},
	}
	//db = db.Table("users")

	db, err = entity.BuildWhere(db, where)
	if err != nil {
		fmt.Println(err)
	}

	//db.Limit(2).Offset((2 - 1) * 2).Find(&list)
	//count := 0
	//db.Count(&count)

	db.Preload("UserCard").Find(&list)
	ctx.JSON(http.StatusOK, list)
}
