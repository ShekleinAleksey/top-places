package app

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/ShekleinAleksey/top-places/internal/handler"
	"github.com/ShekleinAleksey/top-places/internal/repository"
	"github.com/ShekleinAleksey/top-places/internal/service"
	"github.com/ShekleinAleksey/top-places/pkg/postgres"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

// @title BestPlace Service
// @version 1.0
// @description API Service for BestPlace App
// @host 95.174.91.82:8080
// @BasePath /
func Run() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	logrus.Info("Initializing db...")

	db, err := postgres.NewDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}
	defer db.Close()

	logrus.Info("Initializing repository...")
	repos := repository.NewRepository(db)
	logrus.Info("Initializing service...")
	services := service.NewService(repos)
	logrus.Info("Initializing handler...")
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("best-place.online"),
		Cache:      autocert.DirCache("/var/www/.cache"), // Папка для хранения сертификатов
	}

	// Настройка TLS
	tlsConfig := &tls.Config{
		GetCertificate: certManager.GetCertificate,
		MinVersion:     tls.VersionTLS12, // Современные безопасные настройки
	}

	// Создание HTTP-сервера с поддержкой HTTPS
	srv := &http.Server{
		Addr:      ":8080",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	// Перенаправление HTTP -> HTTPS
	go func() {
		http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	}()

	logrus.Info("Starting HTTPS server...")
	if err := srv.ListenAndServeTLS("", ""); err != nil { // Пустые строки, так как сертификаты управляются autocert
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
