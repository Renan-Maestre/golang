package contexto

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)

	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()

	if !s.cancelled {
		s.t.Errorf("store não foi avisada para cancelar")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()

	if s.cancelled {
		s.t.Errorf("store foi avisada para cancelar")
	}
}

func TestServer(t *testing.T) {
	data := "olá, mundo"

	t.Run("retorna dados da store", func(t *testing.T) {
		store := &SpyStore{
			response: data,
		}

		server := Server(store)

		request := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(
				`resultado "%s", esperado "%s"`,
				response.Body.String(),
				data,
			)
		}

		store.assertWasCancelled()
	})

	t.Run("cancela a store se a requisição for cancelada", func(t *testing.T) {
		store := &SpyStore{
			response: data,
		}

		server := Server(store)

		request := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		cancellingCtx, cancel := context.WithCancel(
			request.Context(),
		)

		time.AfterFunc(
			5*time.Millisecond,
			cancel,
		)

		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		store.assertWasNotCancelled()
	})
}
