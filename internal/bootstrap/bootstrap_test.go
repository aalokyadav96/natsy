package bootstrap

import (
	"testing"

	"nae/internal/shared/infra"
)

func TestBuildModules(t *testing.T) {
	modules, err := BuildModules(&infra.Clients{})
	if err != nil {
		t.Fatalf("build modules returned error: %v", err)
	}

	if len(modules) == 0 {
		t.Fatal("expected at least one module to be built")
	}
}
