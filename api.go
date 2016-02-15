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
	vendorName        string
	versions          map[string]*Version
	defaultVersionKey string
	versionSeparator  string
}

func (v *VendorMiddleware) version(versionName string) (*Version, error) {
	for key := range v.versions {
		if key == versionName {
			return v.versions[key], nil
		}
	}
	return nil, errors.New("Version not found")
}

func (v *VendorMiddleware) defaultVersion() *Version {
	return v.versions[v.defaultVersionKey]
}

// VendorName returns the vendorName used
// to determine the vendor used in the
// "Accept" header.
func (v *VendorMiddleware) VendorName() string {
	return vendorSeparator + v.vendorName
}

// NewVendorMiddleware returns a new middleware.
func NewVendorMiddleware(name string, versions ...*Version) (*VendorMiddleware, error) {
	middleware := &VendorMiddleware{
		vendorName:       name,
		versions:         make(map[string]*Version, len(versions)),
		versionSeparator: "-",
	}

	for _, version := range versions {
		if _, ok := middleware.versions[version.Version()]; !ok {
			middleware.versions[version.Version()] = version
		} else {
			return nil, errors.New("Version with same identifer already present")
		}
	}

	middleware.defaultVersionKey = versions[len(versions)-1].Version()
	return middleware, nil
}

func (v *VendorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if acceptHeader := r.Header.Get("Accept"); acceptHeader == "" || !strings.Contains(acceptHeader, v.VendorName()) {
		http.Error(w, "Wrong vendor identifier", http.StatusNotFound)
	} else {

		lastIndex := strings.LastIndex(acceptHeader, v.versionSeparator)
		if lastIndex == -1 {
			v.defaultVersion().handler.ServeHTTP(w, r)
			return
		}
		versionIndex := len(v.versionSeparator) + lastIndex
		version, err := v.version(acceptHeader[versionIndex:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		version.handler.ServeHTTP(w, r)
	}

}
