package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/config"
	"media-manager/internal/testdb"
)

func TestOpenDatabaseCreatesRiverSchema(t *testing.T) {
	cfg := config.Config{
		DatabaseURL:   testdb.Create(t),
		AdminUsername: "admin",
		AdminPassword: "password",
	}
	for range 2 {
		pool, err := openDatabase(context.Background(), cfg)
		if err != nil {
			t.Fatal(err)
		}
		var jobsExist, queuesExist bool
		err = pool.QueryRow(context.Background(), `
			select to_regclass('river_job') is not null, to_regclass('river_queue') is not null
		`).Scan(&jobsExist, &queuesExist)
		pool.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !jobsExist || !queuesExist {
			t.Fatalf("River tables missing: river_job=%v river_queue=%v", jobsExist, queuesExist)
		}
	}
}

func TestDevelopmentRouterDoesNotServeFrontendFallback(t *testing.T) {
	api := chi.NewRouter()
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := appRouter(config.Config{AppEnv: "development"}, api)

	apiResponse := httptest.NewRecorder()
	handler.ServeHTTP(apiResponse, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if apiResponse.Code != http.StatusNoContent {
		t.Fatalf("api status = %d", apiResponse.Code)
	}

	pageResponse := httptest.NewRecorder()
	handler.ServeHTTP(pageResponse, httptest.NewRequest(http.MethodGet, "/movies/example", nil))
	if pageResponse.Code != http.StatusNotFound {
		t.Fatalf("frontend fallback status = %d", pageResponse.Code)
	}
}
