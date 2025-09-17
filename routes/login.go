package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
	"github.com/uberswe/golang-base-project/middleware"
	"github.com/uberswe/golang-base-project/models"
	"github.com/uberswe/golang-base-project/ulid"
	"golang.org/x/crypto/bcrypt"
)

// Login renders the HTML of the login page
func Login(c *gin.Context) {
	pd := DefaultPageData(c)
	pd.Title = pd.Trans("Login")
	c.HTML(http.StatusOK, "login.gohtml", pd)
}

// LoginPost handles login requests and returns the appropriate HTML and messages
func LoginPost(c *gin.Context) {
	pd := DefaultPageData(c)
	loginError := pd.Trans("Could not login, please make sure that you have typed in the correct email and password. If you have forgotten your password, please click the forgot password link below.")
	pd.Title = pd.Trans("Login")

	db := infra.LairInstance().GetDb()

	email := c.PostForm("email")
	user := models.User{Email: email}

	res := db.Where(&user).First(&user)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		slog.Error("LoginPost", "error", res.Error)
		c.HTML(http.StatusInternalServerError, "login.gohtml", pd)
		return
	}

	if res.RowsAffected == 0 {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.gohtml", pd)
		return
	}

	if user.ActivatedAt == nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: pd.Trans("Account is not activated yet."),
		})
		c.HTML(http.StatusBadRequest, "login.gohtml", pd)
		return
	}

	password := c.PostForm("password")
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		c.HTML(http.StatusBadRequest, "login.gohtml", pd)
		return
	}

	// Generate a ulid for the current session
	sessionIdentifier := ulid.Generate()

	ses := models.Session{
		Identifier: sessionIdentifier,
	}

	// Session is valid for 1 hour
	ses.ExpiresAt = time.Now().Add(time.Hour)
	ses.UserID = user.ID

	res = db.Save(&ses)
	if res.Error != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		slog.Error("LoginPost", "error", res.Error)
		c.HTML(http.StatusInternalServerError, "login.gohtml", pd)
		return
	}

	session := middleware.DefaultSessionWithOptions(c)
	session.Set(middleware.SessionIDKey, sessionIdentifier)
	// Safari strictness requires the following

	err = session.Save()
	if err != nil {
		pd.Messages = append(pd.Messages, Message{
			Type:    "error",
			Content: loginError,
		})
		slog.Error("LoginPost", "error", err)
		c.HTML(http.StatusInternalServerError, "login.gohtml", pd)
		return
	}

	//c.Redirect(http.StatusTemporaryRedirect, "/admin")
	c.Redirect(http.StatusMovedPermanently, "/")
}
