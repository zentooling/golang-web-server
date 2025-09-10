// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"

	// "github.com/uberswe/golang-base-project/config"

	"github.com/uberswe/golang-base-project/middleware"
)

// PageData holds the default data needed for HTML pages to render
type PageData struct {
	Title           string
	Messages        []Message
	IsAuthenticated bool
	CacheParameter  string
	Trans           func(s string) string
}

// Message holds a message which can be rendered as responses on HTML pages
type Message struct {
	Type    string // success, warning, error, etc.
	Content string
}

// isAuthenticated checks if the current user is authenticated or not
func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	return exists
}

func DefaultPageData(c *gin.Context) PageData {
	langService := infra.NewService(c, infra.LairInstance().GetBundle())
	return PageData{
		Title:           "Home",
		Messages:        nil,
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  infra.LairInstance().GetConfig().CacheParameter,
		Trans:           langService.Trans,
	}
}
