package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

func TestAdminUsersRouteIsNotShadowedByStaticFiles(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := api.NewHandler(db)

	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true, CreatedAt: time.Now()}
	if err := db.CreateUser(admin); err != nil {
		t.Fatalf("failed to create admin: %v", err)
	}

	student := &models.User{Username: "student", Email: "student@test.com", CreatedAt: time.Now()}
	if err := db.CreateUser(student); err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	New(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected /admin/users to reach the admin route, got %d with body %q", rec.Code, rec.Body.String())
	}

	var users []*models.User
	if err := json.NewDecoder(rec.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode users response: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
