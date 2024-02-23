package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *application) serveHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      app.routes(),
		ErrorLog:     log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	certFile := "/etc/letsencrypt/live/persha.lutsk.ua/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/persha.lutsk.ua/privkey.pem"

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	}

	srv.TLSConfig = tlsConfig

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	// Redirect HTTP to HTTPS
	go func() {
		httpRedirect := &http.Server{
			Addr:    ":80",
			Handler: http.HandlerFunc(redirectHandler),
		}

		err := httpRedirect.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP to HTTPS redirect server error: %v\n", err)
		}
	}()

	log.Printf("Starting server on %s\n", srv.Addr)

	err := srv.ListenAndServeTLS(certFile, keyFile)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	log.Printf("Stopped server on %s\n", srv.Addr)

	app.wg.Wait()
	return nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://persha.lutsk.ua"+r.URL.String(), http.StatusMovedPermanently)
}
