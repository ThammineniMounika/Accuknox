package main

import (
	"log"
    "strings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	
)

type User struct {
	gorm.Model
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}
type Login struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
type LoginResponse struct {
	SID uint32 `json:"id"`
}
type Notes struct {
	gorm.Model
	UserID   uint32 `form:"id"`
	Note string `form:"note"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

var db *gorm.DB

func main() {
	
	var err error
	db, err = gorm.Open("mysql", "root:Mouni@9797@tcp(127.0.0.1:3306)/accuknox?charset=utf8&parseTime=True")
	
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{},&Notes{})

	r := gin.Default()

	r.POST("/signup", CreateUser)
	r.POST("/login", UserLogin)
	r.GET("/notes", ListNotes)
	r.POST("/notes", CreateNote)
	r.DELETE("/notes", DeleteNote)

	if err := r.Run(":8280"); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(c *gin.Context) {
	var user User
	c.Bind(&user)
	db.Where("name=? AND email=? AND password=?",user.Name,user.Email,user.Password).Find(&user)

	if len(user.Name)<3 && len(user.Password)<3{
		c.JSON(400,"bad request")
		return
	}else if user.Name== "" && user.Password == ""{
		c.JSON(201, "invalid name or password")
	}
	if !checkEmailFormat(user.Email) {
		c.JSON(200,  "please make sure email contains @ and .com")
		return
	}

	
    var Response ErrorResponse = ErrorResponse{"user created sussesfully"}
	c.JSON(200, Response)
}
func checkEmailFormat(email string) bool {
	
	return strings.Contains(email, "@") && strings.Contains(email, ".com")
}

func UserLogin(c *gin.Context) {
	
	var request Login
	c.Bind(&request)

	var user User
	result := db.Where("email = ? AND password= ?", request.Email,request.Password).Find(&user)
	if result.Error != nil {
		c.JSON(400, "bad request")
		return
	}

	if user.Password != request.Password && user.Email!= request.Email{
		c.JSON(401,  "unauthorized user")
		return
	}
		userID := uint32(user.ID)
		c.JSON(200,"login successfull ")
	   c.JSON(200,LoginResponse{SID: userID})
}

func ListNotes(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		c.JSON(400,"bad request")
		return
	}

	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		c.JSON(401,  "unauthorized")
		return
	}

	var notes []Notes
	result = db.Where("user_id = ?", user.ID).Find(&notes)
	if result.Error != nil {
		c.JSON(201,  "Failed to retrieve notes")
		return
	}

	c.JSON(200,  notes)
}

func CreateNote(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		c.JSON(400,  "bad request")
		return
	}

	
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		c.JSON(401," unauthorized User")
		return
	}

	var note Notes
	newNote := Notes{
		UserID: uint32(user.ID),
		Note:   note.Note,
	}

	result = db.Create(&newNote)
	if result.Error != nil {
		c.JSON(200, "Failed to create note")
		return
	}

	c.JSON(200, newNote)
}

func DeleteNote(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(200, "User ID is missing")
		return
	}

	noteID := c.Query("note_id")
	if noteID == "" {
		c.JSON(200, "Note ID is missing")
		return
	}

	
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		c.JSON(401,  "unauthorized user")
		return
	}

	result = db.Delete(&Notes{}, "id = ? AND user_id = ?", noteID, user.ID)
	if result.Error != nil {
		c.JSON(400, "Failed to delete note")
		return
	}

	c.JSON(200, "Note deleted successfully")
}
