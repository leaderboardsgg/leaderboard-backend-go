package oauth

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
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

func OauthCallback(c *gin.Context) {
	log.Printf("%s oauth callback", c.Param("provider"))

	// Technically goth works with URL params but not gins it seems
	// So we manually set a query paramater here
	// Its a bit hacky but :shrug:
	// goth really should just allow people to pass in provider...
	queryString := c.Request.URL.Query()
	queryString.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = queryString.Encode()

	providerUserData, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	existingUser, err := Oauth.GetUserByTwitterID(providerUserData.UserID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, OauthErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if existingUser != nil {
		userTemplate := `
<p><a href="/logout/twitter">logout</a></p>
<p>Name: {{.Email}}</p>
`
		c.Status(http.StatusOK)
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(c.Writer, existingUser)
		//@TODO: Setup JWT
		// c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	//TODO: Attempt to use their twitter username but fallback to random if failed
	newUser := model.User{
		Username:  fmt.Sprintf("Runner-%d", randNum.Int()),
		Email:     providerUserData.Email,
		TwitterID: &providerUserData.UserID,
	}
	createdUser, err := Oauth.CreateUser(newUser)

	if err != nil {
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
	//TODO: Setup JWT
	// c.Redirect(http.StatusOK, "/")
	userTemplate := `
<p><a href="/logout/twitter">logout</a></p>
<p>Name: {{.Email}}</p>
`
	c.Status(http.StatusOK)
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(c.Writer, createdUser)

}
