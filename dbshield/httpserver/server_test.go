package httpserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/config"
	"github.com/nim4/DBShield/dbshield/training"
)

func TestMain(t *testing.T) {
	config.Config.HTTPAddr = ":0"
	config.Config.HTTPPassword = "foo"

	tmpfile, err := ioutil.TempFile("", "testdb")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	path := tmpfile.Name()
	training.DBCon, err = bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	if err := training.DBCon.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("queries"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("abnormal"))
		return err
	}); err != nil {
		panic(err)
	}
}

func TestServe(t *testing.T) {
	err := Serve()
	if err == nil {
		t.Error("Expected error")
	}
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

	r, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}
	w = httptest.NewRecorder()
	setSession(w)
	r.Header.Set("Cookie", w.HeaderMap.Get("Set-Cookie"))
	mainHandler(w, r)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("Got an error ", err)
	}
	if strings.Index(string(body), "\"Logout\"") == -1 {
		t.Error("Expected report page")
	}
}

func TestAPIHandler(t *testing.T) {
	defer recover()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}
	w := httptest.NewRecorder()
	apiHandler(w, r)
	body, err := ioutil.ReadAll(w.Body)
	if len(body) != 0 {
		t.Error("Expected 0 length got", len(body))
	}
	setSession(w)
	r.Header.Set("Cookie", w.HeaderMap.Get("Set-Cookie"))
	apiHandler(w, r)
	body, err = ioutil.ReadAll(w.Body)
	var j struct {
		Total    int
		Abnormal int
	}
	j.Total = 1
	j.Abnormal = 1

	err = json.Unmarshal(body, &j)
	if err != nil {
		t.Error("Got an error ", err)
	}
	if j.Total != 0 || j.Abnormal != 0 {
		t.Error("Expected 0 got", j)
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

	r, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Got an error ", err)
	}
	w := httptest.NewRecorder()
	setSession(w)
	r.Header.Set("Cookie", w.HeaderMap.Get("Set-Cookie"))

	if !checkLogin(r) {
		t.Error("Expected true got false")
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
