package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/skvoch/burst/internal/app/model"
	"github.com/skvoch/burst/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServerHandleTypesCreateAndGet(t *testing.T) {
	log := logrus.New()
	s := newServer(teststore.New(), log)

	types := make([]*model.Type, 0)

	types = append(types, &model.Type{ID: 0, Name: "C++ books"})
	types = append(types, &model.Type{ID: 1, Name: "C# books"})
	types = append(types, &model.Type{ID: 2, Name: "Go books"})
	types = append(types, &model.Type{ID: 3, Name: "Math books"})

	for _, _type := range types {
		json, err := json.Marshal(_type)
		assert.NoError(t, err)

		reader := bytes.NewReader(json)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1.0/types/create/", reader)

		s.ServeHTTP(rec, req)

		assert.Equal(t, rec.Code, http.StatusCreated)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1.0/types/", nil)

	s.ServeHTTP(rec, req)
	foundTypes := make([]*model.Type, 0)
	json.Unmarshal(rec.Body.Bytes(), &foundTypes)

	for index := 0; index < len(types); index++ {
		assert.Equal(t, types[index].ID, foundTypes[index].ID)
		assert.Equal(t, types[index].Name, foundTypes[index].Name)
	}
}

func TestServerHandleGetBooksIDs(t *testing.T) {
	type Response struct {
		BooksIDs []int `json:"books_ids"`
	}

	log := logrus.New()
	s := newServer(teststore.New(), log)

	s.store.Types().Create(&model.Type{ID: 0, Name: "Go books"})
	s.store.Types().Create(&model.Type{ID: 1, Name: "C++ books"})

	s.store.Books().Create(model.NewTestBookWithType(0))
	s.store.Books().Create(model.NewTestBookWithType(0))
	s.store.Books().Create(model.NewTestBookWithType(0))
	s.store.Books().Create(model.NewTestBookWithType(0))
	s.store.Books().Create(model.NewTestBookWithType(0))

	s.store.Books().Create(model.NewTestBookWithType(1))
	s.store.Books().Create(model.NewTestBookWithType(1))
	s.store.Books().Create(model.NewTestBookWithType(1))
	s.store.Books().Create(model.NewTestBookWithType(1))
	s.store.Books().Create(model.NewTestBookWithType(1))

	check := func(t *testing.T, typeID int, expectedLength int) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1.0/types/"+strconv.Itoa(typeID)+"/books/", nil)
		s.ServeHTTP(rec, req)
		assert.Equal(t, rec.Code, http.StatusOK)

		response := Response{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, expectedLength, len(response.BooksIDs))
	}

	check(t, 0, 5)
	check(t, 1, 5)
}
