package routes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
)

type ConfigPageData struct {
	PageData
	Config   *infra.Config
	LogLevel string
}

func pageData(c *gin.Context) *ConfigPageData {
	return &ConfigPageData{
		PageData: DefaultPageData(c),
		Config:   infra.LairInstance().GetConfig(),
		LogLevel: infra.LairInstance().GetLoggingLevel().Level().String(),
	}
}
func ConfigRouteHandler(c *gin.Context) {
	pd := pageData(c)
	pd.Title = pd.Trans("Configuration")

	currentLevel := infra.LairInstance().GetLoggingLevel()
	slog.Error("ConfigRouteHandler", "current logging Level", currentLevel)

	c.HTML(http.StatusOK, "config.gohtml", pd)
}

// LoggingRouteHandlerPost handles change in logging submitted by user
func LoggingRouteHandlerPost(c *gin.Context) {

	levelStr := c.PostForm("log-select")

	level, err := infra.StringToLevel(levelStr)

	if err == nil {
		infra.LairInstance().GetLoggingLevel().Set(level)
	}

	slog.Error("set logging level to " + levelStr)
	slog.Warn("set logging level to " + levelStr)
	slog.Info("set logging level to " + levelStr)
	slog.Debug("set logging level to " + levelStr)

	// read updated state
	pd := pageData(c)
	pd.Title = pd.Trans("Configuration")
	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pd.Trans("Logging level set to " + levelStr),
	})

	c.HTML(http.StatusOK, "config.gohtml", pd)
}

func ConfigRouteHandlerPost(c *gin.Context) {

	// this points to config struct so changes are for duration of this process
	prevCfg := infra.LairInstance().GetConfig()

	dirty := false

	newValue := c.PostForm("base_url")
	if newValue != prevCfg.BaseURL {
		slog.Info("BaseUrl", "newValue", newValue)
		prevCfg.BaseURL = newValue
		dirty = true
	}
	newValue = c.PostForm("smtp_host")
	if newValue != prevCfg.SMTPHost {
		slog.Info("SmtpHost", "newValue", newValue)
		prevCfg.SMTPHost = newValue
		dirty = true
	}
	newValue = c.PostForm("smtp_port")
	if newValue != prevCfg.SMTPPort {
		slog.Info("SMTPPort", "newValue", newValue)
		prevCfg.SMTPPort = newValue
		dirty = true
	}
	newValue = c.PostForm("smtp_sender")
	if newValue != prevCfg.SMTPSender {
		slog.Info("SMTPSender", "newValue", newValue)
		prevCfg.SMTPSender = newValue
		dirty = true
	}
	newValue = c.PostForm("smtp_username")
	if newValue != prevCfg.SMTPUsername {
		slog.Info("SMTPUsername", "newValue", newValue)
		prevCfg.SMTPUsername = newValue
		dirty = true
	}
	newValue = c.PostForm("smtp_password")
	if newValue != prevCfg.SMTPPassword {
		slog.Info("SMTPPassword", "newValue", newValue)
		prevCfg.SMTPPassword = newValue
		dirty = true
	}
	newValue = c.PostForm("request_per_minute")
	newValueInt, err := strconv.Atoi(newValue)

	messages := make([]Message, 0)

	if err != nil {
		messages = append(messages, Message{
			Type:    "error",
			Content: err.Error(),
		})
	}

	if err == nil && newValueInt != prevCfg.RequestsPerMinute {
		slog.Info("RequestsPerMinute", "newValue", newValue)
		prevCfg.RequestsPerMinute = newValueInt
		dirty = true
	}
	newValue = c.PostForm("cache_parameter")
	if newValue != prevCfg.CacheParameter {
		slog.Info("CacheParameter", "newValue", newValue)
		prevCfg.CacheParameter = newValue
		dirty = true
	}
	newValue = c.PostForm("cache_max_age")
	// reuse from above
	newValueInt, err = strconv.Atoi(newValue)
	if err != nil {
		messages = append(messages, Message{
			Type:    "error",
			Content: err.Error(),
		})
	}
	if err == nil && newValueInt != prevCfg.CacheMaxAge {
		slog.Info("CacheMaxAge", "newValue", newValue)
		prevCfg.CacheMaxAge = newValueInt
		dirty = true
	}
	if dirty {
		messages = append(messages, Message{
			Type:    "success",
			Content: "Environment updated",
		})
	}

	// read updated state
	pd := pageData(c)
	pd.Title = pd.Trans("Configuration")
	// add collected messages along the way
	pd.Messages = messages

	c.HTML(http.StatusOK, "config.gohtml", pd)
}
