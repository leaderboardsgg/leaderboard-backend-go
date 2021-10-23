package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/model"
	"github.com/speedrun-website/leaderboard-backend/utils"
)

type UserResponse struct {
	ID       uint
	Username string
}

type GetUserErrorResponse struct {
	Error string `json:"error"`
}

type GetUserResponse struct {
	User *model.UserIdentifier `json:"user"`
}

func GetUser(c *gin.Context) {
	// Maybe we shouldn't use the increment ID but generate a UUID instead to avoid
	// exposing the amount of users registered in the database.
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := database.Users.GetUserIdentifierById(id)

	if err != nil {
		var code int
		if errors.Is(err, database.UserNotFoundError{ID: id}) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(code, GetUserErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUserResponse{
		User: user,
	})
}

type RegisterUserConflictResponse struct {
	Error string `json:"error"`
}

type RegisterUserResponse struct {
	User *model.UserIdentifier `json:"user"`
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

	err = database.Users.CreateUser(user)

	if err != nil {
		if uniquenessErr, ok := err.(database.UserUniquenessError); ok {
			/*
			 * TODO: we probably don't want to reveal if an email is already in use.
			 * Maybe just give a 201 and send an email saying that someone tried to sign up as you.
			 * --Ted W
			 *
			 * I still think we should do as above, but for my refactor 2021/10/22 I left
			 * what was already here.
			 * --RageCage
			 */
			c.AbortWithStatusJSON(http.StatusConflict, RegisterUserConflictResponse{
				Error: uniquenessErr.Error(),
			})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/users/%d", user.ID))
	c.JSON(http.StatusCreated, RegisterUserResponse{
		User: &model.UserIdentifier{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}

type MeResponse struct {
	User *model.UserPersonal `json:"user"`
}

func Me(c *gin.Context) {
	rawUser, ok := c.Get(middleware.JwtConfig.IdentityKey)
	if ok {
		user, ok := rawUser.(*model.UserPersonal)
		if ok {
			userInfo, err := database.Users.GetUserPersonalById(uint64(user.ID))

			if err == nil {
				c.JSON(http.StatusOK, MeResponse{
					User: userInfo,
				})
			}
		}
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
