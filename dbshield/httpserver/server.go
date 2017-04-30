package httpserver

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/training"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

var singleHTTP sync.Once
var assets = "assets"

func serve() {
	_, serverGoLocation, _, ok := runtime.Caller(0)
	if ok {
		assets = filepath.Clean(serverGoLocation + "/../../../" + assets)
	}

	if fi, err := os.Stat(assets); os.IsNotExist(err) || !fi.Mode().IsDir() {
		logger.Warning("Could not find 'assets' directory, Web UI disabled.")
		return
	}

	http.Handle("/js/", http.FileServer(http.Dir(assets)))
	http.Handle("/css/", http.FileServer(http.Dir(assets)))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/report.htm", mainHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
}

//Serve HTTP
func Serve() error {
	singleHTTP.Do(serve)
	if config.Config.HTTPSSL {
		return http.ListenAndServeTLS(config.Config.HTTPAddr, config.Config.TLSCertificate, config.Config.TLSPrivateKey, nil)
	}
	return http.ListenAndServe(config.Config.HTTPAddr, nil)
}

//mainHandler for html
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if checkLogin(r) {
		http.ServeFile(w, r, assets+"/report.htm")
		return
	}
	http.ServeFile(w, r, assets+"/index.htm")
}

//apiHandler returns state
func apiHandler(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(r) {
		return
	}
	out, _ := json.Marshal(struct {
		Total    uint64
		Abnormal uint64
	}{
		training.QueryCounter,
		training.AbnormalCounter,
	})
	w.Write(out)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	pass := r.FormValue("password")
	redirectTarget := "/"
	if pass == config.Config.HTTPPassword {
		setSession(w)
		redirectTarget = "/report.htm"
		logger.Warningf("Successful login to web UI [%s] [%s]", r.RemoteAddr, r.UserAgent())
	} else {
		logger.Warningf("Failed login to web UI [%s] [%s]", r.RemoteAddr, r.UserAgent())
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		clearSession(w)
		http.Redirect(w, r, "/", 302)
	}
}

func checkLogin(r *http.Request) bool {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			check := cookieValue["Login"]
			return check != ""
		}
	}
	return false
}

func setSession(w http.ResponseWriter) {
	value := map[string]string{
		"Login": time.Now().String(),
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:     "session",
			Value:    encoded,
			Path:     "/",
			Secure:   config.Config.HTTPSSL,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
