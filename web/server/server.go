package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cody6750/web-crawler/web/handler"
	"github.com/cody6750/web-crawler/web/options"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// WebCrawlerServer contains all functions and dependencies to create
// a web server that hosts a web crawler.
type WebCrawlerServer struct {
	server  *http.Server
	logger  *logrus.Logger
	Options *options.Options
}

// New creates and returns the web crawler server object with default options
func New() *WebCrawlerServer {
	return NewWithOptions(options.New())
}

// NewWithOptions creates and returns the web crawler server object with custom options
func NewWithOptions(o *options.Options) *WebCrawlerServer {
	return &WebCrawlerServer{Options: o}
}

// init intializes required objects for the web crawler server. Gets all necessary environment variables
// that override options and default variables.
func (wcs *WebCrawlerServer) init() {
	wcs.logger = logrus.New()
	wcs.logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	wcs.logger.Info("Initializing web server")
	err := wcs.processEnvironmentVariables()
	if err != nil {
		wcs.logger.Errorf("Error processing environment variables %e", err)
	}
	wcs.logger.Info("Successfully initialized web server")
}

// Run serves as the entrypoint for intializing the web server. It generates all of the
// handlers that handle http request and attaches those to the web server. Once initialized
// the server begings to listen and serve request. Any interrupt or kill signal from the bot
// host will terminate the web server which will gracefully shutdown.
func (wcs *WebCrawlerServer) Run() {
	wcs.init()
	wcs.logger.Info("Starting up web server")
	serverMux := wcs.generateHandlers()
	wcs.server = &http.Server{
		Addr:         wcs.Options.Port,
		Handler:      serverMux,
		IdleTimeout:  wcs.Options.IdleTimeout,
		ReadTimeout:  wcs.Options.ReadTimeout,
		WriteTimeout: wcs.Options.WriteTimeout,
	}

	// start the server
	go func() {
		err := wcs.server.ListenAndServe()
		if err != nil {
			wcs.logger.WithError(err).Error("Web server shutting down")
		}
	}()
	wcs.logger.Info("Successfully started up web server, listening for traffic")

	// trap sigterm or interupt and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received.
	sig := <-sigChan
	wcs.logger.Infof("Recieved teriminate, graceful shutdown ", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	wcs.server.Shutdown(tc)
}

// generateHandlers generates and ataches the handlers to the web server
func (wcs *WebCrawlerServer) generateHandlers() *mux.Router {
	wcs.logger.Debug("Starting to generate handlers for web server")
	serverMux := mux.NewRouter()
	crawler := handler.NewCrawler(wcs.logger)

	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/crawler/item", crawler.GetItem)
	getRouter.Use(crawler.MiddlewareItemValidation)

	wcs.logger.WithField("Handlers", serverMux).Debug("Successfully generated handlers for web server")
	return serverMux
}
