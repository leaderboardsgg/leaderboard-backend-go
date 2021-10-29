package oauth

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
	Oauth = store
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

func Test_OauthCallbackUserAuthError(t *testing.T) {
	completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{}, errors.New("uh oh something went bad")
	}
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}
	var responseJSON OauthErrorResponse
	jsonError := json.Unmarshal(rec.Body.Bytes(), &responseJSON)
	if jsonError != nil {
		t.Fatalf(
			"OauthCallback response expected OauthErrorResponse format, unmarshal failed with %s",
			jsonError,
		)
	}
}

func Test_OauthCallbackUserFetchError(t *testing.T) {
	initStore()
	completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: "error",
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, httpRequestErr := http.NewRequest("POST", "/", nil)
	if httpRequestErr != nil {
		t.Fatal(httpRequestErr)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}

	var responseJSON OauthErrorResponse
	jsonError := json.Unmarshal(rec.Body.Bytes(), &responseJSON)
	if jsonError != nil {
		t.Fatalf(
			"OauthCallback response expected OauthErrorResponse format, unmarshal failed with %s",
			jsonError,
		)
	}
	if responseJSON.Error != "Issue finding user" {
		t.Fatalf("Expected %s but recieved %s", "Issue finding user", responseJSON.Error)
	}
}

func Test_OauthCallbackUserCreationError(t *testing.T) {
	initStore()
	completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: "0",
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, httpRequestErr := http.NewRequest("POST", "/", nil)
	if httpRequestErr != nil {
		t.Fatal(httpRequestErr)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	jsonError := json.Unmarshal(rec.Body.Bytes(), &responseJSON)
	if jsonError != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			jsonError,
		)
	}
}

func Test_OauthCallbackReturnsExistingUser(t *testing.T) {
	store := initStore()
	expectedUser, createUserErr := store.CreateUser(model.User{
		Email:    "oauthcallbackereturn@example.com",
		Username: "iliketurtles",
	})

	if createUserErr != nil {
		t.Fatalf("issue creating mock user %s", createUserErr)
	}

	completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return goth.User{
			UserID: strconv.Itoa(int(expectedUser.ID)),
		}, nil
	}
	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	newQuery := url.Values{
		"provider": []string{"twitter"},
	}
	req, httpRequestErr := http.NewRequest("POST", "/", nil)
	if httpRequestErr != nil {
		t.Fatal(httpRequestErr)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	jsonError := json.Unmarshal(rec.Body.Bytes(), &responseJSON)
	if jsonError != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			jsonError,
		)
	}

	if responseJSON.ID != expectedUser.ID {
		t.Fatalf("Expected user %+v but got %+v", expectedUser, responseJSON)
	}
}

func Test_OauthCallbackCreatesNewUser(t *testing.T) {
	store := initStore()
	completeUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
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

	req, httpRequestErr := http.NewRequest("POST", "/", nil)
	if httpRequestErr != nil {
		t.Fatal(httpRequestErr)
	}
	req.URL.RawQuery = newQuery.Encode()
	ctx.Request = req
	OauthCallback(ctx)
	result := rec.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, result.StatusCode)
	}

	var responseJSON model.UserIdentifier
	jsonError := json.Unmarshal(rec.Body.Bytes(), &responseJSON)
	if jsonError != nil {
		t.Fatalf(
			"OauthCallback response expected UserIdentifier format, unmarshal failed with %s",
			jsonError,
		)
	}
	expectedUser := store.Users[responseJSON.ID]

	if expectedUser == nil {
		t.Fatalf("Expected a user to be created but was nil response was %+v", responseJSON)
	}
}

func Test_InitializeProviders(t *testing.T) {
	setEnvErr := os.Setenv("ENABLED_PROVIDERS", "twitter")
	if setEnvErr != nil {
		t.Fatalf("issue setting environment variable: %s", setEnvErr)
	}

	InitializeProviders()

	_, getProviderErr := goth.GetProvider(twitterProvider.Name())
	if getProviderErr != nil {
		t.Fatalf("issue getting twitter provider %s", getProviderErr)
	}
	goth.ClearProviders()

	setEnvErr = os.Setenv("ENABLED_PROVIDERS", "fakeoauthprovider")
	if setEnvErr != nil {
		t.Fatalf("issue setting environment variable: %s", setEnvErr)
	}

	provider, getProviderErr := goth.GetProvider(twitterProvider.Name())
	if getProviderErr == nil {
		t.Fatalf("expected fake provider to error but got a provider back instead %s", provider.Name())
	}

}
