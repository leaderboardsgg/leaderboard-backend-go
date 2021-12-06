package oauth_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/speedrun-website/leaderboard-backend/handlers/oauth"
	"github.com/speedrun-website/leaderboard-backend/model"
)

type mockOauthStore struct {
	Users map[uint]*model.User
}

func makePointerString(str string) *string {
	return &str
}

var rollingID uint = 0

func initStore() *mockOauthStore {
	usersMap := make(map[uint]*model.User)
	usersMap[rollingID] = &model.User{
		Username:  "FastMcGo",
		Email:     "fast@example.com",
		TwitterID: makePointerString(strconv.Itoa(int(rollingID))),
	}
	rollingID++

	store := mockOauthStore{
		Users: usersMap,
	}
	oauth.Oauth = store
	return &store
}

func (s mockOauthStore) GetUserByTwitterID(twitterID string) (*model.User, error) {
	if twitterID == "error" {
		return nil, errors.New("Issue finding user")
	}

	for _, user := range s.Users {
		if *user.TwitterID == twitterID {
			return user, nil
		}
	}
	return nil, nil
}

func (s mockOauthStore) CreateUser(user model.User) (*model.User, error) {
	if strings.Contains(user.Email, "error") {
		return nil, errors.New("Issue creating user")
	}
	rollingID++
	newUser := &model.User{
		ID:        rollingID,
		Email:     user.Email,
		Username:  user.Username,
		TwitterID: makePointerString(strconv.Itoa(int(rollingID))),
	}
	s.Users[rollingID] = newUser
	return newUser, nil
}

func TestOauthCallbackUserAuthError(t *testing.T) {
	oauth.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{}, errors.New("uh oh something went bad")
	}
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	oauth.OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}
	var responseJSON oauth.OauthErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &responseJSON); err != nil {
		t.Fatalf(
			"OauthCallback response expected OauthErrorResponse format, unmarshal failed with %s",
			err,
		)
	}
}

func TestOauthCallbackUserFetchError(t *testing.T) {
	initStore()
	oauth.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: "error",
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	oauth.OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}

	var responseJSON oauth.OauthErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &responseJSON); err != nil {
		t.Fatalf(
			"OauthCallback response expected OauthErrorResponse format, unmarshal failed with %s",
			err,
		)
	}
	if responseJSON.Error != "Issue finding user" {
		t.Fatalf("Expected %s but recieved %s", "Issue finding user", responseJSON.Error)
	}
}

func TestOauthCallbackUserCreationError(t *testing.T) {
	initStore()
	oauth.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: "0",
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	oauth.OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	if err := json.Unmarshal(rec.Body.Bytes(), &responseJSON); err != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			err,
		)
	}
}

func TestOauthCallbackReturnsExistingUser(t *testing.T) {
	store := initStore()
	expectedUser, err := store.CreateUser(model.User{
		Email:    "oauthcallbackereturn@example.com",
		Username: "iliketurtles",
	})

	if err != nil {
		t.Fatalf("issue creating mock user %s", err)
	}

	oauth.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: strconv.Itoa(int(expectedUser.ID)),
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	oauth.OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	if err := json.Unmarshal(rec.Body.Bytes(), &responseJSON); err != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			err,
		)
	}

	if responseJSON.ID != expectedUser.ID {
		t.Fatalf("Expected user %+v but got %+v", expectedUser, responseJSON)
	}
}

func TestOauthCallbackCreatesNewUser(t *testing.T) {
	store := initStore()
	oauth.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		rollingID++
		return goth.User{
			UserID: strconv.Itoa(int(rollingID)),
			Email:  "test@example.com",
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	oauth.OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	if err := json.Unmarshal(rec.Body.Bytes(), &responseJSON); err != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			err,
		)
	}
	expectedUser := store.Users[responseJSON.ID]

	if expectedUser == nil {
		t.Fatalf("Expected a user to be created but was nil response was %+v", responseJSON)
	}
}

func TestInitializeProviders(t *testing.T) {
	if err := os.Setenv("ENABLED_PROVIDERS", "twitter"); err != nil {
		t.Fatalf("issue setting environment variable: %s", err)
	}

	oauth.InitializeProviders()

	if _, err := goth.GetProvider("twitter"); err != nil {
		t.Fatalf("issue getting twitter provider %s", err)
	}

	goth.ClearProviders()

	if err := os.Setenv("ENABLED_PROVIDERS", "fakeoauthprovider"); err != nil {
		t.Fatalf("issue setting environment variable: %s", err)
	}

	if provider, err := goth.GetProvider("twitter"); err == nil {
		t.Fatalf("expected fake provider to error but got a provider back instead %s", provider.Name())
	}

}
