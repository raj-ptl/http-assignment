package book

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
)

type Book struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
}

func NewBook(title string, description string, price float32) *Book {

	book := Book{
		Id:          uuid.New(),
		Title:       title,
		Description: description,
		Price:       price,
	}

	return &book
}

func (b *Book) SaveToDisk() error {

	bookString, err := json.Marshal(*b)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(b.Id.String()+".json", []byte(fmt.Sprintf("%v", string(bookString))), 0666)
}

func ReadFromDisk(fileName string) (Book, error) {
	jsonByteString, errRead := ioutil.ReadFile(fileName + ".json")

	if errRead != nil {
		return Book{}, errRead
	}

	var bookOnDisk Book

	errUnmarshal := json.Unmarshal(jsonByteString, &bookOnDisk)

	if errUnmarshal != nil {
		return Book{}, errUnmarshal
	}

	return bookOnDisk, nil
}
