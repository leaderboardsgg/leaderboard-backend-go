package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"speedrun.website/database"
	"speedrun.website/graph/model"
)

func UsersHandler(c *gin.Context) {
	db, err := database.GetDatabase()

	if err != nil {
		log.Println("Unable to connect to database", err)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer db.Close()

	var users []model.User
	db.Model(model.User{}).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
