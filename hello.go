package main

import "fmt"
import "hello/lib/files"
import "hello/lib/database"


func main() {
  fmt.Println("right.")
  files.NoOp()
  // files.Databases()
  database.Databases()
}
