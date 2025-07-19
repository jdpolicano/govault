package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/jdpolicano/govault/internal/server"
)

func TestJSONResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	server.JSONResponse(rec, server.NewResponse(200, "ok", nil))

	if rec.Code != 200 {
		t.Fatalf("expected status 200 got %d", rec.Code)
	}
	var res server.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if res.Data != "ok" || res.Error != "" || res.Code != 200 {
		t.Errorf("unexpected response %+v", res)
	}
}
