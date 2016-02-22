package apiversion

import "net/http"

// Version represents an API version.
type Version struct {
	version  string
	handler  http.Handler
	obsolete bool
}

// NewVersion creates a new API with a specified name.
func NewVersion(version string, handler http.Handler) *Version {
	return &Version{version: version, handler: handler, obsolete: false}
}

// Version returns the version of the API.
func (a *Version) Version() string {
	return a.version
}

// Handler returns an handler.
func (a *Version) Handler() http.Handler {
	return a.handler
}

// MakeObsolete indicates an API is obsolete.
func (a *Version) MakeObsolete() {
	a.obsolete = true
}
