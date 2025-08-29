package api

import (
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestApi struct {
}

// TestEndpoint is a test endpoint.
// swagger:operation GET /api/test TestEndpoint
// ---
// summary: Test endpoint
// description: This is a test endpoint to verify Swagger annotations.
// responses:
//   '200':
//     description: Test successful
//     schema:
//       "$ref": "#/definitions/Response"
func (t *TestApi) TestEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Test successful",
		Data: nil,
	})
}