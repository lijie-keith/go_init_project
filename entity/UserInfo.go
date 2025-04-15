package entity

import (
	"fmt"
	"github.com/lijie-keith/go_init_project/commonUtils"
)

type UserInfo struct {
	BaseModel
	UserName string `gorm:"type:varchar(100);comment:'用户名称'" json:"username"`
	Age      int    `gorm:"type:int;comment:'年龄'" json:"age"`
	Sex      int    `gorm:"type:int;comment:'性别0-男1-女'" json:"sex"`
	Phone    string `gorm:"type:varchar(15);comment:'电话'" json:"phone"`
}

type UserInfoPage struct {
	BasePageInfo
	UserInfo
}

func init() {
	err := commonUtils.DB.AutoMigrate(&UserInfo{})
	if err != nil {
		fmt.Println("------------------- admin migrate err -------------------")
		panic(err)
	}
}

func (*UserInfo) TableName() string {
	return "user_info"
}

func (userInfo *UserInfo) CreateUserInfo() int {
	commonUtils.DB.Create(userInfo)
	return userInfo.Id
}

func (userInfo *UserInfo) GetUserInfoById(id int) {
	commonUtils.DB.Where("is_deleted = 0").First(&userInfo, id)
}

func (userInfo *UserInfo) GetUserInfoByName(userName string) {
	commonUtils.DB.Model(&UserInfo{}).Limit(1).Where("is_deleted = 0").Where(UserInfo{UserName: userName}).First(userInfo)
}

func (userInfo *UserInfo) UpdateUserInfoById() {
	commonUtils.DB.Save(&userInfo)
}

func (userInfo *UserInfo) PageUserInfo(userInfoPage *UserInfoPage) []UserInfo {
	var users []UserInfo

	// 假设我们要查询第2页，每页10条记录
	current := userInfoPage.Current
	size := userInfoPage.Size

	// 计算 Offset
	offset := (current - 1) * size

	// 执行分页查询
	commonUtils.DB.Limit(size).Offset(offset).Where("is_deleted = 0").Find(&users)
	return users
}
