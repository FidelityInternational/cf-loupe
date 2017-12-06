package helpers

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

type FakeApi struct {
	Mux                 *http.ServeMux
	Server              *httptest.Server
	fakeUAAServer       *httptest.Server
	TokenCounter        int
	TokenRefreshCounter int
	MaxTokenRefresh     int
	TokenExpiresIn      int
}

func (api *FakeApi) TeardownFakeApi() {
	api.Server.Close()
	api.fakeUAAServer.Close()
}

func NewFakeApi() *FakeApi {
	api := FakeApi{}

	api.TokenCounter = 0
	api.TokenRefreshCounter = 0
	api.MaxTokenRefresh = MaxInt
	api.TokenExpiresIn = 1

	api.Mux = http.NewServeMux()
	api.Server = httptest.NewServer(api.Mux)
	api.fakeUAAServer = newFakeUAAServer(&api)

	m := martini.New()

	m.Use(render.Renderer())
	r := martini.NewRouter()
	r.Get("/v2/info", func(r render.Render) {
		r.JSON(200, map[string]interface{}{
			"authorization_endpoint": api.fakeUAAServer.URL,
			"token_endpoint":         api.fakeUAAServer.URL,
			"logging_endpoint":       api.Server.URL,
		})

	})
	m.Action(r.Handle)
	api.Mux.Handle("/", m)

	return &api
}

func newFakeUAAServer(fapi *FakeApi) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	fapi.TokenCounter = 0
	fapi.TokenRefreshCounter = 0

	r.Post("/oauth/token", func(r render.Render, req *http.Request) {
		fapi.TokenCounter = fapi.TokenCounter + 1

		grant_type := req.PostFormValue("grant_type")

		if grant_type == "refresh_token" {
			fapi.TokenRefreshCounter = fapi.TokenRefreshCounter + 1
			if fapi.TokenRefreshCounter > fapi.MaxTokenRefresh {
				r.JSON(403, map[string]interface{}{})
				return
			}
		}

		r.JSON(200, map[string]interface{}{
			"token_type":    "bearer",
			"access_token":  "foobar" + strconv.Itoa(fapi.TokenCounter),
			"refresh_token": "barfoo",
			"expires_in":    fapi.TokenExpiresIn,
		})
	})
	r.NotFound(func() string { return "" })
	m.Action(r.Handle)
	mux.Handle("/", m)
	return server
}
