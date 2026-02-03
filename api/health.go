package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// ReadinessResponse represents the readiness check response
type ReadinessResponse struct {
	Status    string                `json:"status"`
	Timestamp time.Time             `json:"timestamp"`
	Checks    map[string]CheckState `json:"checks"`
}

// CheckState represents the state of an individual health check
type CheckState struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// MetricsResponse represents Prometheus-style metrics
type MetricsResponse struct {
	GoVersion      string  `json:"go_version"`
	GoRoutines     int     `json:"goroutines"`
	MemoryAlloc    uint64  `json:"memory_alloc_bytes"`
	MemorySys      uint64  `json:"memory_sys_bytes"`
	MemoryTotal    uint64  `json:"memory_total_alloc_bytes"`
	GCCycles       uint32  `json:"gc_cycles"`
	LastGC         uint64  `json:"last_gc_ns"`
	Uptime         float64 `json:"uptime_seconds"`
}

var (
	// startTime records when the server started
	startTime = time.Now()
	// version can be set at build time
	version = "dev"
)

// HealthHandler handles GET /health
// Simple health check that returns 200 OK if the server is running
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(response)
	}
}

// ReadyHandler handles GET /ready
// Readiness check that verifies the application can serve traffic
// Checks database connectivity and other critical dependencies
func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checks := make(map[string]CheckState)
	allHealthy := true

	// Check database connectivity
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "ok" {
		allHealthy = false
	}

	// Determine overall status
	status := "ready"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	response := ReadinessResponse{
		Status:    status,
		Timestamp: time.Now(),
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// checkDatabase verifies database connectivity
func (h *Handler) checkDatabase() CheckState {
	// Try to perform a simple database operation
	// For SQLite, we can try to list users with a limit
	_, err := h.db.ListUsers()
	if err != nil {
		return CheckState{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	return CheckState{
		Status:  "ok",
		Message: "database connection healthy",
	}
}

// MetricsHandler handles GET /metrics
// Returns application metrics in a format compatible with Prometheus
func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate uptime
	uptime := time.Since(startTime).Seconds()

	// For JSON response
	if r.Header.Get("Accept") == "application/json" {
		response := MetricsResponse{
			GoVersion:      runtime.Version(),
			GoRoutines:     runtime.NumGoroutine(),
			MemoryAlloc:    m.Alloc,
			MemorySys:      m.Sys,
			MemoryTotal:    m.TotalAlloc,
			GCCycles:       m.NumGC,
			LastGC:         m.LastGC,
			Uptime:         uptime,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Default Prometheus text format
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	// Write Prometheus metrics
	metrics := `# HELP athena_go_version Go version info
# TYPE athena_go_version gauge
athena_go_version{version="` + runtime.Version() + `"} 1

# HELP athena_goroutines Number of goroutines
# TYPE athena_goroutines gauge
athena_goroutines ` + formatUint64(uint64(runtime.NumGoroutine())) + `

# HELP athena_memory_alloc_bytes Bytes allocated and in use
# TYPE athena_memory_alloc_bytes gauge
athena_memory_alloc_bytes ` + formatUint64(m.Alloc) + `

# HELP athena_memory_sys_bytes Bytes obtained from system
# TYPE athena_memory_sys_bytes gauge
athena_memory_sys_bytes ` + formatUint64(m.Sys) + `

# HELP athena_memory_total_alloc_bytes Total bytes allocated (even if freed)
# TYPE athena_memory_total_alloc_bytes counter
athena_memory_total_alloc_bytes ` + formatUint64(m.TotalAlloc) + `

# HELP athena_gc_cycles_total Number of completed GC cycles
# TYPE athena_gc_cycles_total counter
athena_gc_cycles_total ` + formatUint64(uint64(m.NumGC)) + `

# HELP athena_uptime_seconds Time since server started
# TYPE athena_uptime_seconds gauge
athena_uptime_seconds ` + formatFloat64(uptime) + `

# HELP athena_build_info Build information
# TYPE athena_build_info gauge
athena_build_info{version="` + version + `"} 1
`

	w.Write([]byte(metrics))
}

// Helper functions for formatting metrics
func formatUint64(v uint64) string {
	return fmt.Sprintf("%d", v)
}

func formatFloat64(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// SetVersion sets the application version (typically called during initialization)
func SetVersion(v string) {
	version = v
}
