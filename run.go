package baseproject

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/uberswe/golang-base-project/infra"
	"github.com/uberswe/golang-base-project/middleware"
	"github.com/uberswe/golang-base-project/routes"
)

// staticFS is an embedded file system
//
//go:embed web/*
var staticFS embed.FS

// Run is the main function that runs the entire package and starts the webserver, this is called by /cmd/base/main.go
func Run() {

	// We load environment variables, these are only read when the application launches
	conf := infra.LoadEnvVariables()

	// this logger can be used to change runtime setting through web
	loggingLevel := new(slog.LevelVar)
	loggingLevel.Set(slog.LevelDebug) // default to debug
	level, err := infra.StringToLevel(conf.LogLevel)
	if err == nil {
		loggingLevel.Set(level)
	}

	// set default logger level - TODO make this configurable
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggingLevel, // Set the default log level to DEBUG
	}))
	// Set the logger as the default
	slog.SetDefault(logger)
	// Translations
	bundle := infra.LoadLanguageBundles()

	// We connect to the database using the configuration generated from the environment variables.
	db, err := infra.ConnectToDatabase(conf)
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(2)
	}

	// set global app-wide cfg parms
	infra.InitLair(db, conf, bundle, loggingLevel)

	// Once a database connection is established we run any needed migrations
	err = infra.MigrateDatabase(db)
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(3)
	}
	// t will hold all our html templates used to render pages
	var t *template.Template

	// We parse and load the html files into our t variable
	t, err = loadTemplates()
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(4)
	}

	// A gin Engine instance with the default configuration
	r := gin.Default()
	// proxies below should be correct for most internal networks
	err = r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.1/24", "10.0.0.0/8"})
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(5)
	}

	// We create a new cookie store with a key used to secure cookies with HMAC
	store := cookie.NewStore([]byte(conf.CookieSecret))

	// We define our session middleware to be used globally on all routes
	r.Use(sessions.Sessions("golang_base_project_session", store))

	// log request/responses to see info about bots hitting us
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ResponseLogger())

	// We pass our template variable t to the gin engine so it can be used to render html pages
	r.SetHTMLTemplate(t)

	// our assets are only located in a section of our file system. so we create a sub file system.
	subFS, err := fs.Sub(staticFS, "dist/assets")
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(6)
	}

	// All static assets are under the /assets path so we make this its own group called assets
	assets := r.Group("/assets")

	// This middleware sets the Cache-Control header and is applied to the assets group only
	assets.Use(middleware.Cache(conf.CacheMaxAge))

	// All requests to /assets will use the sub fil system which contains all our static assets
	assets.StaticFS("/", http.FS(subFS))

	// Session middleware is applied to all groups after this point.
	r.Use(middleware.Session(db))

	// A General middleware is defined to add default headers to improve site security
	r.Use(middleware.General())

	// Any request to / will call controller.Index
	r.GET("/", routes.Index)

	// We want to handle both POST and GET requests on the /search route. We define both but use the same function to handle the requests.
	r.GET("/search", routes.Search)
	r.POST("/search", routes.Search)
	r.Any("/search/:page", routes.Search)
	r.Any("/search/:page/:query", routes.Search)

	// We define our 404 handler for when a page can not be found
	r.NoRoute(routes.NoRoute)

	// noAuth is a group for routes which should only be accessed if the user is not authenticated
	noAuth := r.Group("/")
	noAuth.Use(middleware.NoAuth())

	noAuth.GET("/login", routes.Login)
	noAuth.GET("/register", routes.Register)
	noAuth.GET("/activate/resend", routes.ResendActivation)
	noAuth.GET("/activate/:token", routes.Activate)
	noAuth.GET("/user/password/forgot", routes.ForgotPassword)
	noAuth.GET("/user/password/reset/:token", routes.ResetPassword)
	noAuth.GET("/log/:level", routes.Loglevel)

	noAuth.GET("/config", routes.ConfigRouteHandler)
	// We make a separate group for our post requests on the same endpoints so that we can define our throttling middleware on POST requests only.
	noAuthPost := noAuth.Group("/")
	noAuthPost.Use(middleware.Throttle(conf.RequestsPerMinute))

	noAuthPost.POST("/config", routes.ConfigRouteHandlerPost)
	noAuthPost.POST("/loglevel", routes.LoggingRouteHandlerPost)
	noAuthPost.POST("/login", routes.LoginPost)
	noAuthPost.POST("/register", routes.RegisterPost)
	noAuthPost.POST("/activate/resend", routes.ResendActivationPost)
	noAuthPost.POST("/user/password/forgot", routes.ForgotPasswordPost)
	noAuthPost.POST("/user/password/reset/:token", routes.ResetPasswordPost)

	// the admin group handles routes that should only be accessible to authenticated users
	admin := r.Group("/")
	admin.Use(middleware.Auth())
	admin.Use(middleware.Sensitive())

	admin.GET("/admin", routes.Admin)
	// We need to handle post from the login redirect
	admin.POST("/admin", routes.Admin)
	admin.GET("/logout", routes.Logout)

	// This starts our webserver, our application will not stop running or go past this point unless
	// an error occurs or the web server is stopped for some reason. It is designed to run forever.
	err = r.Run(":" + conf.Port)
	if err != nil {
		slog.Error("Run", "error", err)
		os.Exit(7)
	}
}
