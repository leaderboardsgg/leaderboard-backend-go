package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/model"
	"github.com/speedrun-website/leaderboard-backend/utils"
	"gorm.io/gorm"
)

type RegisterResponse struct {
	ID       uint
	Username string
}

func RegisterHandler(c *gin.Context) {
	var registerValue model.UserRegister
	var alreadyExist model.User
	var result *gorm.DB

	if err := c.Bind(&registerValue); err != nil {
		log.Println("Unable to bind value", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result = database.DB.Where(model.User{Email: registerValue.Email}).Find(&alreadyExist)

	if result.Error != nil {
		log.Fatal(result.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if result.RowsAffected != 0 {
		// warning: maybe return a 201 instead for security reason?
		// more: https://stackoverflow.com/a/53144807/2816588
		c.JSON(http.StatusConflict, gin.H{
			"message": "An account with this email already exists",
		})
		return
	}

	hash, err := utils.HashAndSalt([]byte(registerValue.Password))

	if err != nil {
		log.Fatal(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username: registerValue.Username,
		Email:    registerValue.Email,
		Password: hash,
	}

	result = database.DB.Create(&user)

	if result.Error != nil {
		log.Fatal(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/users/%d", user.ID))
	c.JSON(http.StatusCreated, gin.H{
		"data": &RegisterResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}
