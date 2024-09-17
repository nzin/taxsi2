package config

import "time"

// Config is the whole configuration of the app
var Config = struct {
	// Host - golang-skeleton server host
	Host string `env:"HOST" envDefault:"localhost"`
	// Port - golang-skeleton server port
	Port int `env:"PORT" envDefault:"18000"`

	// LogrusLevel sets the logrus logging level
	LogrusLevel string `env:"TAXSI2_LOGRUS_LEVEL" envDefault:"info"`
	// LogrusFormat sets the logrus logging formatter
	// Possible values: text, json
	LogrusFormat string `env:"TAXSI2_LOGRUS_FORMAT" envDefault:"json"`

	// MiddlewareVerboseLoggerEnabled - to enable the negroni-logrus logger for all the endpoints useful for debugging
	MiddlewareVerboseLoggerEnabled bool `env:"TAXSI2_MIDDLEWARE_VERBOSE_LOGGER_ENABLED" envDefault:"true"`
	// MiddlewareVerboseLoggerExcludeURLs - to exclude urls from the verbose logger via comma separated list
	MiddlewareVerboseLoggerExcludeURLs []string `env:"TAXSI2_MIDDLEWARE_VERBOSE_LOGGER_EXCLUDE_URLS" envDefault:"" envSeparator:","`
	// MiddlewareGzipEnabled - to enable gzip middleware
	MiddlewareGzipEnabled bool `env:"TAXSI2_MIDDLEWARE_GZIP_ENABLED" envDefault:"true"`

	/**
	    DBDriver and DBConnectionStr define how we can write and read data.
		For databases, taxsi2 supports sqlite3, mysql and postgres.

		Examples
		GOLANG_SKELETON_DBDRIVER     GOLANG_SKELETON_DBCONNECTIONSTR
		=========================     =======================================
		"sqlite3"                     "/tmp/file.db"
		"sqlite3"                     ":memory"
		"mysql"                       "root:@tcp(127.0.0.1:3306)/golangskeleton?parseTime=true"
		"postgres"                    "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	*/
	DBDriver                  string        `env:"TAXSI2_DBDRIVER" envDefault:"sqlite3"`
	DBConnectionStr           string        `env:"TAXSI2_DBCONNECTIONSTR" envDefault:"golang-skeleton.sqlite3"`
	DBConnectionRetryAttempts uint          `env:"TAXSI2_DBCONNECTION_RETRY_ATTEMPTS" envDefault:"9"`
	DBConnectionRetryDelay    time.Duration `env:"TAXSI2_DBCONNECTION_RETRY_DELAY" envDefault:"100ms"`

	/*
		Output can be empty, or can contains a list of outputs that:
		- start with "blocked:", "all:"
		- continue with "logrus", "stdout", or a filename
		For example
		- "blocked:logrus"
		- "all:/tmp/all.txt"
	*/
	WafOutput string `env:"TAXSI2_WAF_OUTPUT" envDefault:"blocked:logrus"`
	/*
	  we can use the different variables in the waf output format string:
	  - {{.Date}} (human readable)
	  - {{.Timestamp}} (UTC timestamp)
	  - {{.Url}}
	  - {{.UrlHostname}}
	  - {{.UrlPath}}
	  - {{.Method}}
	  - {{.Remoteaddr}}
	  - {{.Scanresult}} (blocked, dryrun, pass)
	*/
	WafOutputFormat string `env:"TAXSI2_WAF_OUTPUT_FORMAT" envDefault:"{{.Remoteaddr}} {{.Method}} {{.UrlHostname}}:{{.UrlPath}} {{.Scanresult}}"`
}{}
