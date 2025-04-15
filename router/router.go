package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lijie-keith/go_init_project/controller"
	"github.com/lijie-keith/go_init_project/validator/admin"
)

func InitRouter(router *gin.Engine) {
	userInfoRouter := router.Group("/userInfo")
	{
		userInfoRouter.POST("/create", controller.CreateUserInfo)
		userInfoRouter.GET("/getById/:id", controller.GetUserInfoById)
		userInfoRouter.DELETE("/deleteById/:id", controller.DeleteUserInfoById)
		userInfoRouter.PUT("/update", controller.UpdateUserInfoById)
		userInfoRouter.POST("/page", controller.PageUserInfo)
	}

	blockChainRouter := router.Group("/blockChain")
	{
		blockChainRouter.GET("/getAccountBalance", controller.QueryAccountBalance)
		blockChainRouter.GET("/queryTokenBalance", controller.QueryTokenBalance)
		blockChainRouter.POST("/createNewWallet", controller.CreateNewWallet)
		blockChainRouter.GET("/queryLastBlock", controller.QueryLastBlock)
		blockChainRouter.GET("/queryFirstTransactionFromBlock", controller.QueryFirstTransactionFromBlock)
		blockChainRouter.POST("/transferBalance", controller.TransferBalance)
		blockChainRouter.POST("/transferTokenBalance", controller.TransferTokenBalance)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("UserNameUnique", admin.UserNameUniqueValidator)
		if err != nil {
			return
		}
	}
}
