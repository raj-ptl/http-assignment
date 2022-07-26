package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/raj-ptl/http-assignment/book"
)

func ServeRequests() {
	fmt.Println("Serving now...")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/getBook/", getBookHandler)
	http.HandleFunc("/addBook/", addBookHandler)
	http.ListenAndServe("127.0.0.1:9090", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to server")
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/getBook/")

	if id == "" {
		jsonIdMissing, _ := json.Marshal("book id missing in url path")
		w.Write(jsonIdMissing)
	} else {
		newBook, errRead := book.ReadFromDisk(id)
		if errRead != nil {
			jsonErrRead, _ := json.Marshal("book with specified id does not exist")
			w.Write(jsonErrRead)
		} else {
			bookJson, errMarshal := json.Marshal(newBook)
			if errMarshal != nil {
				jsonErrMarshal, _ := json.Marshal(errMarshal)
				w.Write(jsonErrMarshal)
			} else {
				w.Write(bookJson)
			}
		}
	}
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		jsonInvalidMethod, _ := json.Marshal("POST Method expected on this endpoint")
		w.Write(jsonInvalidMethod)
	} else {
		body, errBodyParse := ioutil.ReadAll(r.Body)

		if errBodyParse != nil {
			jsonErrBodyParse, _ := json.Marshal(errBodyParse)
			w.Write(jsonErrBodyParse)
		} else {
			fmt.Println(string(body))

			var b book.Book

			errUnmarshal := json.Unmarshal(body, &b)

			if errUnmarshal != nil {
				jsonErrUnmarshal, _ := json.Marshal(errUnmarshal)
				w.Write(jsonErrUnmarshal)
			} else {
				b = *book.NewBook(b.Title, b.Description, b.Price)
				saveErr := b.SaveToDisk()

				if saveErr != nil {
					jsonSaveErr, _ := json.Marshal(saveErr)
					w.Write(jsonSaveErr)
				} else {
					saveSuccessMessage, _ := json.Marshal(fmt.Sprintf("Saved the book to disk, id : %s", b.Id))
					w.Write(saveSuccessMessage)
				}
			}
		}
	}
}
