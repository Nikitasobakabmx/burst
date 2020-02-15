package store_test


import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skvoch/burst/internal/app/model"
	"github.com/skvoch/burst/internal/app/store"
)

func TestTypeRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("books")

	repo := s.Types()

	_type := &model.Type {
		ID:0,
		Name: "Nothing",
	}
	
	err := repo.Create(_type)
	assert.NoError(t, err)
}