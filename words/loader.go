package words

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Word struct {
	Word      string
	Sylabes   []*Sylabe
	AudioFile string
}

type Sylabe struct {
	Sylabe    string
	AudioFile string
}

type Loader interface {
	Load() (*GameData, error)
}

type GameData struct {
	Words   []*Word
	Sylabes map[int]*Sylabe
}

type GameDataAccessor interface {
	GetRandomWord(exclude []*Word) (*Word, error)
	GetUnusedSylabes(amount int, word *Word) ([]*Sylabe, error)
}

func isExcluded(tmpWord *Word, exclude []*Word) bool {
	is := false
	for i := 0; i < len(exclude) && !is; i++ {
		if exclude[i] == tmpWord {
			is = true
		}
	}
	return is
}

func (data *GameData) GetRandomWord(exclude []*Word) (*Word, error) {
	rand.Seed(time.Now().UnixNano())
	var word *Word = nil
	length := len(data.Words)
	// Atenti: Esto puede ser un loop infinito si todas las palabras están excluídas.
	for word == nil {
		index := rand.Int() % length
		tmpWord := data.Words[index]
		is := false
		for i := 0; i < len(exclude) && !is; i++ {
			if exclude[i] == tmpWord {
				is = true
			}
		}
		if !is {
			word = tmpWord
		}
	}
	return word, nil
}

// func (data *GameData) GetRandomWord(exclude []*Word) (*Word, error) {
// 	rand.Seed(time.Now().UnixNano())
// 	var word *Word = nil
// 	length := len(data.Words)
// 	// Atenti: Esto puede ser un loop infinito si todas las palabras están excluídas.
// 	for word == nil {
// 		index := rand.Int() % length
// 		tmpWord := data.Words[index]
// 		is := false
// 		for i := 0; i < len(exclude) && !is; i++ {
// 			if exclude[i] == tmpWord {
// 				is = true
// 			}
// 		}
// 		if !is {
// 			word = tmpWord
// 		}
// 	}
// 	return word, nil
// }

func (data *GameData) GetUnusedSylabes(amount int, word *Word) ([]*Sylabe, error) {
	// Atenti: Verificar que hay suficientes sílabas
	var sylabes []*Sylabe = make([]*Sylabe, amount)
	count := 0
	length := len(data.Sylabes)
	if len(word.Sylabes)+amount > length {
		return nil, errors.New(fmt.Sprintf("Not enough sylabes to display. Needed %d, have %d", len(word.Sylabes)+amount, length))
	}

	for k, v := range data.Sylabes {
		fmt.Println(fmt.Sprintf("Sylabe k: %d = %s", k, v.Sylabe))
	}

	sectionSize := length / amount
	fmt.Println(fmt.Sprintf("Section size: %d", sectionSize))

	rand.Seed(time.Now().UnixNano())

	for count < amount {
		offset := sectionSize*count + 1
		fmt.Println(fmt.Sprintf("Offset: %d", offset))
		rnd := rand.Int() % sectionSize
		fmt.Println(fmt.Sprintf("Random number: %d", rnd))
		index := offset + rnd
		fmt.Println(fmt.Sprintf("Index: %d", index))
		tmpSylabe := data.Sylabes[index]
		is := false
		for i := 0; i < len(word.Sylabes); i++ {
			if word.Sylabes[i] == tmpSylabe {
				is = true
			}
		}
		if !is {
			fmt.Println(fmt.Sprintf("Adding sylabe: %v", tmpSylabe))
			sylabes[count] = tmpSylabe
			count = count + 1
		}
	}
	return sylabes, nil
}
