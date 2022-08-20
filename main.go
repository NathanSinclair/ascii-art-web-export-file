package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type error struct {
	errormsg   string
	errorFound bool
}

var currentError error

// The custom function genAscii generates our map
// returns them as a string
func genAscii(banner string, input string) string {
	var inputString []string
	var tempOut string
	asciiMap := make(map[rune][]string)
	file, err := os.Open(banner)
	if err != nil {
		currentError.errorFound = true
		currentError.errormsg = fmt.Sprintf("404 %s\n", http.StatusText(404))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	char := 31
	for scanner.Scan() {
		if scanner.Text() == "" {
			char++
		} else {
			asciiMap[rune(char)] = append(asciiMap[rune(char)], scanner.Text())
		}
	}
	if len(input) >= 1 {
		if strings.Contains(input, "\\n") {
			inputString = strings.Split(input, "\\n")
		} else {
			inputString = strings.Split(input, "\r") //"\r" means 'carriage return', moves the cursor to the beginning of the line.
		}
	}

	for _, str := range inputString {
		for i := 0; i < 8; i++ {
			for _, srune := range str {
				if srune != rune(10) && srune != rune(13) { // if rune is not 'line feed' or 'carriage return'
					tempOut += asciiMap[srune][i]
				}
			}
			tempOut += "\n"
		}
	}
	//creating two output files of format doc and txt
	err = os.WriteFile("download.doc", []byte(tempOut), 0666)
	if err != nil {
		panic(err)
	}
	err1 := os.WriteFile("download.txt", []byte(tempOut), 0666)
	if err1 != nil {
		panic(err1)
	}
	return tempOut
}

// The custom function Ascii routes our template file to our server.
// Errors are returned if the correct path to the file is not found.
func Ascii(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		t, err := template.ParseFiles("tmpl/index.html")
		if err != nil {
			currentError.errorFound = true
			currentError.errormsg = fmt.Sprintf("404 %s\n", http.StatusText(404))
			log.Fatal(err)
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		banner := r.Form["banner"]
		input := r.Form["asciiInput"]
		banner = append(banner, "")
		input = append(input, "")
		if banner[0] != "standard.txt" && banner[0] != "shadow.txt" && banner[0] != "thinkertoy.txt" || len(input[0]) < 1 {
			if currentError.errorFound {
				t, _ := template.ParseFiles("tmpl/index.html")
				t.Execute(w, currentError.errormsg)
			} else {
				t, _ := template.ParseFiles("tmpl/index.html")
				t.Execute(w, fmt.Sprintf("400 %s\n", http.StatusText(400)))
			}
		} else if currentError.errorFound {
			t, _ := template.ParseFiles("tmpl/index.html")
			t.Execute(w, currentError.errormsg)
		} else {
			asciiOutput := genAscii(banner[0], input[0])
			if asciiOutput == "" {
				t, _ := template.ParseFiles("tmpl/index.html")
				t.Execute(w, fmt.Sprintf("500 %s\n", http.StatusText(500)))
				fmt.Printf("500 %s\n", http.StatusText(500))
			} else {
				t, _ := template.ParseFiles("tmpl/index.html")
				t.Execute(w, asciiOutput)
				fmt.Printf("200 %s\n", http.StatusText(200)) // status 200 OK means request has succeeded
			}
		}
	}
}

func download(w http.ResponseWriter, r *http.Request) {

	formatType := r.FormValue("fileformat")

	f, _ := os.Open("download." + formatType)
	defer f.Close()

	file, _ := f.Stat()
	fsize := file.Size()

	sfSize := strconv.Itoa(int(fsize))

	w.Header().Set("Content-Disposition", "attachment; filename=asciiresults."+formatType)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", sfSize)

	io.Copy(w, f)
}

func main() {
	// setting router rule
	http.HandleFunc("/", Ascii)
	http.HandleFunc("/index.html", Ascii)
	http.HandleFunc("/right", download)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// setting listening port and logging errors to not finding the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("500 %s : %s\n", http.StatusText(500), err.Error()))
	}
}
