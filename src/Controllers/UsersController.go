package Controllers

import (
	"errors"
	"strconv"
	"time"

    "github.com/gin-gonic/gin"
	"github.com/KuangjuX/Lang-Huan-Blessed-Land/Models"
	"github.com/KuangjuX/Lang-Huan-Blessed-Land/Help/json"
	"github.com/KuangjuX/Lang-Huan-Blessed-Land/Help/auth"
)

// User register
func Register(c *gin.Context)  {
    username := c.PostForm("username")
    password := c.PostForm("password")
	email    := c.PostForm("email")
	
	if len(username) == 0 || len(password) == 0 || len(email) == 0 {
		json.JsonMsgWithError(c, "Fail to register.", errors.New("Too few arguments."))
		return
	}

	exist, err := Models.IsExistUser(username)
	if err != nil{
		json.JsonMsgWithError(c, "Fail to register", err)
		return
	}

	if exist{
		json.JsonMsgWithSuccess(c, "Username has exist. Please enter a new username.")
		return
	}

    message, err := Models.CreatUser(username, password, email)

	
	if err != nil {
		json.JsonMsgWithError(c, message, err)
		return
	}

	json.JsonMsgWithSuccess(c, message)
}


func LoginByUsername(c *gin.Context){
	var user Models.User
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	data, err := user.Login()
	if err == nil{
		json.JsonDataWithSuccess(c, data)
	}else{
		json.JsonError(c, err)
	}
}


func LoginByEmail(c *gin.Context)  {
	var user Models.User
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	data, err := user.Login()
	if err == nil{
		json.JsonDataWithSuccess(c, data)
	}else{
		json.JsonError(c, err)
	}
}


func ModifyNickname(c *gin.Context){
	new_nickname := c.PostForm("new_nickname")
	user, _ := auth.GetUserByToken(c)
	user_id := user.(Models.User).ID

	err := Models.ModifyNickname(user_id, new_nickname)
	if err != nil{
		json.JsonError(c, err)
		return
	}
	json.JsonSuccess(c)

}

func ModifyEmail(c *gin.Context){
	new_email := c.PostForm("new_email")
	user, _ := auth.GetUserByToken(c)
	user_id := user.(Models.User).ID

	err := Models.ModifyEmail(user_id, new_email)
	if err != nil{
		json.JsonError(c, err)
		return
	}
	json.JsonSuccess(c)
}

func ModifyAvatar(c *gin.Context){
	avatar, _ := c.FormFile("new_avatar")
	avatar_name := strconv.FormatInt(time.Now().Unix(), 10) + "-" + avatar.Filename
	avatar_path := "storage/app/image/" + avatar_name
	

	if err := c.SaveUploadedFile(avatar, avatar_path); err != nil {
		json.JsonMsgWithError(c, "Fail to upload image", err)
		return
	}
	
	user, _ := auth.GetUserByToken(c)
	user_id := user.(Models.User).ID
	if err := Models.ModifyAvatar(user_id, avatar_name); err != nil {
		json.JsonMsgWithError(c, "Fail to store image into database", err)
		return
	}

	json.JsonMsgWithSuccess(c, "Success to upload image")

}

func UploadArticle(c *gin.Context){
	user, _ := auth.GetUserByToken(c)
	user_id := user.(Models.User).ID
	tag_id, _ := strconv.ParseInt(c.PostForm("tag_id"), 10, 64)
	title := c.PostForm("title")
	content := c.PostForm("content")

	if err := Models.UploadArticle(user_id, tag_id, title, content); err != nil {
		json.JsonError(c, err)
		return
	}

	json.JsonSuccess(c)
}

func ModifyArticle(c *gin.Context){
	article_id, _ := strconv.ParseInt(c.PostForm("article_id"), 10, 64)
	tag_id, _ := strconv.ParseInt(c.PostForm("tag_id"), 10, 64)
	title := c.PostForm("title")
	content := c.PostForm("content")

	if err := Models.ModifyArticle(article_id, tag_id, title, content); err != nil {
		json.JsonError(c, err)
		return
	}

	json.JsonSuccess(c)
}

func DeleteArticle(c *gin.Context){
	article_id, _ := strconv.ParseInt(c.PostForm("article_id"), 10, 64)

	if err := Models.DeleteArticle(article_id); err != nil {
		json.JsonError(c, err)
		return
	}

	json.JsonSuccess(c)
}

func LikeArticle(c *gin.Context){
	article_id, _ := strconv.ParseInt(c.PostForm("article_id"), 10, 64)
	user, err := auth.GetUserByToken(c)

	if err != nil{
		json.JsonError(c, err)
	}

	user_id := user.(Models.User).ID

	if err := Models.LikeArticle(article_id, user_id); err != nil{
		json.JsonError(c, err)
		return
	}
	json.JsonSuccess(c)
}

func Follow(c *gin.Context){
	follower_id, _ := strconv.Atoi(c.Param("follower_id"))
	
	data, err := auth.GetUserByToken(c)

	if err != nil{
		json.JsonError(c, err)
	}

	user := data.(Models.User)

	err = user.Follow(follower_id)
	if err != nil{
		json.JsonError(c, err)
	}

	json.JsonSuccess(c)

}

func UnFollow(c *gin.Context){
	follower_id, _ := strconv.Atoi(c.Param("follower_id"))
	
	data, err := auth.GetUserByToken(c)

	if err != nil{
		json.JsonError(c, err)
	}

	user := data.(Models.User)

	is_exist, err := user.IsExistFollow(follower_id)
	if err != nil && is_exist{
		user.UnFollow(follower_id)
		json.JsonSuccess(c)
		return 
	}

	json.JsonMsgWithError(c, "Fail to unfollow", err)
}


func GetFollowers(c *gin.Context){
	data, err := auth.GetUserByToken(c)

	if err != nil{
		json.JsonError(c, err)
	}

	user := data.(Models.User)
	user_id := int(user.ID)

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	followers, err := Models.GetFollowersByPage(page, size, user_id)
	if err != nil{
		json.JsonError(c, err)
	}
	json.JsonDataWithSuccess(c, followers)
}