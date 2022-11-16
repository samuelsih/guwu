package guwu

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/mail"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/routes"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	email := os.Getenv("MAIL_USER")
	if email == "" {
		panic("email is required")
	}

	password := os.Getenv("MAIL_PASSWORD")
	if password == "" {
		panic("password is empty")
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		smtpServer = "smtp.gmail.com"
	}

	var mailPort int

	smtpPort := os.Getenv("MAIL_PORT")
	if smtpPort == "" {
		mailPort = 587
	} else {
		mailPort, _ = strconv.Atoi(smtpPort)
	}

	redis := config.NewRedis("")
	session := model.SessionDeps{Conn: redis}
	routes.InitDependency(session, "", "")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	
	authDeps := routes.AuthDeps{}

	r.Mount("/api", routes.AuthRoutes(authDeps))

	server := &http.Server{
		Addr: ":" + port,
		Handler: r,
	}

	mailer := mail.NewMailer(smtpServer, mailPort, email, password)

	go mailer.Listen()

	log.Info().Msg("Listening on port :8080")
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal().Msg("Cant start server: " + err.Error())
		}
	}
}
