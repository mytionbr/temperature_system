package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jarcoal/httpmock"
)

func newChiRequest(method, path, cep string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", cep)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestCEPHandler_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		"https://viacep.com.br/ws/01001000/json",
		httpmock.NewStringResponder(200, `{
			"cep":"01001-000",
			"logradouro":"Praça da Sé",
			"complemento":"lado ímpar",
			"bairro":"Sé",
			"localidade":"São Paulo",
			"uf":"SP",
			"ibge":"3550308",
			"gia":"1004",
			"ddd":"11",
			"siafi":"7107"
		}`),
	)

	req := newChiRequest("GET", "/cep/01001000", "01001000")
	rec := httptest.NewRecorder()

	CEPHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("wanted 200, got %d. body=%s", rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("error parsing json: %v", err)
	}
	if got := body["localidade"]; got != "São Paulo" {
		t.Fatalf("wanted localidade=São Paulo, got %v", got)
	}
}
func TestCEPHandler_InvalidCEP(t *testing.T) {
	req := newChiRequest("GET", "/cep/123", "123")
	rec := httptest.NewRecorder()

	CEPHandler(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("wanted 422, got %d. body=%s", rec.Code, rec.Body.String())
	}
}

func TestCEPHandler_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		"https://viacep.com.br/ws/99999999/json",
		httpmock.NewStringResponder(400, `{"erro": "true"}`),
	)

	req := newChiRequest("GET", "/cep/99999999", "99999999")
	rec := httptest.NewRecorder()

	CEPHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d. body=%s", rec.Code, rec.Body.String())
	}
}
