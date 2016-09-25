package dbshield

import "github.com/nim4/DBShield/dbshield/utils"

func worker(tasks <-chan utils.DBMS, results chan<- error) {
	for dbms := range tasks {
		err := dbms.Handler()
		dbms.Close()
		results <- err
	}
}
