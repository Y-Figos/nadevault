package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Y-Figos/nadevault/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeRepo struct {
	getNadeByIDFn       func(ctx context.Context, id int64) (*domain.Nade, error)
	listNadesByMapIDFn  func(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error)
	getMapByCodeFn      func(ctx context.Context, code string) (*domain.CsMap, error)
	addNadeFn           func(ctx context.Context, nade domain.Nade) (string, error)
	getnadeByPublicIDFn func(ctx context.Context, publicID string) (*domain.Nade, error)
}

func (f fakeRepo) GetNadeByID(ctx context.Context, id int64) (*domain.Nade, error) {
	return f.getNadeByIDFn(ctx, id)
}

func (f fakeRepo) ListNadesByMapID(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error) {
	return f.listNadesByMapIDFn(ctx, mapID, limit, offset)
}

func (f fakeRepo) GetMapByCode(ctx context.Context, code string) (*domain.CsMap, error) {
	return f.getMapByCodeFn(ctx, code)
}

func (f fakeRepo) AddNade(ctx context.Context, nade domain.Nade) (string, error) {
	return f.addNadeFn(ctx, nade)
}

func (f fakeRepo) GetNadeByPublicID(ctx context.Context, publicID string) (*domain.Nade, error) {
	return f.getnadeByPublicIDFn(ctx, publicID)
}
func (f fakeRepo) GetMapList(ctx context.Context) ([]domain.CsMap, error) {
	panic("not implemented")
}

