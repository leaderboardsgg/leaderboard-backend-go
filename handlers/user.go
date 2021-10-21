package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/model"
	"github.com/speedrun-website/leaderboard-backend/utils"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID       uint
	Username string
}

func GetUser(c *gin.Context) {
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

func RegisterUser(c *gin.Context) {
	var registerValue model.UserRegister

	if err := c.BindJSON(&registerValue); err != nil {
		log.Println("Unable to bind value", err)
		return
	}

	hash, err := utils.HashAndSalt([]byte(registerValue.Password))

	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username: registerValue.Username,
		Email:    registerValue.Email,
		Password: hash,
	}

	err = database.DB.WithContext(c).Create(&user).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"errors": [1]gin.H{{"constraint": pgErr.ConstraintName, "message": pgErr.Detail}},
			})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/users/%d", user.ID))
	c.JSON(http.StatusCreated, gin.H{
		"data": &UserResponse{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}
