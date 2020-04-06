package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ninomaj/snippetbox/pkg/models/psql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *psql.SnippetModel
	templateCache map[string]*template.Template
	users         *psql.UserModel
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "Postgres data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	db, err := openDB(dsn)
	// db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	fmt.Println("You connected to db.")

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true                      // When using TLS
	session.SameSite = http.SameSiteStrictMode // Default: Lax

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &psql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &psql.UserModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		// Uncomment to support min TLS version (more secure, but limiting older browsers to connect)
		// MinVersion: tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Closing all keep-alive connections after 1 min of inactivity
		IdleTimeout: time.Minute,
		// If request header or body is being read for more than 5 s,
		// close underline connection to mitigrate risk of slow-client attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
