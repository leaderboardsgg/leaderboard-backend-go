package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"gorm.io/gorm"
)

func MeHandler(c *gin.Context) {
	user, _ := c.Get(middleware.JwtConfig.IdentityKey)
	var me model.User

	result := database.DB.Where(model.User{
		Email: user.(*model.User).Email,
	}).First(&me)

	if result.Error != nil {
		var code = http.StatusInternalServerError

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}

		c.AbortWithStatusJSON(code, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, me)
}
