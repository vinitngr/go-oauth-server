package oauth

import "net/http"

type OauthProcessor interface {
	Connect(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}
