package main
/
import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strconv"
	"sync"

	"persha.maxg95/internal/database"
	"persha.maxg95/internal/smtp"
	"persha.maxg95/internal/version"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	db struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
}

type application struct {
	config config
	db     *database.DB
	logger *slog.Logger
	mailer *smtp.Mailer
	wg     sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = os.Getenv("BASE_URL")
	cfg.httpPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	cfg.cookie.secretKey = os.Getenv("COOKIE_SECRET_KEY")
	cfg.db.dsn = os.Getenv("DB_DSN")
	cfg.smtp.host = os.Getenv("SMTP_HOST")
	cfg.smtp.port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	cfg.smtp.username = os.Getenv("SMTP_USERNAME")
	cfg.smtp.password = os.Getenv("SMTP_PASSWORD")
	cfg.smtp.from = os.Getenv("SMTP_FROM")

	showVersion := os.Getenv("SHOW_VERSION")

	if showVersion != "" {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	if err != nil {
		return err
	}

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		mailer: mailer,
	}

	return app.serveHTTP()
}
