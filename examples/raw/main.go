package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/yageek/apiversion"
)

var routingAPI *apiversion.VendorMiddleware

func init() {
	routingAPI, _ = apiversion.NewVendorMiddleware("gecker.io", ApiV1(), ApiV2())
}
func main() {
	http.ListenAndServe(":"+os.Getenv("PORT"), routingAPI)
}

func ApiV1() *apiversion.Version {

	r := pat.New()
	r.Get("/", HandlerV1)

	return apiversion.NewVersion("v1", r)
}

func ApiV2() *apiversion.Version {
	r := pat.New()
	r.Get("/", HandlerV2)

	return apiversion.NewVersion("v2", r)
}

func HandlerV1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello v1")
}

func HandlerV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello v2")
}
