package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubModule struct {
	name string
}

func (m stubModule) Name() string { return m.name }

func (m stubModule) Register(c *Container) error {
	c.RegisterRoutes(func(mux *http.ServeMux) {
		mux.HandleFunc("/stub", func(w http.ResponseWriter, r *http.Request) {})
	})
	c.RegisterStartup(func(ctx context.Context) error {
		if ctx == nil {
			return nil
		}
		return nil
	})
	return nil
}

func TestContainerRegisterModules(t *testing.T) {
	c := NewContainer()

	if err := c.RegisterModules(stubModule{name: "demo"}); err != nil {
		t.Fatalf("register modules returned error: %v", err)
	}

	if got := len(c.modules); got != 1 {
		t.Fatalf("expected 1 registered module, got %d", got)
	}

	req := httptest.NewRequest(http.MethodGet, "/stub", nil)
	res := httptest.NewRecorder()
	c.Mux().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200 from registered route, got %d", res.Code)
	}
}
