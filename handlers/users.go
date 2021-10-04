package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
)

func UsersHandler(c *gin.Context) {
	var users []model.User
	database.DB.Model(model.User{}).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
