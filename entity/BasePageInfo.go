package entity

type BasePageInfo struct {
	Current int `default:"1" json:"current"`
	Size    int `default:"10" json:"size"`
}
