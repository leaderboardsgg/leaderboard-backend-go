package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/handlers"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/model"
)

// Mock utilities

type mockUserStore struct {
	Users map[uint64]*model.User
}

func setupMockUserStore() *mockUserStore {
	store := mockUserStore{
		Users: map[uint64]*model.User{
			1: {
				Username: "RageCage",
				Email:    "rage@cage.com",
			},
			2: {
				Username: "Squiddo",
				Email:    "she@squiddo.com",
			},
			3: {
				Username: "SiriusCord",
				Email:    "sirius@cord.com",
			},
		},
	}
	database.Users = &store
	return &store
}

func (s mockUserStore) GetUserIdentifierById(userId uint64) (*model.UserIdentifier, error) {
	for id, user := range s.Users {
		if id == userId {
			userIdentifier := model.UserIdentifier{
				ID:       user.ID,
				Username: user.Username,
			}
			return &userIdentifier, nil
		}
	}
	return nil, database.UserNotFoundError
}

func (s mockUserStore) GetUserPersonalById(userId uint64) (*model.UserPersonal, error) {
	for id, user := range s.Users {
		if id == userId {
			userPersonal := model.UserPersonal{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}
			return &userPersonal, nil
		}
	}
	return nil, database.UserNotFoundError
}

func (s mockUserStore) GetUserByEmail(email string) (*model.User, error) {
	return nil, errors.New("This method is unused in this file")
}

func (s mockUserStore) CreateUser(newUser model.User) error {
	var maxId uint64
	maxId = 0
	for id, user := range s.Users {
		if user.Email == newUser.Email {
			return database.UserUniquenessError{
				ErrorField: "email",
			}
		}
		if user.Username == newUser.Username {
			return database.UserUniquenessError{
				ErrorField: "username",
			}
		}
		if id > maxId {
			maxId = id
		}
	}
	s.Users[maxId+1] = &newUser
	return nil
}

// Tests

func TestGetUser400WithoutID(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.GetUser(c)

	testExpectedStatusCode(t, w.Result(), http.StatusBadRequest)
}

func TestGetUser404WithNoUser(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	nonExistentUserId := 69
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: fmt.Sprint(nonExistentUserId),
		},
	}

	handlers.GetUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusNotFound)
	testContentTypeHeader(t, r, "application/json")

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("GetUserResponse", w.Body.String()))
	}
}

func TestGetUser200WithRealUser(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testUserId := 1
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: fmt.Sprint(testUserId),
		},
	}

	handlers.GetUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusOK)
	testContentTypeHeader(t, r, "application/json")

	var response handlers.UserIdentifierResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("GetUserResponse", w.Body.String()))
	}
}

func TestRegisterUser400WithImproperRequestFormat(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	if err := makeJsonBodyPostRequest(c, "{}"); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusBadRequest)
}

func TestRegisterUser400IfPasswordConfirmDoesNotMatch(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        "butterfingers",
		Email:           "doesnot@matter.com",
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrrrrrrrrr",
	}
	if err := makeJsonBodyPostRequest(c, registerBody); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusBadRequest)
}

func TestRegisterUser409WithNonUniqueUsername(t *testing.T) {
	store := setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        store.Users[1].Username,
		Email:           "doesnot@matter.com",
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrd",
	}
	if err := makeJsonBodyPostRequest(c, registerBody); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusConflict)
	testContentTypeHeader(t, r, "application/json")

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("RegisterUserConflictResponse", w.Body.String()))
	}
}

func TestRegisterUser409WithNonUniqueEmail(t *testing.T) {
	store := setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        "somedude",
		Email:           store.Users[1].Email,
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrd",
	}
	if err := makeJsonBodyPostRequest(c, registerBody); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusConflict)
	testContentTypeHeader(t, r, "application/json")

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("RegisterUserConflictResponse", w.Body.String()))
	}
}

func TestRegisterUser201SatisfyingAllRequirements(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        "NewUser",
		Email:           "new@email.com",
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrd",
	}
	err := makeJsonBodyPostRequest(c, registerBody)
	if err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusCreated)
	testContentTypeHeader(t, r, "application/json")

	if !strings.Contains(w.Header().Get("Location"), "/api/v1/users") {
		t.FailNow()
	}

	var response handlers.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("RegisterUserResponse", w.Body.String()))
	}
}

func TestMe500WhenJwtConfigFails(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.Me(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusInternalServerError)
}

func TestMe500WhenRawUserDataCannotBeCasted(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, struct{}{})

	handlers.Me(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusInternalServerError)
}

func TestMe500WhenUserInJWTIsNotReal(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, &model.UserPersonal{
		ID: 69,
	})

	handlers.Me(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusInternalServerError)
}

func TestMe200WhenUserInJWTIsReal(t *testing.T) {
	setupMockUserStore()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, &model.UserPersonal{
		ID: 1,
	})

	handlers.Me(c)

	r := w.Result()
	testExpectedStatusCode(t, r, http.StatusOK)
	testContentTypeHeader(t, r, "application/json")

	var response handlers.UserPersonalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(jsonParseFailureMessage("MeResponse", w.Body.String()))
	}
}

// Utilities

func makeJsonBodyPostRequest(c *gin.Context, content interface{}) error {
	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")

	reqBodyBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	test := string(reqBodyBytes)
	fmt.Print(test)

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))

	return nil
}

func testExpectedStatusCode(t *testing.T, r *http.Response, expectedStatusCode int) {
	if r.StatusCode != expectedStatusCode {
		t.Fatalf(
			"Expected status code %d, got %d",
			expectedStatusCode,
			r.StatusCode,
		)
	}
}

func testContentTypeHeader(t *testing.T, r *http.Response, expectedType string) {
	contentType := r.Header.Get("Content-Type")

	// Sometimes the content-type includes other data
	// like charset that isn't a pass/failure case
	if !strings.HasPrefix(contentType, expectedType) {
		t.Fatalf(
			"Expected Content-Type header %s, got %s",
			expectedType,
			contentType,
		)
	}
}

func jsonParseFailureMessage(typeName, data string) string {
	return fmt.Sprintf(
		"Could not parse data into %s: %s",
		typeName,
		data,
	)
}
