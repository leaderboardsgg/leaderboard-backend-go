package oauth

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/speedrun-website/leaderboard-backend/model"
)

var twitterProvider = twitter.NewAuthenticate(
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

//OauthCallback handles the oauth(1.0a/2) callback mechanism
//The frontend needs to make sure to append ?provider={provider}
func OauthCallback(c *gin.Context) {
	//TODO: Connect to JWT
	provider := c.Query("provider")
	log.Printf("%s oauth callback", provider)

	providerUserData, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var existingUser *model.UserIdentifier
	var userLookupErr error
	if provider == twitterProvider.Name() {
		existingUser, userLookupErr = Oauth.GetUserByTwitterID(providerUserData.UserID)
	}

	if userLookupErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, OauthErrorResponse{
			Error: userLookupErr.Error(),
		})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusOK, gin.H{
			//We "copy" the struct here to ensure that response is consistent
			//even if the type of the existing user changes
			"data": &model.UserIdentifier{
				ID:       existingUser.ID,
				Username: existingUser.Username,
			},
		})
		return
	}
	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	username := fmt.Sprintf("Runner-%d", randNum.Int())
	var createdUser *model.User
	var userCreationError error
	if provider == twitterProvider.Name() {
		//TODO: Attempt to use their twitter username but fallback to random if failed
		newUser := model.User{
			Username:  username,
			Email:     providerUserData.Email,
			TwitterID: &providerUserData.UserID,
		}
		createdUser, userCreationError = Oauth.CreateUser(newUser)
	}

	if userCreationError != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			//TODO Create a standard interface for unique violations
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"errors": [1]gin.H{{"constraint": pgErr.ConstraintName, "message": pgErr.Detail}},
			})
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}
	c.JSON(http.StatusOK, gin.H{
		//We "copy" the struct here to ensure that response is consistent
		//even if the type of the existing user changes
		//to prevent unwanted data from being sent
		"data": &model.UserIdentifier{
			ID:       createdUser.ID,
			Username: createdUser.Username,
		},
	})
}
