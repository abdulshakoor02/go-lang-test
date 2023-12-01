// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"sync"

// 	_ "github.com/go-sql-driver/mysql"
// )

// type Trips struct {
// 	Id        int `json:"id"`
// 	Unit      int `json:"unit"`
// 	TripBegin int `json:"tripBegin"`
// 	TripEnd   int `json:"tripEnd"`
// }

// func worker(id int, jobs <-chan int, data []Trips, db *sql.DB, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for j := range jobs {
// 		fmt.Printf("Worker %d started job\n", id)
// 		// do some work here...

// 		newq := fmt.Sprintf("call write_events2(%v,%v,%v)", data[j].Unit, data[j].TripBegin, data[j].TripEnd)

// 		res, err := db.Exec(newq)
// 		if err != nil {
// 			fmt.Println(err)

// 		}
// 		rowsAffected, _ := res.RowsAffected()
// 		fmt.Println(rowsAffected)

// 		fmt.Printf("Worker %d finished job\n", id)
// 		// results <- id * 2 // send result back to results channel
// 	}
// }

// func main() {

// 	db, err := sql.Open("mysql", "doadmin:ORjrFFFzoMTJR6dZ@tcp(opalserver-do-user-7849719-0.b.db.ondigitalocean.com:25060)/opal_staging")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer db.Close()
// 	q := `select id,unit,tripBegin,tripEnd from opal_trips where callproc!=1 limit 0,100`
// 	rows, err := db.Query(q)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer rows.Close()
// 	var data []Trips

// 	for rows.Next() {
// 		var id, unit, tripBegin, tripEnd int

// 		if err := rows.Scan(&id, &unit, &tripBegin, &tripEnd); err != nil {
// 			fmt.Println(err)
// 		}

// 		data = append(data, Trips{Id: id, Unit: unit, TripBegin: tripBegin, TripEnd: tripEnd})
// 	}

// 	jobs := make(chan int, len(data))
// 	// results := make(chan int, len(data))
// 	var wg sync.WaitGroup

// 	// start workers
// 	numWorkers := 5
// 	for w := 1; w <= numWorkers; w++ {
// 		wg.Add(1)
// 		go worker(w, jobs, data, db, &wg)
// 	}

// 	// send jobs to workers
// 	for i, _ := range data {
// 		jobs <- i
// 	}
// 	close(jobs)

// 	// collect results
// 	// for a := 1; a <= len(data); a++ {
// 	// 	result := <-results
// 	// 	fmt.Println("Result:", result)
// 	// }

// 	// wait for all workers to finish
// 	wg.Wait()
// }
