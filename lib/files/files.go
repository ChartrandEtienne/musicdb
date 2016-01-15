package files

import "fmt"
import "io/ioutil"
import "path/filepath"
import "strings"
import "os"
// import "hello/lib/database/postgres"
import id3 "github.com/mikkyang/id3-go"

// Now I'm not too sure how to do that. 
// I want walkFunc to close over the variable containing the data I gather while walking. 
// I could either declare both walkFunc and this variable at the top level
// or declare both from inside the Walk function.
// let's go with top level just to test the init function stuff described here
// http://thenewstack.io/understanding-golang-packages/

var bands map[string]int

func init() {
  bands = make(map[string]int)
}

// this is a bad habit I've already taken
// every one of my modules has that so I don't have to de-import stuff I don't use. 
func NoOp() {
}

func getSongData(path string) (string, string, string) {
  mp3File, err := id3.Open(path)
  if err != nil {
    return "error", "error", "error"
  }
  defer mp3File.Close()
  whataboutnullchars(mp3File, path)
  return mp3File.Artist(), mp3File.Album(), mp3File.Title()
}

func stripTerminalNull(to_strip string) string {
  // does it strip the final space? Where does it come from?
  return strings.TrimRight(to_strip, "\x00 ")
}

func whataboutnullchars(file *id3.File, path string) {
  taggerVersion := file.Tagger.Version()
  trimmed_artist := stripTerminalNull(file.Artist())
  fmt.Println(file.Artist(), "| to |", trimmed_artist, "|file ", path, " version ", taggerVersion)
  // here I just found out that the id3v1 parser doesn't strip the null padding bytes at the end of id3v1 fields. 
  // also maybe I have some songs with 16bit-wide metadata fields?
  // look
  // artist  Space Witch  |file  /home/usr/music/Holy_mountain/Space Witch - Sleep's Holy Mountain- The CVLT Nation Sessions - 07 From Beyond.mp3  version  2.3.0
}

// now that's nice. 
// time to put that stuff in database. 
// the first entry will involve inserting a band, an album and a title. 
// inserting the album requires the band id
// inserting the title requires the album id

func walkFunc(path string, info os.FileInfo, err error) error {
  // artist, album, title := getSongData(path)
  // I disagree with this design choice
  artist, _, _ := getSongData(path)
  band_id, existing := bands[artist]
  if !existing {
    // fmt.Println("|", artist, "| to |", trimmed_artist, "|in file ", path)
    // band_id = postgres.InsertBandReturnId(artist)
    bands[artist] = band_id
  }
  return nil
}

func Walk() {
  // https://golang.org/src/path/filepath/path.go?s=11458:11503#L381
  // jackpot 
  err := filepath.Walk("/home/usr/music", walkFunc)
  if err != nil {
    panic(err.Error()) 
  }
  fmt.Println("and we get: ", bands)
}


func Files() { fmt.Printf("hello, world\n")
  contents, _ := ioutil.ReadDir("/home/usr/music")
  // fmt.Println("wtf", contents)
  for _, element := range contents {
    fmt.Println("er", element.Name())
  }
}

