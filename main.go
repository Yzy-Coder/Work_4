package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

//用户注册函数
func UserRegistration(c *gin.Context) {

	//初始化 ， 获得数据库句柄DB, 注册需要的参数:UserName , UserPassWord
	DB := DataBaseInit()
	UserName := c.PostForm("UserName")
	UserPassWord := c.PostForm("UserPassWord")

	//执行注册操作获得UserId
	UserId := Registration(DB, UserName, UserPassWord)

	//根据UserId的返回值来判断注册是否成功 ，-1表示失败 ，其他则成功
	if UserId == -1 {

		c.JSON(http.StatusOK, gin.H{

			"Erro":         "注册失败",
			"UserName":     UserName,
			"UserPassWord": UserPassWord,
		})
	} else {

		c.JSON(http.StatusOK, gin.H{

			"UserId":       UserId,
			"UserName":     UserName,
			"UserPassWord": UserPassWord,
		})
	}

}

//用户登录函数
func UserSignIn(c *gin.Context) {

	//初始化 ， 获得数据库句柄DB, 登录需要的参数:UserName , UserPassWord
	DB := DataBaseInit()
	userName := c.PostForm("UserName")
	userPassWord := c.PostForm("UserPassWord")

	//执行登录操作获得resulet参数
	resulet := SignIn(DB, userName, userPassWord)

	//利用Switch语句进行分类 ， 1则表示成功 ， -1表示数据库操作出错 ， -2表示账号密码不正确
	switch resulet {
	case 1:
		c.JSON(http.StatusOK, gin.H{

			"Login information": "Success",
			"YourName":          userName,
		})
	case -1:
		c.JSON(http.StatusOK, gin.H{

			"Login information": "Fail",
			"YourName":          userName,
			"ErrMessage":        "服务器端出错",
		})
	case -2:
		c.JSON(http.StatusOK, gin.H{

			"Login information": "Fail",
			"YourName":          userName,
			"ErrMessage":        "账号或密码错误",
		})

	}

}

//发布博客函数
func BlogPublish(c *gin.Context) {

	//初始化 ， 获得数据库句柄DB, 发表博客需要的参数:UserId , UserName ,BlogContent
	DB := DataBaseInit()
	UserId := c.PostForm("UserId")
	BlogName := c.PostForm("BlogName")
	BlogContent := c.PostForm("BlogContent")

	//利用博客发布操作获得BlogId
	BlogId := Publish(DB, UserId, BlogName, BlogContent)

	//如果返回值BlogId == -1 ，说明发布失败 ， 否则则成功
	if BlogId == -1 {
		c.JSON(http.StatusOK, gin.H{

			"PublishInformation": "Fail",
			"BlogId":             "Null",
			"BlogName":           BlogName,
			"BlogContent":        BlogContent,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{

			"PublishInformation": "Success",
			"BlogId":             BlogId,
			"BlogName":           BlogName,
			"BlogContent":        BlogContent,
		})
	}

}

//点赞操作函数
func BlogThumbsUp(c *gin.Context) {

	//初始化 ， 获得数据库句柄DB, 注册需要的参数:UserId , BlogId
	DB := DataBaseInit()
	UserId := c.PostForm("UserId")
	BlogId := c.PostForm("BlogId")

	//执行点赞操作获得该篇博客现在的点赞数
	BlogThumbsUpNum := ThumbsUp(DB, UserId, BlogId)

	//如果返回值BlogThumbsUpNum == -1 ，说明点赞失败 ， 否则成功
	if BlogThumbsUpNum == -1 {

		c.JSON(http.StatusOK, gin.H{

			"ThumbsUpInformation": "Fail",
			"UserId":              UserId,
			"BlogId":              BlogId,
			"BlogThumbsUpNum":     "NULL",
		})
	} else {

		c.JSON(http.StatusOK, gin.H{

			"ThumbsUpInformation": "Success",
			"UserId":              UserId,
			"BlogId":              BlogId,
			"BlogThumbsUpNum":     BlogThumbsUpNum,
		})
	}
}

func main() {

	Route := gin.Default()

	Route.POST("/register", UserRegistration)
	Route.POST("/signin", UserSignIn)
	Route.POST("/publish", BlogPublish)
	Route.POST("/thumbsup", BlogThumbsUp)

	Route.Run(":9090")
}
