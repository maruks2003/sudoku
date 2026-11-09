package docs

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

const page = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Sudoku API</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({ url: "/openapi.json", dom_id: "#swagger-ui" });
  </script>
</body>
</html>`

func Register(router gin.IRouter, spec *openapi3.T) error {
	specJSON, err := spec.MarshalJSON()
	if err != nil {
		return err
	}

	router.GET("/openapi.json", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "application/json", specJSON)
	})
	router.GET("/docs", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(page))
	})
	return nil
}
