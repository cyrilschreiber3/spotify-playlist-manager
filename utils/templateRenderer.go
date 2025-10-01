package utils

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RenderTemplate(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}
