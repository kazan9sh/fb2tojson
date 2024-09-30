package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// FictionBook structure for parsing simplified FB2
type FictionBook struct {
	XMLName  xml.Name `xml:"FictionBook"`
	Author   Author   `xml:"description>title-info>author"`
	Sequence Sequence `xml:"description>title-info>sequence"`
	Genres   []string `xml:"genres>genre"` // Changed to reflect new structure
	Body     Body     `xml:"body>section"`
}

type Author struct {
	FirstName string `xml:"first-name"`
	LastName  string `xml:"last-name"`
}

type Sequence struct {
	Name   string `xml:"name,attr"`
	Number string `xml:"number,attr"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"` // Captures all paragraphs
}

type Paragraph struct {
	Text string `xml:",chardata"` // To store paragraph text
}

type Word struct {
	Text string `json:"w"` // Structure for word representation
}

func main() {
	// Reading the fb2 file
	filePath := "book.fb2"
	xmlFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Ошибка открытия файла:", err)
	}
	defer xmlFile.Close()

	// Reading file contents
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal("Ошибка чтения файла:", err)
	}

	// XML parsing
	var book FictionBook
	err = xml.Unmarshal(byteValue, &book)
	if err != nil {
		log.Fatal("Ошибка парсинга XML:", err)
	}

	// Check for paragraphs

	// Prepare structure for JSON
	result := map[string]interface{}{
		"author": map[string]string{
			"first_name": book.Author.FirstName,
			"last_name":  book.Author.LastName,
		},
		"sequence": map[string]string{
			"name":   book.Sequence.Name,
			"number": book.Sequence.Number,
		},
		"genres": book.Genres,
		"words":  extractWords(book.Body.Paragraphs), // Extract words from the book text
	}

	// Convert to JSON
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Ошибка создания JSON:", err)
	}

	// Write JSON to a new file, overwriting each time
	jsonFilePath := "book.json"
	err = ioutil.WriteFile(jsonFilePath, resultJSON, 0644)
	if err != nil {
		log.Fatal("Ошибка записи в файл:", err)
	}

	fmt.Println("JSON сохранен в", jsonFilePath)
}

func extractWords(paragraphs []Paragraph) []Word {
	var words []Word
	for _, para := range paragraphs {
		for _, word := range strings.Fields(para.Text) {
			words = append(words, Word{Text: word}) 
		}
	}
	return words
}
