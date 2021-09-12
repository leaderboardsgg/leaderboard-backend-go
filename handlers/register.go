package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"github.com/speedrun-website/leaderboard-backend/utils"
)

func RegisterHandler(c *gin.Context) {
	var registerValue model.Register
	c.Bind(&registerValue)

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

	var alreadyExist model.User
	result := db.Where(model.User{Email: registerValue.Email}).Find(&alreadyExist)

	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "An account with this email already exists",
		})
		return
	}

	hash, _ := utils.HashAndSalt([]byte(registerValue.Password))

	db.Create(model.User{
		Username: registerValue.Username,
		Email:    registerValue.Email,
		Password: hash,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
