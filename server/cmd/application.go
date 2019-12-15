package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/batect/abacus/server/api"
	"github.com/sirupsen/logrus"
)

func main() {
	initLogging()

	srv := createServer(getPort())
	runServer(srv)
}

func initLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		logrus.Fatal("PORT environment variable is not set.")
	}

	return port
}

func createServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", api.Ping)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	return srv
}

func runServer(srv *http.Server) {
	connectionDrainingFinished := shutdownOnInterrupt(srv)

	logrus.WithField("address", srv.Addr).Info("Server starting.")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithField("err", err).Fatal("Could not start HTTP server.")
	}

	<-connectionDrainingFinished

	logrus.Info("Server shut down.")
}

func shutdownOnInterrupt(srv *http.Server) chan struct{} {
	connectionDrainingFinished := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logrus.Info("Interrupt received, draining connections...")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithField("err", err).Error("Shutting down HTTP server failed.")
		}

		close(connectionDrainingFinished)
	}()

	return connectionDrainingFinished
}
