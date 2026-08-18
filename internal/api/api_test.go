package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	scaffoldgoapi "github.com/alrayyes/scaffold-go-api/internal/api"
	"github.com/stretchr/testify/require"
)

func TestHealthAnswersOK(t *testing.T) {
	mux := scaffoldgoapi.NewMux(nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body scaffoldgoapi.Health
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "ok", body.Status)
}

func TestGetWidget(t *testing.T) {
	widgets := map[string]scaffoldgoapi.Widget{
		"hammer": {ID: "hammer", Name: "Claw hammer"},
	}
	mux := scaffoldgoapi.NewMux(widgets)

	t.Run("known id returns the widget", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/widgets/hammer", nil)

		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var body scaffoldgoapi.Widget
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, scaffoldgoapi.Widget{ID: "hammer", Name: "Claw hammer"}, body)
	})

	t.Run("unknown id returns 404 with an error body", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/widgets/nope", nil)

		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var body scaffoldgoapi.Error
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, "unknown widget", body.Error)
	})
}
