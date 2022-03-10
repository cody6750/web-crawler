package server

import (
	"os"

	env "github.com/cody6750/web-crawler/shared"
	"github.com/sirupsen/logrus"
)

// processEnvironmentVariables gets and sets all specified environemnt variables for the
// web crawler server
func (wcs *WebCrawlerServer) processEnvironmentVariables() error {
	if os.Getenv("PORT") != "" {
		wcs.Options.Port = os.Getenv("PORT")
		wcs.logger.WithField("PORT", wcs.Options.Port).Debug("PORT overide found. Overriding with value %v", wcs.Options.Port)
	}

	if os.Getenv("LOG_LEVEL") != "" {
		wcs.Options.LogLevel = os.Getenv("LOG_LEVEL")
		wcs.logger.WithField("LOG_LEVEL", wcs.Options.LogLevel).Debug("LOG_LEVEL overide found. Overriding with value %v", wcs.Options.LogLevel)
		switch wcs.Options.LogLevel {
		case "INFO":
			wcs.logger.SetLevel(logrus.InfoLevel)
		case "WARN":
			wcs.logger.SetLevel(logrus.WarnLevel)
		case "DEBUG":
			wcs.logger.SetLevel(logrus.DebugLevel)
		default:
			wcs.logger.Error("unsupported log type %v, using default logger", wcs.Options.LogLevel)
		}
	}
	if os.Getenv("IDLE_TIMEOUT") != "" {
		time, err := env.GetEnvTime(os.Getenv("IDLE_TIMEOUT"))
		wcs.logger.WithField("IDLE_TIMEOUT", time.Seconds()).Debug("IDLE_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		wcs.Options.IdleTimeout = time
	}
	if os.Getenv("READ_TIMEOUT") != "" {
		time, err := env.GetEnvTime(os.Getenv("READ_TIMEOUT"))
		wcs.logger.WithField("READ_TIMEOUT", time.Seconds()).Debug("READ_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		wcs.Options.ReadTimeout = time
	}
	if os.Getenv("WRITE_TIMEOUT") != "" {
		time, err := env.GetEnvTime(os.Getenv("WRITE_TIMEOUT"))
		wcs.logger.WithField("WRITE_TIMEOUT", time.Seconds()).Debug("WRITE_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		wcs.Options.WriteTimeout = time
	}
	return nil
}
