package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lijie-keith/go_init_project/controller"
	"github.com/lijie-keith/go_init_project/validator/admin"
)

func InitRouter(router *gin.Engine) {
	adminRouter := router.Group("/admin")
	{
		adminRouter.POST("/create", controller.CreateUserInfo)
		adminRouter.GET("/getById/:id", controller.GetUserInfoById)
		adminRouter.DELETE("/deleteById/:id", controller.DeleteUserInfoById)
		adminRouter.PUT("/update", controller.UpdateUserInfoById)
		adminRouter.POST("/page", controller.PageUserInfo)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("UserNameUnique", admin.UserNameUniqueValidator)
		if err != nil {
			return
		}
	}
}
