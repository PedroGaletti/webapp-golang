package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Todo struct {
	Count int
	ToDos []string
}

func ErrorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Values(filename string) []string {
	var lines []string

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}

	ErrorCheck(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	ErrorCheck(scanner.Err())

	return lines
}

func Write(writer http.ResponseWriter, message string) {
	_, err := writer.Write([]byte(message))
	ErrorCheck(err)
}

func EnglishHandler(writer http.ResponseWriter, request *http.Request) {
	Write(writer, "Hello internet")
}

func PortgueseHandler(writer http.ResponseWriter, request *http.Request) {
	Write(writer, "Olá internet")
}

func InteractHandler(writer http.ResponseWriter, request *http.Request) {
	vals := Values("todos.txt")
	fmt.Printf("%#v\n", vals)

	template, err := template.ParseFiles("view.html")
	ErrorCheck(err)

	todos := Todo{
		Count: len(vals),
		ToDos: vals,
	}

	template.Execute(writer, todos)
}

func NewHandler(writer http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("new.html")
	ErrorCheck(err)
	template.Execute(writer, nil)
}

func CreateHandler(writer http.ResponseWriter, request *http.Request) {
	todo := request.FormValue("todo")

	file, err := os.OpenFile("todos.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0600))
	ErrorCheck(err)

	_, err = fmt.Fprintf(file, "\n"+todo)
	ErrorCheck(err)

	err = file.Close()
	ErrorCheck(err)

	http.Redirect(writer, request, "/interact", http.StatusFound)
}

func main() {
	http.HandleFunc("/hello", EnglishHandler)
	http.HandleFunc("/ola", PortgueseHandler)
	http.HandleFunc("/interact", InteractHandler)
	http.HandleFunc("/new", NewHandler)
	http.HandleFunc("/create", CreateHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
