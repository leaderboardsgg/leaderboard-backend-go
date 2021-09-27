package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"gorm.io/gorm"
)

func UserHandler(c *gin.Context) {
	// Maybe we shouldn't use the increment ID but generate a UUID instead to avoid
	// exposing the amount of users registered in the database.
	var user model.User
	ID := c.Param("id")
	db, err := database.GetDatabase()

	if err != nil {
		log.Println("Unable to connect to database", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := db.First(&user, ID)

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
		"ID":       user.ID,
		"Username": user.Username,
	})
}
