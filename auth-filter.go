package auth

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

// Filter provides an authentication filter for requests. Several supported authentications
// can be provided as well as exceptions (path not covered by the filter).
type Filter struct {
	supportedAuth []Auth
	exceptions    []string
}

// Auth defines an interface for a supported auth. This can be used to extend other auth
// besides basic auth and oidc
type Auth interface {
	Authorize(req *restful.Request) bool
}

// NewAuthFilter creates a new AuthFilter
func NewAuthFilter(supportedAuth []Auth, exceptions []string) *Filter {
	return &Filter{
		supportedAuth: supportedAuth,
		exceptions:    exceptions,
	}
}

// Filter runs all authentication filters against the request. If one of the filters return true, then
// the request is authenticated.
func (af *Filter) Filter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	// proceed if no auth is defined
	if len(af.supportedAuth) == 0 {
		chain.ProcessFilter(req, res)
	}

	// proceed if request path is excempted
	if af.exceptions != nil {
		uri := req.Request.URL.RequestURI()
		for _, exception := range af.exceptions {
			if strings.HasPrefix(uri, exception) {
				chain.ProcessFilter(req, res)
				return
			}
		}
	}

	// proceed if at least on auth suceeds
	var success bool
	for _, auth := range af.supportedAuth {
		success = success || auth.Authorize(req)
		if success {
			break
		}
	}
	if !success {
		log.Debug("Unauthorized request rejected")
		res.WriteErrorString(401, "401: Unauthorized")
		return
	}
	chain.ProcessFilter(req, res)
}
