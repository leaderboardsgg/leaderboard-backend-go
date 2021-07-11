package data

import "time"

type Run struct {
	Runner *User
	Game   *Game
	Time   time.Duration
}
