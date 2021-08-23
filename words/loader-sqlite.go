package words

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteLoader struct {
	dbPath string
}

func NewSqliteLoader(dbPath string) *SqliteLoader {
	fmt.Println("Opening database " + dbPath)
	loader := new(SqliteLoader)
	loader.dbPath = dbPath
	return loader

}

func splitIntoIntArray(sylabesIndexes string) []int {
	var indexes []int = make([]int, 0)
	var stringIndexes []string = strings.Split(sylabesIndexes, ",")
	for i := 0; i < len(stringIndexes); i++ {
		intIndex, err := strconv.Atoi(stringIndexes[i])
		if err != nil {
			panic("Something wrong happened loading sylabes: Non-int index.")
		}
		indexes = append(indexes, intIndex)
	}
	return indexes
}

func getSylabesById(sylabesList map[int]*Sylabe, sylabesIndexes []int) []*Sylabe {
	var sylabes []*Sylabe = make([]*Sylabe, 0)
	for i := 0; i < len(sylabesIndexes); i++ {
		sylabes = append(sylabes, sylabesList[sylabesIndexes[i]])
	}
	return sylabes
}

func loadWords(db *sql.DB, sylabes map[int]*Sylabe) ([]*Word, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT w.word, w.audio_file, group_concat(ws.sylabe_id, ',') FROM words as w JOIN words_sylabes ws ON ws.word_id = w.id group by w.id order by w.id, ws.position")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	var word, sylabesIndexes, audioFile string
	var words []*Word = make([]*Word, 0)
	for rows.Next() {
		err = rows.Scan(&word, &audioFile, &sylabesIndexes)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fmt.Println(word, sylabesIndexes, audioFile)
		sylabesIntIndexes := splitIntoIntArray(sylabesIndexes)
		wordPointer := new(Word)
		wordPointer.Word = strings.ToUpper(word)
		wordPointer.Sylabes = getSylabesById(sylabes, sylabesIntIndexes)
		wordPointer.AudioFile = audioFile
		words = append(words, wordPointer)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return words, nil
}

func loadSylabes(db *sql.DB) (map[int]*Sylabe, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query("SELECT s.id, s.sylabe, s.audio_file FROM sylabes s order by s.id")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	var sylabe, audioFile string
	var id int
	var sylabes map[int]*Sylabe = make(map[int]*Sylabe)
	for rows.Next() {
		err = rows.Scan(&id, &sylabe, &audioFile)
		fmt.Println(id, sylabe, audioFile)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		sylabeData := new(Sylabe)
		sylabeData.Sylabe = strings.ToUpper(sylabe)
		sylabeData.AudioFile = audioFile
		sylabes[id] = sylabeData
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return sylabes, nil
}

func (loader *SqliteLoader) Load() (GameDataAccessor, error) {
	var err error
	var db *sql.DB
	var gameData = new(GameData)
	var sylabes map[int]*Sylabe
	var words []*Word

	db, err = sql.Open("sqlite3", loader.dbPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	sylabes, err = loadSylabes(db)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	gameData.Sylabes = sylabes

	words, err = loadWords(db, sylabes)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	gameData.Words = words

	if len(gameData.Words) == 0 {
		return nil, errors.New("No words loaded")
	}

	return gameData, nil
}
