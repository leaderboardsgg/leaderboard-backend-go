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

var Users = map[uint64]*model.User{
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
}

type mockUserStore struct{}

func (s mockUserStore) GetUserIdentifierById(userId uint64) (*model.UserIdentifier, error) {
	for id, user := range Users {
		if id == userId {
			userIdentifier := model.UserIdentifier{
				ID:       user.ID,
				Username: user.Username,
			}
			return &userIdentifier, nil
		}
	}
	return nil, database.UserNotFoundError{ID: userId}
}

func (s mockUserStore) GetUserPersonalById(userId uint64) (*model.UserPersonal, error) {
	for id, user := range Users {
		if id == userId {
			userPersonal := model.UserPersonal{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}
			return &userPersonal, nil
		}
	}
	return nil, database.UserNotFoundError{ID: userId}
}

func (s mockUserStore) GetUserByEmail(email string) (*model.User, error) {
	return nil, errors.New("This method is unused in this file")
}

func (s mockUserStore) CreateUser(newUser model.User) error {
	var maxId uint64
	maxId = 0
	for id, user := range Users {
		if user.Email == newUser.Email {
			return database.UserUniquenessError{
				User:       newUser,
				ErrorField: "email",
			}
		}
		if user.Username == newUser.Username {
			return database.UserUniquenessError{
				User:       newUser,
				ErrorField: "username",
			}
		}
		if id > maxId {
			maxId = id
		}
	}
	Users[maxId+1] = &newUser
	return nil
}

func TestGetUser400WithoutID(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.GetUser(c)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fail()
	}
}

func TestGetUser404WithNoUser(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testUserId := 69
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: fmt.Sprint(testUserId),
		},
	}

	handlers.GetUser(c)

	if w.Result().StatusCode != http.StatusNotFound {
		t.FailNow()
	}
}

func TestGetUser200WithNoUser(t *testing.T) {
	database.Users = mockUserStore{}
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

	if w.Result().StatusCode != http.StatusOK {
		t.FailNow()
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

	if w.Result().StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
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

	if w.Result().StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
}

func TestRegisterUser409WithNonUniqueUsername(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        Users[1].Username,
		Email:           "doesnot@matter.com",
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrd",
	}
	if err := makeJsonBodyPostRequest(c, registerBody); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	if w.Result().StatusCode != http.StatusConflict {
		t.FailNow()
	}
}

func TestRegisterUser409WithNonUniqueEmail(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	registerBody := model.UserRegister{
		Username:        "somedude",
		Email:           Users[1].Email,
		Password:        "str0ngP4sswrd",
		PasswordConfirm: "str0ngP4sswrd",
	}
	if err := makeJsonBodyPostRequest(c, registerBody); err != nil {
		t.FailNow()
	}

	handlers.RegisterUser(c)

	if w.Result().StatusCode != http.StatusConflict {
		t.FailNow()
	}
}

func TestRegisterUser201SatisfyingAllRequirements(t *testing.T) {
	database.Users = mockUserStore{}
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

	if w.Result().StatusCode != http.StatusCreated {
		t.FailNow()
	}
	if !strings.Contains(w.Header().Get("Location"), "/api/v1/users") {
		t.FailNow()
	}
}

func TestMe500WhenJwtConfigFails(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.Me(c)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.FailNow()
	}
}

func TestMe500WhenRawUserDataCannotBeCasted(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, struct{}{})

	handlers.Me(c)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.FailNow()
	}
}

func TestMe500WhenUserInJWTIsNotReal(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, &model.UserPersonal{
		ID: 69,
	})

	handlers.Me(c)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.FailNow()
	}
}

func TestMe200WhenUserInJWTIsReal(t *testing.T) {
	database.Users = mockUserStore{}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JwtConfig.IdentityKey, &model.UserPersonal{
		ID: 1,
	})

	handlers.Me(c)

	if w.Result().StatusCode != http.StatusOK {
		t.FailNow()
	}
}

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
