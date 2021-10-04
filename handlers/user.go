package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/model"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID       uint
	Username string
}

func UserHandler(c *gin.Context) {
	// Maybe we shouldn't use the increment ID but generate a UUID instead to avoid
	// exposing the amount of users registered in the database.
	var user model.User
	result := database.DB.First(&user, c.Param("id"))

	if result.Error != nil {
		code := http.StatusInternalServerError

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}

		c.AbortWithStatusJSON(code, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}