func reqWithURLParam(req *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func readErrorCode(t *testing.T, body []byte) string {
	t.Helper()
	var resp struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("invalid json response: %v body=%s", err, string(body))
	}
	return resp.Error.Code
}
func TestGetNadeByID(t *testing.T) {

	t.Run("invalid param returns 400 invalid_param", func(t *testing.T) {

		h := NewNadeHandler(fakeRepo{
			getNadeByIDFn: func(ctx context.Context, id int64) (*domain.Nade, error) {
				t.Fatal("repo should not be called when param is invalid")
				return nil, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/api/nades/abc", nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()

		rctx.URLParams.Add("nadeID", "abc")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetNadeByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		responseCode := readErrorCode(t, w.Body.Bytes())
		if responseCode != "invalid_param" {
			t.Fatalf("expected error.code=%q got %q", "invalid_param", responseCode)
		}

	})
	t.Run("not found returns 404 not_found", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			getNadeByIDFn: func(ctx context.Context, id int64) (*domain.Nade, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/api/nades/123", nil)
		w := httptest.NewRecorder()

		req = reqWithURLParam(req, "nadeID", "123")

		h.GetNadeByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected %d got %d body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
		responseCode := readErrorCode(t, w.Body.Bytes())
		if responseCode != "not_found" {
			t.Fatalf("expected error.code=%q got %q", "not_found", responseCode)
		}
	})
	t.Run("success returns 200 and nade json", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			getNadeByIDFn: func(ctx context.Context, id int64) (*domain.Nade, error) {
				return &domain.Nade{ID: id, Name: "test"}, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/api/nades/123", nil)
		w := httptest.NewRecorder()

		req = reqWithURLParam(req, "nadeID", "123")
		h.GetNadeByID(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}
		var resp struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json response: %v body=%s", err, w.Body.String())
		}
		if resp.ID != 123 || resp.Name != "test" {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})
}

func TestAddNades(t *testing.T) {
	t.Run("invalid json returns 400 invalid_json", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			addNadeFn: func(ctx context.Context, nade domain.Nade) (string, error) {
				t.Fatal("repo should not be called when json is invalid")
				return "", nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/api/nades", nil)
		w := httptest.NewRecorder()

		h.AddNade(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		responseCode := readErrorCode(t, w.Body.Bytes())
		if responseCode != "invalid_json" {
			t.Fatalf("expected error.code=%q got %q", "invalid_json", responseCode)
		}
	})
	t.Run("unknown fields returns 400 invalid_json", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			addNadeFn: func(ctx context.Context, nade domain.Nade) (string, error) {
				t.Fatal("repo should not be called when json has unknown fields")
				return "", nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/nades", http.NoBody)
		req.Body = io.NopCloser(strings.NewReader(`{"name":"test","unknown_field":"value"}`))
		w := httptest.NewRecorder()

		h.AddNade(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
		responseCode := readErrorCode(t, w.Body.Bytes())
		if responseCode != "invalid_json" {
			t.Fatalf("expected error.code=%q got %q", "invalid_json", responseCode)
		}
	})
	t.Run("validation of fields returns 400 validation_error", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			addNadeFn: func(ctx context.Context, nade domain.Nade) (string, error) {
				t.Fatal("repo should not be called when validation fails")
				return "", nil
			},
		})

		body := `{
		"name": "Mirage Window Smoke",
		"desc": "test",
		"map_id": 1,
		"type": "smoke",
		"common_side": "Terror",
		"from": "T Spawn",
		"to": "Window",
		"mouse_click": "mouse1"
	}`

		req := httptest.NewRequest(http.MethodPost, "/api/nades", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.AddNade(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}

		var resp struct {
			Error struct {
				Code    string `json:"code"`
				Details []struct {
					Field string `json:"field"`
				} `json:"details"`
			} `json:"error"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json response: %v body=%s", err, w.Body.String())
		}

		if resp.Error.Code != "validation_error" {
			t.Fatalf("expected error.code=%q got %q", "validation_error", resp.Error.Code)
		}

		found := false
		for _, d := range resp.Error.Details {
			if d.Field == "common_side" {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("expected validation error for field common_side")
		}
	})
	t.Run("success returns 201 with id", func(t *testing.T) {
		called := false

		h := NewNadeHandler(fakeRepo{
			addNadeFn: func(ctx context.Context, nade domain.Nade) (string, error) {
				called = true

				// checa só o essencial (contrato DTO -> domain)
				if nade.Name != "Mirage Window Smoke" {
					t.Fatalf("expected name %q got %q", "Mirage Window Smoke", nade.Name)
				}
				if nade.MapID != 1 {
					t.Fatalf("expected map_id %d got %d", 1, nade.MapID)
				}
				if string(nade.CommonSide) != "T" {
					t.Fatalf("expected common_side %q got %q", "T", nade.CommonSide)
				}
				if string(nade.Type) != "smoke" {
					t.Fatalf("expected type %q got %q", "smoke", nade.Type)
				}

				return "123", nil
			},
		})

		body := `{
		"name": "Mirage Window Smoke",
		"desc": "Easy instant smoke from T Spawn",
		"map_id": 1,
		"type": "smoke",
		"common_side": "T",
		"from": "T Spawn",
		"to": "Window",
		"mouse_click": "mouse1",
		"is_jumping": true,
		"is_running": false,
		"is_walking": false,
		"images": {},
		"created_by": "user_123"
	}`

		req := httptest.NewRequest(http.MethodPost, "/api/nades", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.AddNade(w, req)

		if !called {
			t.Fatalf("expected repo.AddNade to be called")
		}

		if w.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var resp struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json response: %v body=%s", err, w.Body.String())
		}
		if resp.ID != 123 {
			t.Fatalf("expected id=%d got %d", 123, resp.ID)
		}
	})
	t.Run("conflict returns 409 conflict", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			addNadeFn: func(ctx context.Context, nade domain.Nade) (string, error) {
				return "", &pgconn.PgError{Code: "23505"}
			},
		})

		body := `{
		"name": "Mirage Window Smoke",
		"desc": "Easy instant smoke from T Spawn",
		"map_id": 1,
		"type": "smoke",
		"common_side": "T",
		"from": "T Spawn",
		"to": "Window",
		"mouse_click": "mouse1",
		"images": {},
		"created_by": "user_123"
	}`

		req := httptest.NewRequest(http.MethodPost, "/api/nades", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.AddNade(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected %d got %d body=%s", http.StatusConflict, w.Code, w.Body.String())
		}

		if got := readErrorCode(t, w.Body.Bytes()); got != "conflict" {
			t.Fatalf("expected error.code=%q got %q", "conflict", got)
		}
	})

}
func TestListNadesByMapID(t *testing.T) {
	t.Run("map not found returns 404 not_found", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			getMapByCodeFn: func(ctx context.Context, code string) (*domain.CsMap, error) {
				return nil, pgx.ErrNoRows
			},
			listNadesByMapIDFn: func(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error) {
				t.Fatal("list should not be called when map not found")
				return nil, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/api/maps/mirage/nades", nil)
		req = reqWithURLParam(req, "mapCode", "mirage")
		w := httptest.NewRecorder()

		h.ListNadesByMapID(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected %d got %d body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
		if got := readErrorCode(t, w.Body.Bytes()); got != "not_found" {
			t.Fatalf("expected error.code=%q got %q", "not_found", got)
		}
	})
	t.Run("success returns 200 and list", func(t *testing.T) {
		h := NewNadeHandler(fakeRepo{
			getMapByCodeFn: func(ctx context.Context, code string) (*domain.CsMap, error) {
				return &domain.CsMap{ID: 1, Code: code, DisplayName: "Mirage"}, nil
			},
			listNadesByMapIDFn: func(ctx context.Context, mapID int16, limit int32, offset int32) ([]domain.Nade, error) {
				return []domain.Nade{
					{ID: 1, Name: "A"},
					{ID: 2, Name: "B"},
				}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/api/maps/mirage/nades?limit=20&offset=0", nil)
		req = reqWithURLParam(req, "mapCode", "mirage")
		w := httptest.NewRecorder()

		h.ListNadesByMapID(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var resp []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid json response: %v body=%s", err, w.Body.String())
		}
		if len(resp) != 2 {
			t.Fatalf("expected 2 items got %d", len(resp))
		}
	})

}
