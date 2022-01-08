package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/server/request"
)

func PublicRoutes(r *gin.RouterGroup) {
	r.GET("/ping", pingHandler)
}

func AuthRoutes(r *gin.RouterGroup) {
	r.GET("/authPing", authPingHandler)
}

type PingResponse struct {
	Message string `json:"message"`
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, request.SuccessResponse{
		Data: PingResponse{
			Message: "pong",
		},
	})
}

func authPingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, request.SuccessResponse{
		Data: PingResponse{
			Message: "authenticated pong",
		},
	})
}
