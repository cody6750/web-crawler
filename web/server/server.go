package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cody6750/web-crawler/web/handler"
	"github.com/cody6750/web-crawler/web/tools"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	logger              = logrus.New()
	defaultPort         = ":9090"
	defaultIdleTimeout  = time.Second * 120
	defaultReadTimeout  = time.Second * 60
	defaultWriteTimeout = time.Second * 60
	defaultLogLevel     = logrus.InfoLevel
	port                = defaultPort
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
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
		logger.WithField("PORT", port).Debug("PORT overide found. Overriding with value %v", port)
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
		Addr:         port,
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

func generateHandlers(logger *logrus.Logger) *mux.Router {
	logger.Debug("Starting to generate handlers for web server")
	serverMux := mux.NewRouter()
	crawler := handler.NewCrawler(logger)

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/crawler/item", crawler.GetItem)
	getRouter.Use(crawler.MiddlewareItemValidation)

	logger.WithField("Handlers", serverMux).Debug("Successfully generated handlers for web server")
	return serverMux
}
