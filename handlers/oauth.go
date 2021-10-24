package handlers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/model"
	"gorm.io/gorm"
)

var twitterProvider = twitter.New(
	os.Getenv("TWITTER_OAUTH_KEY"),
	os.Getenv("TWITTER_OAUTH_SECRET"),
	os.Getenv("TWITTER_OAUTH_CALLBACK_URL"),
)

func InitializeProviders() {
	goth.UseProviders(twitterProvider)
}

type OauthErrorResponse struct {
	Error string `json:"error"`
}

func OauthLogin(c *gin.Context) {
	log.Printf("%s oauth authentication", c.Param("provider"))
	// Handles redirecting the user
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OauthCallback(c *gin.Context) {
	log.Printf("%s oauth callback", c.Param("provider"))
	providerUserData, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	var maybeExistingUser model.UserIdentifier
	result := database.DB.Model(&model.User{}).Where("twitter_id = ?", providerUserData.UserID).First(&maybeExistingUser)

	if result.Error != nil {
		//@TODO: Implement error redirects
		c.AbortWithStatusJSON(http.StatusInternalServerError, OauthErrorResponse{
			Error: result.Error.Error(),
		})
		return
	}

	userExists := !errors.Is(err, gorm.ErrRecordNotFound)
	if userExists {
		//@TODO: Setup JWT
		c.Redirect(http.StatusOK, "/")
		return
	}
	//@TODO: Handle failures for random username
	newUser := model.User{
		Username: fmt.Sprintf("Runner-%d", rand.Int()),
		Email:    providerUserData.Email,
	}

	result = database.DB.WithContext(c).Create(&newUser)

	if result.Error != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"errors": [1]gin.H{{"constraint": pgErr.ConstraintName, "message": pgErr.Detail}},
			})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}
	//@TODO: Setup JWT
	c.Redirect(http.StatusOK, "/")

}
