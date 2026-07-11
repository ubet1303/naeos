package monitoring

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("expected registry to be created")
	}
}

func TestRegistryRegister(t *testing.T) {
	reg := NewRegistry()
	reg.Register("test_counter", Counter, "A test counter")

	families := reg.GetFamilies()
	if len(families) != 1 {
		t.Errorf("expected 1 family, got %d", len(families))
	}
	if families[0].Name != "test_counter" {
		t.Errorf("expected name 'test_counter', got %s", families[0].Name)
	}
}

func TestRegistryCounterInc(t *testing.T) {
	reg := NewRegistry()
	reg.Register("requests", Counter, "Total requests")

	reg.CounterInc("requests", nil)
	reg.CounterInc("requests", nil)

	families := reg.GetFamilies()
	if len(families) != 1 {
		t.Fatalf("expected 1 family, got %d", len(families))
	}

	if families[0].Metrics[0].Value != 2 {
		t.Errorf("expected value 2, got %f", families[0].Metrics[0].Value)
	}
}

func TestRegistryCounterWithLabels(t *testing.T) {
	reg := NewRegistry()
	reg.Register("http_requests", Counter, "HTTP requests")

	labels1 := map[string]string{"method": "GET", "path": "/api"}
	labels2 := map[string]string{"method": "POST", "path": "/api"}

	reg.CounterInc("http_requests", labels1)
	reg.CounterInc("http_requests", labels1)
	reg.CounterInc("http_requests", labels2)

	families := reg.GetFamilies()
	if len(families[0].Metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(families[0].Metrics))
	}
}

func TestRegistryGauge(t *testing.T) {
	reg := NewRegistry()
	reg.Register("connections", Gauge, "Active connections")

	reg.GaugeSet("connections", 10, nil)
	reg.GaugeInc("connections", nil)
	reg.GaugeDec("connections", nil)

	families := reg.GetFamilies()
	if families[0].Metrics[0].Value != 10 {
		t.Errorf("expected value 10, got %f", families[0].Metrics[0].Value)
	}
}

func TestRegistryHistogram(t *testing.T) {
	reg := NewRegistry()
	reg.Register("duration", Histogram, "Request duration")

	reg.HistogramObserve("duration", 0.1, nil)
	reg.HistogramObserve("duration", 0.2, nil)

	families := reg.GetFamilies()
	if len(families[0].Metrics) != 2 {
		t.Errorf("expected 2 observations, got %d", len(families[0].Metrics))
	}
}

func TestFormatPrometheus(t *testing.T) {
	reg := NewRegistry()
	reg.Register("test_metric", Counter, "A test metric")
	reg.CounterInc("test_metric", nil)

	output := reg.FormatPrometheus()
	if !strings.Contains(output, "# HELP test_metric") {
		t.Error("expected HELP line")
	}
	if !strings.Contains(output, "# TYPE test_metric counter") {
		t.Error("expected TYPE line")
	}
	if !strings.Contains(output, "test_metric") {
		t.Error("expected metric name")
	}
}

func TestFormatLabels(t *testing.T) {
	labels := map[string]string{"method": "GET", "status": "200"}
	result := formatLabels(labels)

	if !strings.Contains(result, "method=\"GET\"") {
		t.Error("expected method label")
	}
	if !strings.Contains(result, "status=\"200\"") {
		t.Error("expected status label")
	}
}

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics()
	if metrics == nil {
		t.Fatal("expected metrics to be created")
	}

	families := metrics.Registry().GetFamilies()
	if len(families) < 5 {
		t.Errorf("expected at least 5 default metrics, got %d", len(families))
	}
}

func TestMetricsIncRequests(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncRequests("GET", "/api/health", "200")

	families := metrics.Registry().GetFamilies()
	for _, f := range families {
		if f.Name == "naeos_requests_total" {
			if len(f.Metrics) == 0 {
				t.Error("expected at least 1 metric")
			}
			return
		}
	}
	t.Error("expected naeos_requests_total metric")
}

func TestMetricsIncPipelines(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncPipelines("success")
	metrics.IncPipelines("failure")

	families := metrics.Registry().GetFamilies()
	for _, f := range families {
		if f.Name == "naeos_pipelines_total" {
			if len(f.Metrics) != 2 {
				t.Errorf("expected 2 metrics, got %d", len(f.Metrics))
			}
			return
		}
	}
	t.Error("expected naeos_pipelines_total metric")
}

func TestHealthHandler(t *testing.T) {
	handler := HealthHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "healthy") {
		t.Error("expected 'healthy' in response")
	}
}

func TestReadyHandler(t *testing.T) {
	handler := ReadyHandler()

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "ready") {
		t.Error("expected 'ready' in response")
	}
}

func TestPrometheusHandler(t *testing.T) {
	reg := NewRegistry()
	reg.Register("test_metric", Counter, "Test")
	reg.CounterInc("test_metric", nil)

	handler := PrometheusHandler(reg)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test_metric") {
		t.Error("expected metric in response")
	}
}

func TestCollector(t *testing.T) {
	reg := NewRegistry()
	reg.Register("test", Counter, "Test")
	reg.CounterInc("test", nil)

	collector := NewCollector(reg)
	snapshot := collector.Collect()

	if snapshot == nil {
		t.Fatal("expected snapshot")
	}

	if len(snapshot.Families) != 1 {
		t.Errorf("expected 1 family, got %d", len(snapshot.Families))
	}

	if snapshot.Timestamp.IsZero() {
		t.Error("expected timestamp")
	}
}

func TestMetricsMiddleware(t *testing.T) {
	metrics := NewMetrics()
	middleware := MetricsMiddleware(metrics)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(inner)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUptime(t *testing.T) {
	metrics := NewMetrics()
	start := time.Now()
	uptime := time.Since(start).Seconds()
	metrics.SetUptime(uptime)

	families := metrics.Registry().GetFamilies()
	for _, f := range families {
		if f.Name == "naeos_uptime_seconds" {
			if f.Metrics[0].Value <= 0 {
				t.Error("expected positive uptime")
			}
			return
		}
	}
	t.Error("expected naeos_uptime_seconds metric")
}
