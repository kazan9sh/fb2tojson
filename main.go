package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// Структуры для парсинга FB2
type FictionBook struct {
	XMLName xml.Name `xml:"FictionBook"`
	Author  Author   `xml:"description>title-info>author"`
	Sequence Sequence  `xml:"description>title-info>sequence"`
	Genres []string `xml:"description>title-info>genre"`
}

type Author struct {
	FirstName  string `xml:"first-name"`
	LastName   string `xml:"last-name"`
}

type Sequence struct {
	Name   string `xml:"name,attr"`
	Number string `xml:"number,attr"`
}


func main() {
	// Чтение файла fb2
	filePath := "book.fb2"
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer xmlFile.Close()

	// Чтение содержимого файла
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Парсинг XML
	var book FictionBook
	err = xml.Unmarshal(byteValue, &book)
	if err != nil {
		fmt.Println("Ошибка парсинга XML:", err)
		return
	}

	// Формирование структуры автора для JSON
	result := map[string]interface{}{
		"author": map[string]string{
			"first_name":  book.Author.FirstName,
			"last_name":   book.Author.LastName,
		},
		"sequence": map[string]string{
			"name":   book.Sequence.Name,
			"number": book.Sequence.Number,
		},
		"genre": book.Genres,
	}

	// Преобразование в JS
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Ошибка создания JSON:", err)
		return
	}

	// Вывод JSON
	fmt.Println(string(resultJSON))
}
