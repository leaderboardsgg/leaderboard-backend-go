package user

import (
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/server/request"
)

const identityKey = "id"

type TokenResponse struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}

var JwtConfig = &jwt.GinJWTMiddleware{
	Realm:       "test zone",
	Key:         []byte("secret key"),
	Timeout:     time.Hour,
	MaxRefresh:  time.Hour,
	IdentityKey: identityKey,
	PayloadFunc: func(d interface{}) jwt.MapClaims {
		if v, ok := d.(*UserPersonal); ok {
			return jwt.MapClaims{
				identityKey: strconv.FormatUint(uint64(v.ID), 36),
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		idStr := claims[identityKey].(string)
		id, _ := strconv.ParseUint(idStr, 36, 0)
		return &UserPersonal{
			ID: uint(id),
		}
	},
	Authenticator: func(c *gin.Context) (interface{}, error) {
		var loginVals UserLogin
		if err := c.ShouldBindJSON(&loginVals); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		email := loginVals.Email
		password := loginVals.Password

		user, err := Store.GetUserByEmail(email)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		if !ComparePasswords(user.Password, []byte(password)) {
			return nil, jwt.ErrFailedAuthentication
		}

		return user.AsPersonal(), nil
	},
	LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, request.SuccessResponse{
			Data: TokenResponse{
				Token:  token,
				Expiry: expire.Format(time.RFC3339),
			},
		})
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

func GetAuthMiddlewareHandler() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddlware, err := jwt.New(JwtConfig)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	err = authMiddlware.MiddlewareInit()
	if err != nil {
		log.Fatalf("authMiddleware.MiddlewareInit() Error: %s", err)
	}

	return authMiddlware
}
