package oauth

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/speedrun-website/leaderboard-backend/model"
)

var availableProviders = map[string]goth.Provider{
	"twitter": twitter.NewAuthenticate(
		os.Getenv("TWITTER_OAUTH_KEY"),
		os.Getenv("TWITTER_OAUTH_SECRET"),
		os.Getenv("TWITTER_OAUTH_CALLBACK_URL"),
	),
}
var CompleteUserAuth = gothic.CompleteUserAuth

func InitializeProviders() {
	enabledProviders := strings.Split(os.Getenv("ENABLED_PROVIDERS"), ",")
	for _, provider := range enabledProviders {
		cleanedProvider := strings.ToLower(strings.TrimSpace(provider))
		maybeProvider, exists := availableProviders[cleanedProvider]
		if exists {
			goth.UseProviders(maybeProvider)
		} else {
			log.Printf("Unknown provider %s", provider)
		}
	}
}

//TODO Replace with generic error response struct
type OauthErrorResponse struct {
	Error string `json:"error"`
}

//OauthLogin begins the oauth authentication process
//such as fetching tokens and redirecting the user to the providers website.
func OauthLogin(c *gin.Context) {
	log.Printf("%s oauth authentication", c.Query("provider"))
	// Handles redirecting the user
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

//OauthCallback handles the oauth(1.0a/2) callback mechanism
//The frontend needs to make sure to append ?provider={provider}
func OauthCallback(c *gin.Context) {
	//TODO: Connect to JWT
	providerName := c.Query("provider")
	log.Printf("%s oauth callback", providerName)

	//We dont need to check for provider existance further down
	//since this will error if the provider is one of our supported ones
	providerUserData, err := CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &OauthErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var existingUser *model.User
	var userLookupErr error

	//Go doesnt "support" dynamic lookups like say Javascript so we need unique functions for each provider
	//(you can kind of cheat this in less than ideal ways but copy and paste a few times is better than interface{} shenanigans)
	// aka this(and similar) if statement(s) *may* grow but its one of the clearest ways to handle this
	// another option would be to shift the if statement into a seperate function but w/e
	// I dont ever expect more than a few providers honestly
	if providerName == "twitter" {
		existingUser, userLookupErr = Oauth.GetUserByTwitterID(providerUserData.UserID)
	}

	if userLookupErr != nil {
		log.Printf("error looking up user: %s", userLookupErr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, OauthErrorResponse{
			Error: "Issue finding user",
		})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusOK,
			//We "copy" the struct here to ensure that response is consistent
			//even if the type of the existing user changes
			&model.UserIdentifier{
				ID:       existingUser.ID,
				Username: existingUser.Username,
			},
		)
		return
	}
	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	username := fmt.Sprintf("Runner-%d", randNum.Int())
	var createdUser *model.User
	var userCreationError error
	if providerName == "twitter" {
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
			//TODO: Create a standard interface for unique violations
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"errors": [1]gin.H{{"constraint": pgErr.ConstraintName, "message": pgErr.Detail}},
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, OauthErrorResponse{
				Error: "Issue creating user",
			})
		}

		return
	}
	c.JSON(http.StatusOK,
		//We "copy" the struct here to ensure that response is consistent
		//even if the type of the existing user changes
		//to prevent unwanted data from being sent
		&model.UserIdentifier{
			ID:       createdUser.ID,
			Username: createdUser.Username,
		},
	)
}
