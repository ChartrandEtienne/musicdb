package files

import "fmt"
import "io/ioutil"
import "path/filepath"
import "strings"
import "os"
import "hello/lib/database/postgres"
import id3 "github.com/mikkyang/id3-go"

// Now I'm not too sure how to do that. 
// I want walkFunc to close over the variable containing the data I gather while walking. 
// I could either declare both walkFunc and this variable at the top level
// or declare both from inside the Walk function.
// let's go with top level just to test the init function stuff described here
// http://thenewstack.io/understanding-golang-packages/

// band, artist, act, group, sorry about the messy nomenclature
// TODO I guess
var bands map[string]int
var albums map[string]int

func init() {
  bands = make(map[string]int)
  albums = make(map[string]int)
}

// this is a bad habit I've already taken
// every one of my modules has that so I don't have to de-import stuff I don't use. 
func NoOp() {
}

func getSongData(path string) (string, string, string) {
  mp3File, err := id3.Open(path)
  if err != nil {
    // do something better please
    return "error", "error", "error"
  }
  defer mp3File.Close()
  // TODO: maybe fix the library, I dunno?
  return stripTerminalNull(mp3File.Artist()), stripTerminalNull(mp3File.Album()), stripTerminalNull(mp3File.Title())
}

func stripTerminalNull(to_strip string) string {
  // does it strip the final space? Where does it come from?
  // TODO: read more about Unicode I suppose
  return strings.TrimRight(to_strip, "\x00 ")
}

// the appropriate place for this function is in a unit test
// the appropriate place for this style of exploratory code is in unit tests. 
// TODO: figure out TDD in Go
func whataboutnullchars(file *id3.File, path string) {
  taggerVersion := file.Tagger.Version()
  trimmed_artist := stripTerminalNull(file.Artist())
  fmt.Println(file.Artist(), "| to |", trimmed_artist, "|file ", path, " version ", taggerVersion)
  // here I just found out that the id3v1 parser doesn't strip the null padding bytes at the end of id3v1 fields. 
  // also maybe I have some songs with 16bit-wide metadata fields?
  // look
  // er
  // turns out a nul byte, even in a comment, is illegal. 
  // see file notes for originally included offending line. 
  // strings are complicated
}

// now that's nice. 
// time to put that stuff in database. 
// the first entry will involve inserting a band, an album and a title. 
// inserting the album requires the band id
// inserting the title requires the album id

func walkFunc(path string, info os.FileInfo, err error) error {
  // artist, album, title := getSongData(path)
  // declared but unused variables are illegal.
  // I disagree with this design choice.
  artist_name, album_name, track_name := getSongData(path)
  band_id, band_existing := bands[artist_name]
  if !band_existing {
    band_id = postgres.InsertBandReturnId(artist_name)
    bands[artist_name] = band_id
  }
  album_id, album_existing := albums[album_name]
  if !album_existing {
    album_id = postgres.InsertAlbumReturnId(album_name, band_id)
    albums[album_name] = album_id
  }
  // let's pretend one second that individual songs don't have duplicates. 
  postgres.InsertTrack(track_name, album_id)
  return nil
}

func Walk() {
  // https://golang.org/src/path/filepath/path.go?s=11458:11503#L381
  // jackpot 
  err := filepath.Walk("/home/usr/music", walkFunc)
  if err != nil {
    panic(err.Error()) 
  }
  fmt.Println("and we get: ")
  fmt.Println("bands: ", bands)
  fmt.Println("albums: ", albums)
}


func Files() { fmt.Printf("hello, world\n")
  contents, _ := ioutil.ReadDir("/home/usr/music")
  // fmt.Println("wtf", contents)
  for _, element := range contents {
    fmt.Println("er", element.Name())
  }
}

