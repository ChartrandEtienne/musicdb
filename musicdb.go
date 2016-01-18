package main

import "musicdb/lib/files"
import "os"
import "fmt"
import "musicdb/lib/database/postgres"
import "musicdb/lib/webapp"

// what the fuck does this program do anyways
// musicdb read /home/usr/music
// musicdb serve (url? port? dunno?)


func read(args []string) {
  var directory = "."
  if 3 == len(args) {
    directory = args[2]
  }
  _, err := os.Open(directory) 
  if err != nil {
    fmt.Println("no such directory: ", directory)
    return
  }
  postgres.InitDatabase()
  fmt.Println("reading directory ", directory)
  files.Walk(directory)
}

func serve() {
  fmt.Println("implement me")
}

func main() {
  args := os.Args

  if 1 == len(args) {
    fmt.Println("Usage:")
    fmt.Println("musicdb read [directory] to find all music in [directory]")
    fmt.Println("with current directory as default")
  } else if "read" == args[1] {
    read(args)
  } else if "serve" == args[1] {
    webapp.Serve()
  }
}
