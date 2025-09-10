package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
)

type ConfigPageData struct {
	PageData
	Config   *infra.Config
	LogLevel string
}

func ConfigRouteHandler(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Configuration")

	currentLevel := infra.LairInstance().GetLoggingLevel()
	slog.Error("ConfigRouteHandler", "current logging Level", currentLevel)

	cd := ConfigPageData{
		PageData: pd,
		Config:   infra.LairInstance().GetConfig(),
		LogLevel: currentLevel.String(),
	}

	c.HTML(http.StatusOK, "config.html", cd)
}

// LoggingRouteHandlerPost handles change in logging submitted by user
func LoggingRouteHandlerPost(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Configuration")

	levelStr := c.PostForm("log-select")

	level, err := infra.StringToLevel(levelStr)

	if err == nil {
		infra.LairInstance().GetLoggingLevel().Set(level)
	}

	slog.Error("set logging level to " + levelStr)
	slog.Warn("set logging level to " + levelStr)
	slog.Info("set logging level to " + levelStr)
	slog.Debug("set logging level to " + levelStr)

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pd.Trans("Logging level set to " + levelStr),
	})
	cd := ConfigPageData{
		PageData: pd,
		Config:   infra.LairInstance().GetConfig(),
		LogLevel: levelStr,
	}

	c.HTML(http.StatusOK, "config.html", cd)
}

func ConfigRouteHandlerPost(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Configuration")

	prevCfg := infra.LairInstance().GetConfig()

	newCookie := c.PostForm("cookie-secret")

	slog.Info("ConfigRouteHandlerPost", "prevCookie", prevCfg.CookieSecret)
	if newCookie != prevCfg.CookieSecret {
		slog.Info("ConfigRouteHandlerPost", "newCookie", newCookie)
		prevCfg.CookieSecret = newCookie
	}

	cd := ConfigPageData{
		PageData: pd,
		Config:   infra.LairInstance().GetConfig(),
	}

	c.HTML(http.StatusOK, "config.html", cd)
}
