// Package config defines the env configuration variables
package infra

// Config defines all the configuration variables for the golang-base-project
type Config struct {
	LogLevel          string
	Port              string
	CookieSecret      string
	Database          string
	DatabaseHost      string
	DatabasePort      string
	DatabaseName      string
	DatabaseUsername  string
	DatabasePassword  string
	BaseURL           string
	SMTPUsername      string
	SMTPPassword      string
	SMTPHost          string
	SMTPPort          string
	SMTPSender        string
	RequestsPerMinute int
	CacheParameter    string
	CacheMaxAge       int
}
