package routes

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
)

func Loglevel(c *gin.Context) {

	pd := DefaultPageData(c)

	token := strings.ToLower(c.Param("level"))
	slog.Debug("Loglevel", "level", token)
	newLevel := slog.LevelDebug
	switch token {
	case "": // no token passed in url
		{
			newLevel = slog.LevelDebug
		}
	case "debug":
		{
			newLevel = slog.LevelDebug
		}
	case "error":
		{
			newLevel = slog.LevelError
		}
	case "info":
		{
			newLevel = slog.LevelInfo
		}
	case "warn":
		{
			newLevel = slog.LevelWarn
		}
	default:
		token = "debug"
		newLevel = slog.LevelDebug
	}
	slog.SetLogLoggerLevel(newLevel)
	//pdPre.Title = pdPre.Trans("Reset Password")
	//pd := ResetPasswordPageData{
	//	PageData: pdPre,
	//	Token:    token,
	//}

	// set default logger level - TODO make this configurable
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: newLevel,
	}))
	// Set this logger as the default overwriting previous one
	slog.SetDefault(logger)

	// now send a message at each level to indicate operation

	slog.Info("Loglevel", "new-level", newLevel)
	slog.Debug("log at debug level")
	slog.Info("log at info level")
	slog.Error("log at error level")
	slog.Warn("log at warn level")
	slog.Debug("Run", "config", infra.LairInstance().GetConfig())

	pd.Messages = append(pd.Messages, Message{
		Type:    "success",
		Content: pd.Trans("Logging level set to ") + token,
	})

	c.HTML(http.StatusOK, "loglevel.gohtml", pd)
}
