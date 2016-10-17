package httpserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/nim4/DBShield/dbshield/config"
)

func TestMain(t *testing.T) {
	config.Config.HTTPPassword = "foo"
}

func TestMainHandler(t *testing.T) {
	os.Chdir("../../")
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}

	w := httptest.NewRecorder()
	mainHandler(w, r)

	if w.Code != 200 {
		t.Error("Expected 200 got ", w.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	data := url.Values{}
	data.Add("password", "bar")
	r, err := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error("Got an error ", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	loginHandler(w, r)

	if w.Code != 302 {
		t.Error("Expected 302 got ", w.Code)
	}

	if w.HeaderMap.Get("Location") != "/" {
		t.Error("Expected / got ", w.HeaderMap.Get("Location"))
	}

	data.Set("password", config.Config.HTTPPassword)
	r, err = http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error("Got an error ", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	loginHandler(w, r)

	if w.Code != 302 {
		t.Error("Expected 302 got ", w.Code)
	}

	if w.HeaderMap.Get("Location") != "/report.htm" {
		t.Error("Expected /report.htm got ", w.HeaderMap.Get("Location"))
	}
}

func TestLogoutHandler(t *testing.T) {
	r, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}

	w := httptest.NewRecorder()
	logoutHandler(w, r)

	if w.Code != 302 {
		t.Error("Expected 302 got ", w.Code)
	}
}

func TestCheckLogin(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}

	if checkLogin(r) {
		t.Error("Expected false got true")
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    "XYZ",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	r.AddCookie(cookie)
	if checkLogin(r) {
		t.Error("Expected false got true")
	}
}

func TestSetSession(t *testing.T) {
	w := httptest.NewRecorder()
	setSession(w)
	if strings.Index(w.HeaderMap.Get("Set-Cookie"), "session=") == -1 {
		t.Error("Expected session cookie")
	}
}

func TestClearSession(t *testing.T) {
	w := httptest.NewRecorder()
	clearSession(w)
	if strings.Index(w.HeaderMap.Get("Set-Cookie"), "session=;") == -1 {
		t.Error("Expected empty session cookie")
	}
}
