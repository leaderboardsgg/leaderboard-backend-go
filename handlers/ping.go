package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Message string `json:"message"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Data: PingResponse{
			Message: "pong",
		},
	})
}
