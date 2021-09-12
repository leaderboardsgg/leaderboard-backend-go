package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"speedrun.website/database"
	"speedrun.website/graph/model"
	"speedrun.website/middleware"
)

func MeHandler(c *gin.Context) {
	user, _ := c.Get(middleware.JwtConfig.IdentityKey)
	db, err := database.GetDatabase()

	// todo error handler or middleware?
	if err != nil {
		log.Println("Unable to connect to database", err)
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer db.Close()

	var me model.User
	result := db.Where(model.User{
		Email: user.(*model.User).Email,
	}).First(&me)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, me)
}
