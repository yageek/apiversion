package apiversion

import (
	"errors"
	"net/http"
	"strings"
)

const (
	vendorSeparator = "application/vnd."
)

// Common errors
var (
	ErrVersionNotFound  = errors.New("Version not found")
	ErrVersionDuplicate = errors.New("Version with same identifer already present")
)

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
	return nil, ErrVersionNotFound
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
			return nil, ErrVersionDuplicate
		}
	}

	middleware.defaultVersionKey = versions[len(versions)-1].Version()
	return middleware, nil
}

// CheckVendorHandler wraps an handler and call it only if the vendor
// corresponds to the appropriate vendor.
func (v *VendorMiddleware) CheckVendorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if acceptHeader := r.Header.Get("Accept"); acceptHeader == "" || !strings.Contains(acceptHeader, v.VendorName()) {
			http.Error(rw, "Wrong vendor identifier", http.StatusNotFound)
		} else {
			h.ServeHTTP(rw, r)
		}
	})
}

// DispatchVersion returns the handler
// that corresponds to the appropriate version.
func (v *VendorMiddleware) DispatchVersion() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		acceptHeader := r.Header.Get("Accept")
		lastIndex := strings.LastIndex(acceptHeader, v.versionSeparator)

		if lastIndex == -1 {
			v.defaultVersion().handler.ServeHTTP(rw, r)
			return
		}
		versionIndex := len(v.versionSeparator) + lastIndex
		version, err := v.version(acceptHeader[versionIndex:])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		version.handler.ServeHTTP(rw, r)
	})
}

// Default implementation in case of non using a middleware
func (v *VendorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.CheckVendorHandler(v.DispatchVersion()).ServeHTTP(w, r)
}
