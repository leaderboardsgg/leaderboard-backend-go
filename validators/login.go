package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/speedrun-website/leaderboard-backend/model"
)

func LoginValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginValue model.UserLogin
		if err := c.ShouldBind(&loginValue); err == nil {
			validate := validator.New()
			if err := validate.Struct(&loginValue); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
