package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/speedrun-website/leaderboard-backend/model"
)

func RegisterValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerValue model.UserRegister
		if err := c.ShouldBind(&registerValue); err == nil {
			validate := validator.New()
			if err := validate.Struct(&registerValue); err != nil {
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
