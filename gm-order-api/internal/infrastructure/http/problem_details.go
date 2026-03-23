package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func NewProblemDetails(c *gin.Context, status int, title, detail string) {
	c.JSON(status, ProblemDetails{
		Type:     "https://httpstatuses.io/" + http.StatusText(status),
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	})
}
