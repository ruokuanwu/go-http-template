package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

type failedStore struct{}

func (failedStore) Find(string) ([]byte, bool, error)      { return nil, false, nil }
func (failedStore) Commit(string, []byte, time.Time) error { return errors.New("database unavailable") }
func (failedStore) Delete(string) error                    { return nil }

func TestSessionCommitFailureDoesNotPublishSuccess(t *testing.T) {
	sm := scs.New()
	sm.Store = failedStore{}
	sm.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`{"code":10000,"message":"session failed"}`))
	}
	handler := Session(sm, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm.Put(r.Context(), "userID", uint(1))
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"code":0,"data":{"user":1}}`))
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("POST", "/api/auth/register", nil))
	if w.Code != 500 || strings.Contains(w.Body.String(), `"code":0`) || len(w.Result().Cookies()) != 0 {
		t.Fatalf("commit failure leaked success: %d %s", w.Code, w.Body.String())
	}
}
