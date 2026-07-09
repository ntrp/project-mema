package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Media Manager API</title>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
	<style>
		body { margin: 0; background: #fff; }
		.swagger-ui .topbar { display: none; }
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
	<script>
		window.addEventListener("load", function () {
			SwaggerUIBundle({
				url: "/api/openapi.yaml",
				dom_id: "#swagger-ui",
				deepLinking: true,
				persistAuthorization: true,
				displayRequestDuration: true
			});
		});
	</script>
</body>
</html>`

func SwaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerHTML))
}

func OpenAPISpecHandler(w http.ResponseWriter, r *http.Request) {
	path, ok := openAPISpecPath()
	if !ok {
		http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	http.ServeFile(w, r, path)
}

func openAPISpecPath() (string, bool) {
	candidates := []string{filepath.Join("api", "openapi.yaml")}
	if _, file, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(file), "..", "..", "api", "openapi.yaml"))
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}
	return "", false
}
