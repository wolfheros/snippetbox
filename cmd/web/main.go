package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// global data structure, initial in main() funtion.
type application struct {
	errorLog      	*log.Logger
	infoLog       	*log.Logger
	session       	*sessions.Session
	snippets      	*mysql.SnippetModel
	templateCache map[string]*template.Template
	users			*mysql.UserModel
}

func main() {
	// config info input from command line
	addr := flag.String("addr", ":4000", "HTTP network address")
	// define a new command-line flag for the mysql dsn string.
	dsn := flag.String("dsn", "web:199219333@/snippetbox?parseTime=true", "MySQL data source name")

	// Define a 32bits long secrect, which is used to authenticate and encrypt.
	secret := flag.String("secret", "u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4", "Secret Key")
	// parse command line
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open db add connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// make sure anything happened, the db can be close.
	defer db.Close()

	// initial template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// inital session manager, configure expires time
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // set use TLS secure connection

	// Logger is passing to home handler as dependency by use a extral
	// "application" struct, and handlers are become struct's
	// method to access the Logger comming with "application" struct

	// intial the 'application' data struct, Create application struct instance
	app := &application{
		infoLog:       	infoLog,
		errorLog:      	errorLog,
		session:       	session,
		snippets:      	&mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users: 			&mysql.UserModel{DB: db},
	}

	// initial a tls.Config, holding the default TLS settings the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// initial server, set errorlog let server use customised logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	//err = srv.ListenAndServe()

	// Use TLS version of start server method
	srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
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
