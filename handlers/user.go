package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/middleware"
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var user model.UserIdentifier
	err = database.DB.Model(&model.User{}).First(&user, id).Error

	if err != nil {
		var code int
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": &user,
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
			/*
			 * TODO: we probably don't want to reveal if an email is already in use.
			 * Maybe just give a 201 and send an email saying that someone tried to sign up as you.
			 * --Ted W
			 */
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
		"data": &model.UserIdentifier{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

func Me(c *gin.Context) {
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
