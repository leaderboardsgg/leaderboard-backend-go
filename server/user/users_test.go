package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/server/common"
	"github.com/speedrun-website/leaderboard-backend/server/user"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Where's the .env file?")
	}

	database.InitTest()
	user.InitGormStore(nil)
}

func TestAuthFlow(t *testing.T) {
	t.Parallel()

	r := getUsersContext()

	cleanup := []uint{}
	t.Run("Full auth flow", func(t *testing.T) {
		password := "str0ng3stp4ssw0rd"
		userReg := user.UserRegister{
			Username:        "AGoodUser",
			Email:           "cool@email.com",
			Password:        password,
			PasswordConfirm: password,
		}
		u := testRegister(t, r, userReg)
		cleanup = append(cleanup, u.ID)

		userLogin := user.UserLogin{
			Email:    u.Email,
			Password: password,
		}
		token := testLogin(t, r, userLogin)

		testMe(t, r, *u, token)

		testRefreshToken(t, r, token)
	})
	if err := cleanupUsers(cleanup); err != nil {
		t.Fatalf("cleanup failed: %s", err)
	}
}

func testRegister(
	t *testing.T,
	r *gin.Engine,
	registerBody user.UserRegister,
) *user.UserPersonal {
	t.Helper()

	responseBytes, err := testJsonPostRequest(r, "/register", registerBody, http.StatusCreated)
	if err != nil {
		// FIXME
		t.Fatalf("it failed: %s", err)
	}
	var responseData user.UserIdentifierResponse
	_, err = common.UnmarshalSuccessResponseData(responseBytes, &responseData)
	if err != nil {
		// FIXME
		t.Fatal("bad response format")
	}
	user, err := user.Store.GetUserPersonalById(responseData.User.ID)
	if err != nil {
		t.Fatalf("failed to register user: %s", err)
	}

	return user
}

func testLogin(
	t *testing.T,
	r *gin.Engine,
	loginBody user.UserLogin,
) string {
	t.Helper()

	responseBytes, err := testJsonPostRequest(r, "/login", loginBody, http.StatusOK)
	if err != nil {
		// FIXME
		t.Fatal("login failed")
	}
	var response user.TokenResponse
	_, err = common.UnmarshalSuccessResponseData(responseBytes, &response)
	if err != nil {
		// FIXME
		t.Fatal("login failed response bad")
	}

	return response.Token
}

func testMe(
	t *testing.T,
	r *gin.Engine,
	u user.UserPersonal,
	token string,
) {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	responseBytes, err := testGetRequest(r, "/me", http.StatusOK, request)
	if err != nil {
		// FIXME
		t.Fatalf("me failed: %s", err)
	}
	var responseData user.UserPersonalResponse
	_, err = common.UnmarshalSuccessResponseData(responseBytes, &responseData)
	if err != nil {
		// FIXME
		t.Fatal("me failed response bad")
	}

	if u.ID != responseData.User.ID {
		// FIXME
		t.Fatalf("me failed: %s", err)
	}
}

func testRefreshToken(
	t *testing.T,
	r *gin.Engine,
	token string,
) {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/refresh_token", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	responseBytes, err := testGetRequest(r, "/refresh_token", http.StatusOK, request)
	if err != nil {
		// FIXME
		t.Fatalf("refresh_token failed: %s", err)
	}
	var response struct {
		Token string `json:"token"`
	}
	_, err = common.UnmarshalSuccessResponseData(responseBytes, &response)
	if err != nil {
		// FIXME
		t.Fatal("refresh_token failed response bad")
	}
	if token == response.Token {
		// FIXME
		t.Fatal("what the fuck")
	}
}

