package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/navruz-rakhimov/snippetbox/pkg/models"
	"github.com/navruz-rakhimov/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string
const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets interface{
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	session *sessions.Session
	users interface{
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
}

func main() {
	// getting flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// setting custom app Logger's
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// getting a pointer to a pool of connections
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true  // setting the Secure flag on our session cookies

	// initializing our application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session: session,
		users: &mysql.UserModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// initializing our custom server
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog   ,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
		// Idle, Read and Write timeouts to the server for preventing slow-client attacks
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
