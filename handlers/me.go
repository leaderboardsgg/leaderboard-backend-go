package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/model"
	"gorm.io/gorm"
)

func MeHandler(c *gin.Context) {
	user, _ := c.Get(middleware.JwtConfig.IdentityKey)
	var me model.User

	result := database.DB.Where(model.User{
		Email: user.(*model.User).Email,
	}).First(&me)

	if result.Error != nil {
		var code int
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(code, gin.H{
			"errors": [1]gin.H{{"message": result.Error.Error()}},
		})
		return
	}

	c.JSON(http.StatusOK, me)
}
