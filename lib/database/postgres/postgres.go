package postgres

import "database/sql"
import _ "github.com/lib/pq"

// once again I could do something really evil and have an open connection as a module global
// TODO: investigate best practices

var global_db *sql.DB

func init() {
  // TODO: best practices for storing config data
  db, err := sql.Open("postgres", "user=usr password=access dbname=go_test sslmode=disable")
  if err != nil {
    panic(err.Error()) 
  }
  global_db = db
}

func InsertAlbumReturnId(album_name string, band_id int) int {
  stmtIns, err := global_db.Prepare("INSERT INTO album (name, band_id) VALUES ($1, $2)")
  // type stmtIns https://golang.org/pkg/database/sql/driver/#Stmt
  // TODO: get type inference working in vim
  // if possible
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
