package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

type Account struct {
	Id          uint   `gorm:"primary_key;AUTO_INCREMENT"  json:"-"`
	Player_id   string `json:"#account_id"`
	Type        string `gorm:"-" json:"#user_set"`
	Time        string `gorm:"-" json:"#time"`
	Channel_id  string `gorm:"type:varchar(10)" json:"-"`
	Uid         string `gorm:"index:idx_uid_time;NOT NULL;type:varchar(100)"  json:"-"`
	Prop        *Prop  `json:"properties"`
	Active_time int    `gorm:"index:idx_uid_time;type:int(10);NOT NULL"  json:"-"`
}
type Prop struct {
	Active_time string `json:"active_time"`
	Channel_id  string `json:"channel_id"`
	Uid         string `json:"uid"`
}

func init() {
	fmt.Println("conn")
	Db, _ = gorm.Open("mysql", "root:root@(localhost:3306)/wash?charset=utf8mb4&parseTime=True&loc=Local")
	// if err != nil {
	// 	panic(fmt.Errorf("创建数据库连接失败:%v", err))

	// }
	Db.AutoMigrate(Account{})

	// defer Db.Close()
}
