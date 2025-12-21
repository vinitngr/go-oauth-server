package oauth

import "github.com/vinitngr/go-oauth-server/internals/config"

type OauthRegistry struct {
	handlers map[string]OauthProcessor
}

func NewOauthRegistry(cfg config.Config) *OauthRegistry {
	return &OauthRegistry{
		handlers: map[string]OauthProcessor{
			"github": NewGithub(cfg),
			"google" : NewGoogle(cfg),
		},
	}
}


func (r *OauthRegistry) Get(provider string) (OauthProcessor, bool) {
	h, ok := r.handlers[provider]
	return h, ok
}
