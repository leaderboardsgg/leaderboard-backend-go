package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	database "speedrun.website/db"
	"speedrun.website/graph/model"
)

func UsersHandler(c *gin.Context) {
	db, err := database.GetDatabase()

	if err != nil {
		log.Println("Unable to connect to database", err)
		c.Error(err)
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer db.Close()

	var users []model.User
	db.Model(model.User{}).Find(&users)

	c.JSON(200, gin.H{
		"data": users,
	})
}
