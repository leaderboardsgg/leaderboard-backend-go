package user

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/server/common"
)

func PublicRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r.POST("/register", RegisterUserHandler)
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)

	r.GET("/users/:id", GetUserHandler)
}

func AuthRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r.GET("/me", MeHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
}

type UserRegister struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserIdentifierResponse struct {
	User *UserIdentifier `json:"user"`
}

type UserPersonalResponse struct {
	User *UserPersonal `json:"user"`
}

func GetUserHandler(c *gin.Context) {
	// Maybe we shouldn't use the increment ID but generate a UUID instead to avoid
	// exposing the amount of users registered in the database.
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := Store.GetUserIdentifierById(uint(id))

	if err != nil {
		var code int
		if errors.Is(err, ErrUserNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(code, common.ErrorResponse{
			Errors: []error{
				err,
			},
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Data: UserIdentifierResponse{
			User: user,
		},
	})
}

func RegisterUserHandler(c *gin.Context) {
	var registerValue UserRegister
	if err := c.BindJSON(&registerValue); err != nil {
		log.Println("Unable to bind value", err)
		return
	}

	hash, err := HashAndSaltPassword([]byte(registerValue.Password))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user := User{
		Username: registerValue.Username,
		Email:    registerValue.Email,
		Password: hash,
	}

	err = Store.CreateUser(&user)

	if err != nil {
		if errors.Is(err, ErrUserNotUnique) {
			/*
			 * TODO: we probably don't want to reveal if an email is already in use.
			 * Maybe just give a 201 and send an email saying that someone tried to sign up as you.
			 * --Ted W
			 *
			 * I still think we should do as above, but for my refactor 2021/10/22 I left
			 * what was already here.
			 * --RageCage
			 */
			c.AbortWithStatusJSON(http.StatusConflict, common.ErrorResponse{
				Errors: []error{
					err,
				},
			})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.Header("Location", fmt.Sprintf("/api/v1/users/%d", user.ID))
	c.JSON(http.StatusCreated, common.SuccessResponse{
		Data: UserIdentifierResponse{
			User: &UserIdentifier{
				ID:       user.ID,
				Username: user.Username,
			},
		},
	})
}

func MeHandler(c *gin.Context) {
	rawUser, ok := c.Get(JwtConfig.IdentityKey)
	if ok {
		user, ok := rawUser.(*UserPersonal)
		if ok {
			userInfo, err := Store.GetUserPersonalById(uint(user.ID))

			if err == nil {
				c.JSON(http.StatusOK, common.SuccessResponse{
					Data: UserPersonalResponse{
						User: userInfo,
					},
				})
				return
			}
		}
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
