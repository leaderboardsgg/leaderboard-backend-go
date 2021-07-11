package auth

// FIXME lol
func VerifyAPIKey(key string) bool {
	return key == "apikey"
}
