package postgres

import "fmt"
import "database/sql"
import _ "github.com/lib/pq"

func insertAlbum(db *sql.DB) {
  stmtIns, err := db.Prepare("INSERT INTO album (band_id, name) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  result, err := stmtIns.Exec(1, "Dawn of Ash")
  // func (s *Stmt) Exec(args ...interface{}) (Result, error)
  // one might ask themself what happens when one of the Go-side arguments doesn't match the SQL-side type
  // driverArgs found here https://golang.org/src/database/sql/convert.go
  // is the entry point to this whole process. 
  // btw.
  // result type https://golang.org/pkg/database/sql/driver/#Result
  if err != nil {
    panic(err.Error()) 
  }
  rows_affected, _ := result.RowsAffected()
  // oh yeah btw
  // LastInsertId doesn't work with PSQL. 
  // the function happily returns 0, but the error is set. 
  // just sayin'
  last_id, last_insert_id_error := result.LastInsertId()
  if last_insert_id_error != nil {
    fmt.Println("er: ", last_insert_id_error.Error())
  } else {
    fmt.Println("what is going on with last insert id")
  }
  fmt.Println("num: ", rows_affected, ". id: ", last_id)
}

func insertBand(db *sql.DB) {
  stmtIns, err := db.Prepare("INSERT INTO band (name) VALUES ($1)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  result, err := stmtIns.Exec("Liturgy (the tech death one)")
  // result type https://golang.org/pkg/database/sql/driver/#Result
  if err != nil {
    panic(err.Error()) 
  }
  rows_affected, _ := result.RowsAffected()
  last_id, _ := result.LastInsertId()
  fmt.Println("num: ", rows_affected, ". id: ", last_id)
}

func TestInsert() {
  db, err := sql.Open("postgres", "user=usr password=access dbname=go_test sslmode=disable")

  // db type https://golang.org/pkg/database/sql/driver/#Conn
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  insertBand(db)
  insertAlbum(db)
}
