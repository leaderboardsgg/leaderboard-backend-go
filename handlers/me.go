package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/model"
)

func MeHandler(c *gin.Context) {
	rawUser, exists := c.Get(middleware.JwtConfig.IdentityKey)
	user := *rawUser.(*model.UserPersonal)

	if exists {
		err := database.DB.Model(&model.User{}).First(&user, user.ID).Error

		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data": &user,
			})
			return
		}
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
