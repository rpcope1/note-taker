package main
// Borrowing aggressively from: https://gist.github.com/EtienneR/ed522e3d31bc69a9dec3335e639fcf60

import (
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

type Note struct {
	Id int `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Title string `gorm:"not null" form:"title" json:"title"`
        Message string `gorm:"not null" form:"message" json:"message"`
}

func InitDb() *gorm.DB {
       db, err := gorm.Open("sqlite3", "./data.db")

       db.LogMode(true)

       if err != nil {
           panic(err)
       }

       if !db.HasTable(&Note{}) {
           db.CreateTable(&Note{})
       }

       return db
}

func Cors() gin.HandlerFunc {
    return func(c *gin.Context){
        c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
    }
}

func main() {
	r := gin.Default()

    r.Use(Cors())

    v1 := r.Group("api/v1")
    {
	  v1.POST("/notes", PostNote)
          v1.GET("/notes", GetNotes)
    }

    r.Run("127.0.0.1:8080")
}

func PostNote(c *gin.Context) {
    db := InitDb()
    defer db.Close()

    var note Note
    c.Bind(&note)

    if note.Title != "" && note.Message != "" {
         db.Create(&note)
	 c.JSON(201, gin.H{"success": note})
    } else {
         c.JSON(422, gin.H{"error": "Fields are empty"})
    }
}

func GetNotes(c *gin.Context) {
     db := InitDb()
     defer db.Close()

     var notes []Note
     db.Find(&notes)
     c.JSON(200, notes)
}