func TestPOSTRegister400(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		body user.UserRegister
	}{
		{
			name: "Mismatch password",
			body: user.UserRegister{
				Username:        "RageCage",
				Email:           "x@y.com",
				Password:        "beepboopbo",
				PasswordConfirm: "beepboopbop",
			},
		},
		{
			name: "Too short password",
			body: user.UserRegister{
				Username:        "RageCage",
				Email:           "x@y.com",
				Password:        "2",
				PasswordConfirm: "2",
			},
		},
		{
			name: "Invalid email",
			body: user.UserRegister{
				Username:        "RageCage",
				Email:           "bepis",
				Password:        "beepboopbo",
				PasswordConfirm: "beepboopbo",
			},
		},
	}

	for _, testCase := range testCases {
		r := getUsersContext()

		t.Run(testCase.name, func(t *testing.T) {
			_, err := testJsonPostRequest(r, "/register", testCase.body, http.StatusBadRequest)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPOSTRegister409(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		setupUser user.UserRegister
		user      user.UserRegister
	}{
		{
			name: "Conflicting email address",
			setupUser: user.UserRegister{
				Username:        "Squiddo",
				Email:           "same@email.com",
				Password:        "beepboopbop",
				PasswordConfirm: "beepboopbop",
			},
			user: user.UserRegister{
				Username:        "Siriuscord",
				Email:           "same@email.com",
				Password:        "beepboopbop",
				PasswordConfirm: "beepboopbop",
			},
		},
		{
			name: "Conflicting username",
			setupUser: user.UserRegister{
				Username:        "SomeDude",
				Email:           "different@email.com",
				Password:        "beepboopbop",
				PasswordConfirm: "beepboopbop",
			},
			user: user.UserRegister{
				Username:        "SomeDude",
				Email:           "another@email.com",
				Password:        "beepboopbop",
				PasswordConfirm: "beepboopbop",
			},
		},
	}

	cleanup := []uint{}
	for _, testCase := range testCases {
		r := getUsersContext()

		t.Run(testCase.name, func(t *testing.T) {
			_, err := testJsonPostRequest(r, "/register", testCase.setupUser, http.StatusCreated)
			if err != nil {
				t.Fatal(err)
			}
			u, err := user.Store.GetUserByEmail(testCase.setupUser.Email)
			if err != nil {
				t.Fatal(err)
			}
			cleanup = append(cleanup, u.ID)

			_, err = testJsonPostRequest(r, "/register", testCase.user, http.StatusConflict)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	if err := cleanupUsers(cleanup); err != nil {
		t.Fatalf("failed cleanup: %s", err)
	}
}

func TestLogin401(t *testing.T) {
	t.Parallel()

	r := getUsersContext()
	email := "email@cool.com"
	password := "beepboopbop"
	setupUser := user.UserRegister{
		Username:        "UserWhoIsBadAtLoggingIn",
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
	}
	_, err := testJsonPostRequest(r, "/register", setupUser, http.StatusCreated)
	if err != nil {
		t.Fatalf("Registering setup user failed: %s", err)
	}
	u, err := user.Store.GetUserByEmail(setupUser.Email)
	if err != nil {
		t.Fatalf("Retrieving setup user failed: %s", err)
	}

	testCases := []struct {
		name      string
		loginBody user.UserLogin
	}{
		{
			name: "Wrong password",
			loginBody: user.UserLogin{
				Email:    email,
				Password: "othertext",
			},
		},
		{
			name: "Email doesn't exist",
			loginBody: user.UserLogin{
				Email:    "someother@email.com",
				Password: password,
			},
		},
		{
			name: "Email isn't a real email",
			loginBody: user.UserLogin{
				Email:    "garbagetext",
				Password: password,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testJsonPostRequest(r, "/login", testCase.loginBody, http.StatusUnauthorized)
			if err != nil {
				t.Fatalf("register failed: %s", err)
			}
		})
	}

	cleanupUsers([]uint{u.ID})
}

func getUsersContext() *gin.Engine {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	api := r.Group("/")
	authMiddleware := user.GetAuthMiddlewareHandler()
	user.PublicRoutes(api, authMiddleware)
	api.Use(authMiddleware.MiddlewareFunc())
	{
		user.AuthRoutes(api, authMiddleware)
	}
	return r
}

func testJsonPostRequest(
	r *gin.Engine,
	target string,
	content interface{},
	expectedStatusCode int,
) ([]byte, error) {
	reqBodyBytes, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf(
			"could not marshal %s into json",
			content,
		)
	}
	reqBodyBuffer := ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
	req := httptest.NewRequest(http.MethodPost, target, reqBodyBuffer)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	res := w.Result()
	if res.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf(
			"expected status code %d, got %d",
			expectedStatusCode,
			res.StatusCode,
		)
	}
	return w.Body.Bytes(), nil
}

func testGetRequest(
	r *gin.Engine,
	target string,
	expectedStatusCode int,
	request *http.Request,
) ([]byte, error) {
	var req *http.Request
	if request == nil {
		req = httptest.NewRequest(http.MethodGet, target, nil)
	} else {
		req = request
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	res := w.Result()
	if res.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf(
			"expected status code %d, got %d",
			expectedStatusCode,
			res.StatusCode,
		)
	}
	return w.Body.Bytes(), nil
}

func cleanupUsers(usersToDelete []uint) error {
	for _, id := range usersToDelete {
		if err := user.Store.DeleteUser(id); err != nil {
			return err
		}
		if err := user.Store.DumpDeleted(); err != nil {
			return err
		}
	}
	return nil
}
