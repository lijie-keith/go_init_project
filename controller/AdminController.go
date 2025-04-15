package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lijie-keith/go_init_project/commonUtils"
	"github.com/lijie-keith/go_init_project/entity"
	"net/http"
	"strconv"
)

func CreateUserInfo(c *gin.Context) {
	var userInfo entity.UserInfo
	if err := c.ShouldBind(&userInfo); err == nil {
		msg := userInfo.CreateUserInfo()
		c.JSON(http.StatusOK, commonUtils.OK.WithData(msg))
		return
	} else {
		println(err)
		c.JSON(http.StatusOK, commonUtils.ErrParam.WithMsg(err.Error()))
		return
	}
}

func DeleteUserInfoById(c *gin.Context) {
	var userInfo entity.UserInfo
	idStr := c.Param("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, commonUtils.ErrIdBlank)
		return
	}
	id, _ := strconv.Atoi(idStr)

	userInfo.GetUserInfoById(id)
	if userInfo.Id == 0 {
		c.JSON(http.StatusOK, commonUtils.ErrDataNoExist)
		return
	}

	userInfo.IsDeleted = 1
	userInfo.UpdateUserInfoById()
	c.JSON(http.StatusOK, commonUtils.OK)
}

func UpdateUserInfoById(c *gin.Context) {
	var userInfo entity.UserInfo
	var temp entity.UserInfo
	if err := c.ShouldBind(&userInfo); err == nil {
		if userInfo.Id == 0 {
			c.JSON(http.StatusOK, commonUtils.ErrIdBlank)
			return
		}
		temp.GetUserInfoById(userInfo.Id)
		if temp.Id == 0 {
			c.JSON(http.StatusOK, commonUtils.ErrDataNoExist)
			return
		}

		userInfo.UpdateUserInfoById()
		c.JSON(http.StatusOK, commonUtils.OK)
	} else {
		println(err)
		c.JSON(http.StatusOK, commonUtils.ErrParam.WithMsg(err.Error()))
		return
	}
}

func GetUserInfoById(c *gin.Context) {
	var userInfo = new(entity.UserInfo)
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusOK, commonUtils.ErrIdBlank)
		return
	}
	id, _ := strconv.Atoi(idStr)

	userInfo.GetUserInfoById(id)
	c.JSON(http.StatusOK, commonUtils.OK.WithData(userInfo))
	return
}

func PageUserInfo(c *gin.Context) {
	var userInfo = new(entity.UserInfo)
	var userInfoPage = new(entity.UserInfoPage)
	var userInfoList []entity.UserInfo

	if err := c.ShouldBind(&userInfoPage); err == nil {

		userInfoList = userInfo.PageUserInfo(userInfoPage)

		c.JSON(http.StatusOK, commonUtils.OK.WithData(userInfoList))
		return
	} else {
		println(err)
		c.JSON(http.StatusOK, commonUtils.ErrParam.WithMsg(err.Error()))
		return
	}
}
