package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Trips struct {
	Id        int `json:"id"`
	Unit      int `json:"unit"`
	TripBegin int `json:"tripBegin"`
	TripEnd   int `json:"tripEnd"`
}

// var wg sync.WaitGroup
func main() {
	var data Trips
	jobs := make(chan int, 5)
	db, err := sql.Open("mysql", "doadmin:ORjrFFFzoMTJR6dZ@tcp(opalserver-do-user-7849719-0.b.db.ondigitalocean.com:25060)/opal_staging")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Db Connected!")
	q := `select id,unit,tripBegin,tripEnd from opal_trips where callproc!=1`
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var id, unit, tripBegin, tripEnd int

		if err := rows.Scan(&id, &unit, &tripBegin, &tripEnd); err != nil {
			fmt.Println(err)
		}
		thread := count%5 + 1
		jobs <- 1
		data = Trips{Id: id, Unit: unit, TripBegin: tripBegin, TripEnd: tripEnd}
		go worker(thread, data, jobs, db)
		count++
	}

}

func worker(thread int, data Trips, jobs <-chan int, db *sql.DB) {
	fmt.Printf("Thread %v Writing events for %v\n", thread, data.Unit)
	newq := fmt.Sprintf("call write_events2(%v,%v,%v)", data.Unit, data.TripBegin, data.TripEnd)
	_, err := db.Exec(newq)
	if err != nil {
		fmt.Println(err)
	}
	update := fmt.Sprintf("update opal_trips set callproc=1 where id=%v", data.Id)
	res, err := db.Exec(update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Thread %v finished job for unit %v\n", thread, data.Unit, res)
	<-jobs
	return
}
