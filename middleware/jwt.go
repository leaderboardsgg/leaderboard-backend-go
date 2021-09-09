package middleware

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	database "speedrun.website/db"
	"speedrun.website/graph/model"
	"speedrun.website/utils"
)

const identityKey = "email"

var JwtConfig = &jwt.GinJWTMiddleware{
	Realm:       "test zone",
	Key:         []byte("secret key"),
	Timeout:     time.Hour,
	MaxRefresh:  time.Hour,
	IdentityKey: identityKey,
	PayloadFunc: func(d interface{}) jwt.MapClaims {
		if v, ok := d.(*model.User); ok {
			return jwt.MapClaims{
				identityKey: v.Email,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &model.User{
			Email: claims[identityKey].(string),
		}
	},
	Authenticator: func(c *gin.Context) (interface{}, error) {
		var loginVals model.Login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		email := loginVals.Email
		password := loginVals.Password

		db, err := database.GetDatabase()

		if err != nil {
			log.Println("Unable to connect to database", err)
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return nil, err
		}

		defer db.Close()

		var user model.User
		result := db.Where(model.User{
			Email: email,
		}).Limit(1).Find(&user)

		if result.Error != nil {
			return nil, result.Error
		}

		if result.RowsAffected == 0 {
			return nil, jwt.ErrFailedAuthentication
		}

		if utils.ComparePasswords(user.Password, []byte(password)) {
			return &model.User{
				Email: email,
			}, nil
		}

		return nil, jwt.ErrFailedAuthentication
	},
	Authorizator: func(d interface{}, c *gin.Context) bool {
		if _, ok := d.(*model.User); ok {
			return true
		}

		return false
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
