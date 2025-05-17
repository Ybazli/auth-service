package utils

import "fmt"

func SessionKey(sessionID string) string {
	return sessionID
}

func UserSessionsKey(userID uint) string {
	return "user_sessions:" + fmt.Sprint(userID)
}
