package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	database "speedrun.website/db"
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
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer db.Close()

	var me model.User
	db.Model(model.User{
		Username: user.(*model.User).Username,
	}).Limit(1).Find(&me)

	c.JSON(200, me)
}
