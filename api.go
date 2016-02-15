package api

import (
	"errors"
	"net/http"
	"strings"
)

const (
	vendorSeparator = "application/vnd."
)

// Version represents an API version.
type Version struct {
	version  string
	handler  http.Handler
	obsolete bool
}

// NewAPI creates a new API with a specified name.
func NewAPI(version string, handler http.Handler) *Version {
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

// VendorMiddleware dispatches the request
// regarding the wanted version.
type VendorMiddleware struct {
	vendorName string
	versions   map[string]*Version
}

func (v *VendorMiddleware) version(versionName string) (*Version, error) {
	for key := range v.versions {
		if key == versionName {
			return v.versions[key], nil
		}
	}
	return nil, errors.New("Version not found")
}

// VendorName returns the vendorName used
// to determine the vendor used in the
// "Accept" header.
func (v *VendorMiddleware) VendorName() string {
	return vendorSeparator + v.vendorName
}

// NewVendorMiddleware returns a new middleware.
func NewVendorMiddleware(name string) *VendorMiddleware {
	return &VendorMiddleware{vendorName: name}
}

func (v *VendorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if acceptVersion := r.Header.Get("Accept"); acceptVersion != v.VendorName() {
		http.Error(w, "Unknown vendor", http.StatusNotFound)
	} else {

		lastIndex := strings.LastIndex(acceptVersion, vendorSeparator)

		if lastIndex == -1 {
			http.Error(w, "Can not read accepted version", http.StatusNotFound)
			return
		}
		version, err := v.version(acceptVersion[lastIndex:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		version.handler.ServeHTTP(w, r)
	}

}
