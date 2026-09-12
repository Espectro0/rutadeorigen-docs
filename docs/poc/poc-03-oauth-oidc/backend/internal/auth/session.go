package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/workos/workos-go/v10"
)

const sessionCookieName = "session_id"

var sessionStore = struct {
	mu   sync.RWMutex
	data map[string]workos.User
}{data: make(map[string]workos.User)}

func CreateSession(w http.ResponseWriter, user workos.User) error {
	id, err := generateSessionID()
	if err != nil {
		return err
	}

	sessionStore.mu.Lock()
	sessionStore.data[id] = user
	sessionStore.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func GetSession(r *http.Request) (workos.User, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return workos.User{}, false
	}

	sessionStore.mu.RLock()
	user, ok := sessionStore.data[cookie.Value]
	sessionStore.mu.RUnlock()

	return user, ok
}

func DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		sessionStore.mu.Lock()
		delete(sessionStore.data, cookie.Value)
		sessionStore.mu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
