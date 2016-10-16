package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
	"github.com/nim4/DBShield/dbshield/training"
)

//Serve HTTP on given address
func Serve(addr string) {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/api/", handle)

	logger.Infof("HTTP server on %s", addr)
	http.ListenAndServe(addr, nil)
}

//handle api count
func handle(w http.ResponseWriter, r *http.Request) {
	abnormal := 0
	total := 0

	//Count queries
	if err := training.DBCon.Update(func(tx *bolt.Tx) error {
		var contextArray []sql.QueryContext
		b := tx.Bucket([]byte("queries"))
		if b == nil {
			panic(errors.New("Bucket not found"))
		}
		return b.ForEach(func(k []byte, v []byte) error {
			if err := json.Unmarshal(v, &contextArray); err != nil {
				return err
			}
			total += len(contextArray)
			return nil
		})
	}); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//Count abnormal
	if err := training.DBCon.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("abnormal"))
		if b == nil {
			panic(errors.New("Bucket not found"))
		}
		abnormal = b.Stats().KeyN
		return nil
	}); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	out, _ := json.Marshal(struct {
		Total    int
		Abnormal int
	}{
		total,
		abnormal,
	})
	w.Write(out)
}
