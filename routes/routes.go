// Package routes defines all the handling functions for all the routes
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

// Define an enum using iota
type MessageType int

const (
	Error MessageType = iota
	Warning
	Info
	Success
)

// Map for string representation
var messageTypeNames = map[MessageType]string{
	Error:   "error",
	Warning: "warning",
	Info:    "info",
	Success: "success",
}

func (pd *PageData) AddMessage(msgType MessageType, content string) {
	pd.Messages = append(pd.Messages, Message{
		Type:    msgType,
		Content: content,
	})
}

func (c MessageType) String() string {
	if name, ok := messageTypeNames[c]; ok {
		return name
	}
	return "Unknown"
}

// Message holds a message which can be rendered as responses on HTML pages
type Message struct {
	Type    MessageType // success, warning, error, etc.
	Content string
}

// isAuthenticated checks if the current user is authenticated or not
func isAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(middleware.UserIDKey)
	return exists
}

func getUserId(c *gin.Context) uint {
	id, exists := c.Get(middleware.UserIDKey)
	if exists {
		return id.(uint)
	}
	return 0
}

func DefaultPageData(c *gin.Context, bundle *i18n.Bundle, cacheParameter string) PageData {
	langService := infra.NewLangService(c, bundle)
	return PageData{
		Title:           "Home",
		Messages:        nil,
		IsAuthenticated: isAuthenticated(c),
		CacheParameter:  cacheParameter,
		Trans:           langService.Trans,
	}
}

type Service struct {
	env infra.ILair
}

func NewService(env infra.ILair) *Service {
	return &Service{
		env: env,
	}
}
