package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
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

	var users []model.User
	db.Model(model.User{}).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}