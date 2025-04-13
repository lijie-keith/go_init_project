package entity

type BaseModel struct {
	Id         int `gorm:"primaryKey;autoIncrement;not null;comment:'主键id'" json:"id"`
	CreateTime int `gorm:"autoCreateTime;comment:'创建时间'"`
	UpdateTime int `gorm:"autoUpdateTime;comment:'更新时间'"`
	DeleteTime int `gorm:"autoUpdateTime;comment:'删除时间'"`
	IsDeleted  int `gorm:"default:0;comment:'是否删除0-未删除1-已删除'"`
}
