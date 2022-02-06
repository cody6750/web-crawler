package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cody6750/web-crawler/web/handler"
	"github.com/cody6750/web-crawler/web/tools"
	"github.com/sirupsen/logrus"
)

var (
	logger              = logrus.New()
	defaultAddress      = ":9090"
	defaultIdleTimeout  = time.Second * 120
	defaultReadTimeout  = time.Second * 60
	defaultWriteTimeout = time.Second * 60
	defaultLogLevel     = logrus.InfoLevel
	address             = defaultAddress
	logLevel            = defaultLogLevel
	idleTimeout         = defaultIdleTimeout
	readTimeout         = defaultReadTimeout
	writeTimeout        = defaultWriteTimeout
)

func init() {
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	logger.Info("Initializing web server")
	err := processEnvironmentVariables()
	if err != nil {
		logger.Errorf("Error processing environment variables %e", err)
	}
	logger.Info("Successfully initialized web server")
}

func processEnvironmentVariables() error {
	if os.Getenv("ADDRESS") != "" {
		address = os.Getenv("ADDRESS")
		logger.WithField("ADDRESS", address).Debug("ADDRESS overide found. Overriding with value %v", address)
	}
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel := os.Getenv("LOG_LEVEL")
		logger.WithField("LOG_LEVEL", logLevel).Debug("LOG_LEVEL overide found. Overriding with value %v", logLevel)
		switch logLevel {
		case "INFO":
			logger.SetLevel(logrus.InfoLevel)
		case "WARN":
			logger.SetLevel(logrus.WarnLevel)
		case "DEBUG":
			logger.SetLevel(logrus.DebugLevel)
		default:
			logger.Error("unsupported log type %v, using default logger", logLevel)
		}
	}
	if os.Getenv("IDLE_TIMEOUT") != "" {
		time, err := tools.ConvertStringToTime(os.Getenv("IDLE_TIMEOUT"))
		logger.WithField("IDLE_TIMEOUT", time.Seconds()).Debug("IDLE_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		idleTimeout = time
	}
	if os.Getenv("READ_TIMEOUT") != "" {
		time, err := tools.ConvertStringToTime(os.Getenv("READ_TIMEOUT"))
		logger.WithField("READ_TIMEOUT", time.Seconds()).Debug("READ_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		readTimeout = time
	}
	if os.Getenv("WRITE_TIMEOUT") != "" {
		time, err := tools.ConvertStringToTime(os.Getenv("WRITE_TIMEOUT"))
		logger.WithField("WRITE_TIMEOUT", time.Seconds()).Debug("WRITE_TIMEOUT overide found. Overriding with value %v", time.Seconds())
		if err != nil {
			return err
		}
		writeTimeout = time
	}
	return nil
}

func Run() {
	logger.Info("Starting up web server")

	// handlers for API
	serverMux := generateHandlers(logger)

	// create a new server
	server := http.Server{
		Addr:         address,
		Handler:      serverMux,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.WithError(err).Error("Web server shutting down")
		}
	}()
	logger.Info("Successfully started up web server, listening for traffic")

	// trap sigterm or interupt and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received.
	sig := <-sigChan
	logger.Infof("Recieved teriminate, graceful shutdown ", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}

func generateHandlers(logger *logrus.Logger) *http.ServeMux {
	logger.Debug("Starting to generate handlers for web server")
	serverMux := http.NewServeMux()
	crawler := handler.NewCrawler(logger)
	serverMux.Handle("/crawler/item", crawler)
	logger.WithField("Handlers", serverMux).Debug("Successfully generated handlers for web server")
	return serverMux
}
