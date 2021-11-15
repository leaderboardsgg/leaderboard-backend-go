package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/handlers"
)

func TestPingHandlerReturnsOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.Ping(c)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf(
			"Expected Ping to return status code %d, got %d",
			http.StatusOK,
			w.Result().StatusCode,
		)
	}

	// First parse into SuccessResponse to get the data out,
	// then re-marshal the contents of data. This allows the test
	// to still work with the concrete PingResponse type.
	var response handlers.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("expected response to be a valid SuccessResponse")
	}
	responseDataBytes, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatal(err.Error())
	}
	var responseData handlers.PingResponse
	err = json.Unmarshal(responseDataBytes, &responseData)
	if err != nil {
		t.Fatalf(
			"Ping response should have expectedPingResponse format, unmarshal failed with %s",
			err,
		)
	}
	if responseData.Message != "pong" {
		t.Fatalf(
			"Ping should pong, instead responded with %s",
			responseData.Message,
		)
	}
}
