package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/proto"
)

func testState(t *testing.T) (*State, func()) {
	tmpDir, err := os.MkdirTemp("/tmp", "atlast-state-")
	if err != nil {
		t.Fatal(err)
	}
	st, err := NewState(filepath.Join(tmpDir, "my.db"))
	if err != nil {
		t.Fatal(err)
	}
	closeFn := func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatal(err)
		}
	}
	return st, closeFn
}

func TestState_Provider(t *testing.T) {
	state, closeFn := testState(t)
	defer closeFn()

	providers, err := state.ListProviders()
	assert.NoError(t, err)
	assert.Empty(t, providers)

	p := &proto.Provider{
		Id:       "a",
		Name:     "b",
		Provider: "c",
	}
	err = state.CreateProvider(p)
	assert.NoError(t, err)

	providers, err = state.ListProviders()
	assert.NoError(t, err)
	assert.Equal(t, providers[0].Id, "a")
}
