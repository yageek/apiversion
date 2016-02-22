# Versioning API with Go

Simple API versioning with Golang.

# Installation

```
go get -v github.com/yageek/apiversion
```


# Example

```
package main

import (
	"os"
    "github.com/gorilla/pat"
	"github.com/yageek/apiversion"
	"github.com/codegangsta/negroni"
)

func main() {
	n := negroni.Classic()
	n.UseHandler(Router())
	n.Run(":" + os.Getenv("PORT"))
}

func Router() *api.VendorMiddleware {
	routing, _ := api.NewVendorMiddleware("gecker.io", ApiV1(), ApiV2())
	return routing
}

func ApiV1() *api.Version {

	r := pat.New()
	r.Get("/", HandlerV1)

	return api.NewAPI("v1", r)
}

func ApiV2() *api.Version {
	r := pat.New()
	r.Get("/", HandlerV2)

	return api.NewAPI("v2", r)
}

func HandlerV1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello v1")
}

func HandlerV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello v2")
}
```


