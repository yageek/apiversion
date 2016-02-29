package adapters

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/yageek/apiversion"
)

// NegroniVendorAdapter is an adaptor for Negroni.
func NegroniVendorAdapter(v *apiversion.VendorMiddleware) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v.CheckVendorHandler(next).ServeHTTP(rw, r)
	}
}
