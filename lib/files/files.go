package files

import "fmt"
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

// "band" and "artist" are used kind of randomly to refer to the same thing. 
// TODO pick one I guess
var bands map[string]int
var albums map[string]int

type SongFileData struct {
  artist string
  album string
  title string
}

func init() {
  bands = make(map[string]int)
  albums = make(map[string]int)
}

func getSongData(path string) (*SongFileData, error) {
  mp3File, err := id3.Open(path)
  if err != nil {
    fmt.Errorf("error with filepath |", path)
    fmt.Println("error with filepath |", path)
    return nil, err
  }
  defer mp3File.Close()
  // TODO: maybe fix the library, I dunno?
  return &SongFileData{artist: stripTerminalNull(mp3File.Artist()), album: stripTerminalNull(mp3File.Album()), title: stripTerminalNull(mp3File.Title())}, nil
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
// the obvious solution is to store that stuff in a hash

func walkFunc(path string, info os.FileInfo, err error) error {
  is_directory := info.IsDir()
  // directories obvioulsy don't contain music
  if is_directory {
    return nil
  }
  song_data, err := getSongData(path)
  if err != nil {
    // just skip that particular file
    // getSongData error logs the error
    // not sure what's the best approach in this situation
    // there's also the annoyment that returning an error stops the walk for that subtree
    // so just skip that particular file
    return nil
  }
  artist_id, artist_existing := bands[song_data.artist]
  if !artist_existing {
    // blindly assume that database stuff works
    // not sure what's the best approach in this situation
    artist_id = postgres.InsertBandReturnId(song_data.artist)
    bands[song_data.artist] = artist_id
  }
  album_id, album_existing := albums[song_data.album]
  if !album_existing {
    album_id = postgres.InsertAlbumReturnId(song_data.album, artist_id)
    albums[song_data.album] = album_id
  }
  // let's pretend one second that individual songs don't have duplicates. 
  postgres.InsertTrack(song_data.title, album_id)
  return nil
}

func Walk() {
  // https://golang.org/src/path/filepath/path.go?s=11458:11503#L381
  // jackpot 
  // but err won't be set. At least not by walkFunc.
  // TODO: see in what kind of conditions would this return err
  // if it can
  filepath.Walk("/home/usr/music", walkFunc)
  fmt.Println("and we get: ")
  fmt.Println("bands: ", bands)
  fmt.Println("albums: ", albums)
}
