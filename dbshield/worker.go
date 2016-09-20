package dbshield

import "./utils"

func worker(tasks <-chan utils.DBMS, results chan<- error) {
	for dbms := range tasks {
		err := dbms.Handler()
		dbms.Close()
		results <- err
	}
}
