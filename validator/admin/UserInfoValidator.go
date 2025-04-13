package admin

import (
	"github.com/go-playground/validator/v10"
	"github.com/lijie-keith/go_init_project/entity"
)

func UserNameUniqueValidator(fl validator.FieldLevel) bool {
	if s, ok := fl.Field().Interface().(string); ok {
		var admin = entity.UserInfo{}
		admin.GetUserInfoByName(s)
		if admin.Id != 0 {
			return false
		}
	}
	return true
}
