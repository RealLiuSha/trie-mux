package trie_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/itchenyi/akita"
	"github.com/itchenyi/trie-mux/mux"
)

func ExampleAkita() {
	app := akita.New()

	router := akita.NewRouter()
	router.Get("/", func(ctx *akita.Context) error {
		return ctx.HTML(200, "<h1>Hello, Akita!</h1>")
	})
	router.Get("/view/:view", func(ctx *akita.Context) error {
		view := ctx.Param("view")
		if view == "" {
			return &akita.Error{Code: 400, Msg: "Invalid view"}
		}
		return ctx.HTML(200, "View: "+view)
	})

	app.UseHandler(router)
	srv := app.Start(":3000")
	defer srv.Close()

	res, _ := http.Get("http://" + srv.Addr().String() + "/view/users")
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	fmt.Println(res.StatusCode, string(body))
	// Output: 200 View: users
}

func ExampleMux() {
	router := mux.New()
	router.Get("/", func(w http.ResponseWriter, _ *http.Request, _ mux.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("<h1>Hello, Akita!</h1>"))
	})

	router.Get("/view/:view", func(w http.ResponseWriter, _ *http.Request, params mux.Params) {
		view := params["view"]
		if view == "" {
			http.Error(w, "Invalid view", 400)
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte("View: " + view))
		}
	})

	// srv := http.Server{Addr: ":3000", Handler: router}
	// srv.ListenAndServe()
	srv := httptest.NewServer(router)
	defer srv.Close()

	res, _ := http.Get(srv.URL + "/view/users")
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	fmt.Println(res.StatusCode, string(body))
	// Output: 200 View: users
}
