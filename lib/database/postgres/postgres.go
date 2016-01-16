package postgres

import "fmt"
import "database/sql"
import _ "github.com/lib/pq"

// once again I could do something really evil and have an open connection as a module global

var global_db *sql.DB

func init() {
  db, err := sql.Open("postgres", "user=usr password=access dbname=go_test sslmode=disable")
  if err != nil {
    panic(err.Error()) 
  }
  global_db = db
}


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

func InsertAlbumReturnId(album_name string, band_id int) int {
  stmtIns, err := global_db.Prepare("INSERT INTO album (name, band_id) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(album_name, band_id)
  if err != nil {
    panic(err.Error()) 
  }
  var album_id int
  err = global_db.QueryRow("SELECT last_value FROM album_id_seq").Scan(&album_id)
  if err != nil {
    panic(err.Error()) 
  }
  return album_id
}

func InsertBandReturnId(band string) int {
  stmtIns, err := global_db.Prepare("INSERT INTO band (name) VALUES ($1)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(band)
  if err != nil {
    panic(err.Error()) 
  }
  var band_id int
  err = global_db.QueryRow("SELECT last_value FROM band_id_seq").Scan(&band_id)
  if err != nil {
    panic(err.Error()) 
  }
  return band_id
}

// all these inserters should do some error management but eh
func InsertTrack(track_name string, album_id int) {
  stmtIns, err := global_db.Prepare("INSERT INTO track (name, album_id) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  if err != nil {
    panic(err.Error()) 
  }
  defer stmtIns.Close() 
  _, err = stmtIns.Exec(track_name, album_id)
  if err != nil {
    panic(err.Error()) 
  }
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
}
