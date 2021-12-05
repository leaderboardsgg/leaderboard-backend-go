package middleware

import (
	"log"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/model"
	"github.com/speedrun-website/leaderboard-backend/utils"
)

const identityKey = "id"

var JwtConfig = &jwt.GinJWTMiddleware{
	Realm:       "test zone",
	Key:         []byte("secret key"),
	Timeout:     time.Hour,
	MaxRefresh:  time.Hour,
	IdentityKey: identityKey,
	PayloadFunc: func(d interface{}) jwt.MapClaims {
		if v, ok := d.(*model.UserPersonal); ok {
			return jwt.MapClaims{
				identityKey: strconv.FormatUint(uint64(v.ID), 36),
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		id, _ := strconv.ParseUint(claims[identityKey].(string), 36, 0)
		return &model.UserPersonal{
			ID: uint(id),
		}
	},
	Authenticator: func(c *gin.Context) (interface{}, error) {
		var loginVals model.UserLogin
		if err := c.ShouldBindJSON(&loginVals); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		user, err := database.Users.GetUserByEmail(loginVals.Email)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		if user.Password == nil {
			log.Println("User password in database is null indicating they use oauth and not the password flow")
			return nil, jwt.ErrFailedAuthentication
		}
		passwordMatches, err := utils.ComparePasswords(user.Password, []byte(loginVals.Password))
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		if passwordMatches {
			return &model.UserPersonal{
				ID:       user.ID,
				Email:    user.Email,
				Username: user.Username,
			}, nil
		}

		return &model.UserPersonal{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}, nil
	},
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	},
	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	// - "param:<name>"
	TokenLookup: "header: Authorization, query: token, cookie: jwt",
	// TokenLookup: "query:token",
	// TokenLookup: "cookie:token",

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName: "Bearer",

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc: time.Now,
}

func GetGinJWTMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(JwtConfig)

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware
}
