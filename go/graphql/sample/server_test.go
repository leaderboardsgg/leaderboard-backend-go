package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/samsarahq/thunder/graphql"
	"github.com/stretchr/testify/assert"
)

// TestGamesListAll is an example of a "simple" unit test in GoLang.
func TestGamesListAll(t *testing.T) {
	games := []*game{
		{Title: "Game One"},
		{Title: "Game Two"},
	}
	server := &server{
		games: games,
	}

	rr := httptest.NewRecorder()
	handler := graphql.HTTPHandler(server.schema())

	req, err := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query": "query TestQuery {games {title} }"}`))
	assert.NoError(t, err)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

	fullBody, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(string(fullBody), games[0].Title), "Response should contain this game")
	assert.True(t, strings.Contains(string(fullBody), games[1].Title), "Response should contain this game")
}

// TestGamesFiltering is an example of a "table driven" unit test in GoLang.
func TestGamesFiltering(t *testing.T) {
	testCases := []struct {
		desc                        string
		games                       []*game
		titleRegex                  string
		expectedStringsInResponse   []string
		unexpectedStringsInResponse []string
	}{
		{
			desc: "No games",
		},
		{
			desc: "One game - match",
			games: []*game{
				{Title: "First"},
			},
			titleRegex: "irs",
			expectedStringsInResponse: []string{
				"First",
			},
		},
		{
			desc: "One game - no match",
			games: []*game{
				{Title: "First"},
			},
			titleRegex: "Two",
			unexpectedStringsInResponse: []string{
				"First",
			},
		},
		{
			desc: "Three games - two matches",
			games: []*game{
				{Title: "This matches first"},
				{Title: "This will not"},
				{Title: "This matches second"},
			},
			titleRegex: "matches",
			expectedStringsInResponse: []string{
				"This matches first",
				"This matches second",
			},
			unexpectedStringsInResponse: []string{
				"This will not",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			server := &server{
				games: tC.games,
			}
			rr := httptest.NewRecorder()
			handler := graphql.HTTPHandler(server.schema())

			req, err := http.NewRequest("POST", "/graphql", strings.NewReader(fmt.Sprintf(`{"query": "query TestQuery {games(titleRegex: \"%s\") {title} }"}`, tC.titleRegex)))
			assert.NoError(t, err)

			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

			fullBody, err := ioutil.ReadAll(rr.Result().Body)
			assert.NoError(t, err)
			fullBodyStr := string(fullBody)

			for _, expectedString := range tC.expectedStringsInResponse {
				assert.True(t, strings.Contains(fullBodyStr, expectedString), "%s was expected in Response: %s", expectedString, fullBodyStr)
			}
			for _, unexpectedString := range tC.unexpectedStringsInResponse {
				assert.False(t, strings.Contains(fullBodyStr, unexpectedString), "%s was not expected in Response: %s", unexpectedString, fullBodyStr)
			}
		})
	}
}

// BenchmarkGameFiltering is an example of a benchmark in GoLang.
func BenchmarkGameFiltering(b *testing.B) {
	games := []*game{
		{Title: "Game One"},
		{Title: "Game Two"},
	}
	server := &server{
		games: games,
	}

	rr := httptest.NewRecorder()
	handler := graphql.HTTPHandler(server.schema())

	req, _ := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query": "query TestQuery {games(titleRegex: \"One") {title} }"}`))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
