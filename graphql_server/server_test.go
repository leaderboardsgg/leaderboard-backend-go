package graphql_server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/speedrun-website/leaderboard-backend/data"

	"github.com/samsarahq/thunder/graphql"
	"github.com/stretchr/testify/assert"
)

type testSqlDriver struct {
	games []*data.Game
}

func (s *testSqlDriver) GetAllGames(ctx context.Context) ([]*data.Game, error) {
	return s.games, nil
}

func (s *testSqlDriver) InsertGame(ctx context.Context, game *data.Game) error {
	s.games = append(s.games, game)
	return nil
}

// TestGamesListAll is an example of a "simple" unit test in GoLang.
func TestGamesListAll(t *testing.T) {
	games := []*data.Game{
		{Title: "Game One"},
		{Title: "Game Two"},
	}
	server := &Server{
		SqlDriver: &testSqlDriver{games: games},
	}

	rr := httptest.NewRecorder()
	handler := graphql.HTTPHandler(server.Schema())

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
		games                       []*data.Game
		titleRegex                  string
		expectedStringsInResponse   []string
		unexpectedStringsInResponse []string
	}{
		{
			desc: "No games",
		},
		{
			desc: "One game - match",
			games: []*data.Game{
				{Title: "First"},
			},
			titleRegex: "irs",
			expectedStringsInResponse: []string{
				"First",
			},
		},
		{
			desc: "One game - no match",
			games: []*data.Game{
				{Title: "First"},
			},
			titleRegex: "Two",
			unexpectedStringsInResponse: []string{
				"First",
			},
		},
		{
			desc: "Three games - two matches",
			games: []*data.Game{
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
			server := &Server{
				SqlDriver: &testSqlDriver{games: tC.games},
			}
			rr := httptest.NewRecorder()
			handler := graphql.HTTPHandler(server.Schema())

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

func TestUsersListAll(t *testing.T) {
	users := []*data.User{
		{Name: "User One"},
		{Name: "User Two"},
	}
	server := &Server{
		Users: users,
	}

	rr := httptest.NewRecorder()
	handler := graphql.HTTPHandler(server.Schema())

	req, err := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query": "query TestQuery {users {name} }"}`))
	assert.NoError(t, err)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

	fullBody, err := ioutil.ReadAll(rr.Result().Body)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(string(fullBody), users[0].Name), "Response should contain this user")
	assert.True(t, strings.Contains(string(fullBody), users[1].Name), "Response should contain this user")
}

func TestUsersFiltering(t *testing.T) {
	testCases := []struct {
		desc                        string
		users                       []*data.User
		nameRegex                   string
		expectedStringsInResponse   []string
		unexpectedStringsInResponse []string
	}{
		{
			desc: "No games",
		},
		{
			desc: "One user - match",
			users: []*data.User{
				{Name: "User"},
			},
			nameRegex: "User",
			expectedStringsInResponse: []string{
				"User",
			},
		},
		{
			desc: "One user - no match",
			users: []*data.User{
				{Name: "UserOne"},
			},
			nameRegex: "Two",
			unexpectedStringsInResponse: []string{
				"UserOne",
			},
		},
		{
			desc: "Three users - two matches",
			users: []*data.User{
				{Name: "Match1"},
				{Name: "Match2"},
				{Name: "Outlier"},
			},
			nameRegex: "Match",
			expectedStringsInResponse: []string{
				"Match1",
				"Match2",
			},
			unexpectedStringsInResponse: []string{
				"Outlier",
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			server := &Server{
				Users: tC.users,
			}
			rr := httptest.NewRecorder()
			handler := graphql.HTTPHandler(server.Schema())

			req, err := http.NewRequest("POST", "/graphql", strings.NewReader(fmt.Sprintf(`{"query": "query TestQuery {users(nameRegex: \"%s\") {name} }"}`, tC.nameRegex)))
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
	games := []*data.Game{
		{Title: "Game One"},
		{Title: "Game Two"},
	}
	server := &Server{
		SqlDriver: &testSqlDriver{games: games},
	}

	rr := httptest.NewRecorder()
	handler := graphql.HTTPHandler(server.Schema())

	req, _ := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query": "query TestQuery {games(titleRegex: \"One") {title} }"}`))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
