package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/workos/workos-go/v10"
)

func (cfg *Config) LoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := "authkit"
	authURL, err := cfg.Client.GetAuthKitAuthorizationURL(workos.AuthKitAuthorizationURLParams{
		RedirectURI: cfg.RedirectURI,
		Provider:    &provider,
	})
	if err != nil {
		http.Error(w, "An error an ocurred getting the URL.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (cfg *Config) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Parameter 'code' is missing", http.StatusBadRequest)
		return
	}

	resp, err := cfg.Client.UserManagement().AuthenticateWithCode(r.Context(), &workos.UserManagementAuthenticateWithCodeParams{
		Code: code,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("An error an ocurred authenticating the code %s", err), http.StatusUnauthorized)
		return
	}

	if resp.User == nil {
		http.Error(w, "User is not returned", http.StatusInternalServerError)
		return
	}

	if err := CreateSession(w, *resp.User); err != nil {
		http.Error(w, "Error creating the session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (cfg *Config) MeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetSession(r)
	if !ok {
		http.Error(w, "Not Authenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding the user", http.StatusInternalServerError)
	}
}

func (cfg *Config) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	DestroySession(w, r)
	w.WriteHeader(http.StatusNoContent)
}
