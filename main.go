
package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	version     = "v0.1.1"
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>/?"
)

type PageData struct {
	Password string
	Error    string
	Length   int
	Version  string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleHome).Methods("GET")
	r.HandleFunc("/generate", handleGenerate).Methods("POST")

	fmt.Println("Server starting at http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", r))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	tmpl.Execute(w, PageData{Length: 16, Version: version})
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	length, _ := strconv.Atoi(r.FormValue("length"))
	if length < 1 {
		length = 16
	}

	charSet := ""
	if r.FormValue("lower") == "on" {
		charSet += lowerChars
	}
	if r.FormValue("upper") == "on" {
		charSet += upperChars
	}
	if r.FormValue("digits") == "on" {
		charSet += digitChars
	}
	if r.FormValue("symbols") == "on" {
		charSet += symbolChars
	}

	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	data := PageData{Length: length, Version: version}

	if charSet == "" {
		data.Error = "Please select at least one character set"
	} else {
		data.Password = generatePassword(length, charSet)
	}

	tmpl.Execute(w, data)
}

func generatePassword(length int, charSet string) string {
	var password strings.Builder
	charSetLen := big.NewInt(int64(len(charSet)))

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, charSetLen)
		if err != nil {
			return ""
		}
		password.WriteByte(charSet[index.Int64()])
	}

	return password.String()
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Password Generator</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 2rem auto;
            padding: 0 1rem;
            background: #f0f2f5;
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .options {
            margin: 1rem 0;
        }
        .option {
            margin: 0.5rem 0;
        }
        .length-input {
            margin: 1rem 0;
        }
        .length-input input {
            width: 60px;
            padding: 4px;
        }
        button {
            background: #4CAF50;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1rem;
        }
        button:hover {
            background: #45a049;
        }
        .result {
            margin-top: 1rem;
            padding: 1rem;
            background: #f8f9fa;
            border-radius: 4px;
            word-break: break-all;
        }
        .error {
            color: #dc3545;
            margin-top: 1rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Password Generator <small style="font-size: 0.5em; color: #666;">{{.Version}}</small></h1>
        <form method="POST" action="/generate">
            <div class="length-input">
                <label for="length">Password Length:</label>
                <input type="number" id="length" name="length" value="{{.Length}}" min="1" max="100">
            </div>
            <div class="options">
                <div class="option">
                    <input type="checkbox" id="lower" name="lower" checked>
                    <label for="lower">Lowercase letters (a-z)</label>
                </div>
                <div class="option">
                    <input type="checkbox" id="upper" name="upper" checked>
                    <label for="upper">Uppercase letters (A-Z)</label>
                </div>
                <div class="option">
                    <input type="checkbox" id="digits" name="digits" checked>
                    <label for="digits">Digits (0-9)</label>
                </div>
                <div class="option">
                    <input type="checkbox" id="symbols" name="symbols" checked>
                    <label for="symbols">Symbols (!@#$%^&*...)</label>
                </div>
            </div>
            <button type="submit">Generate Password</button>
        </form>
        {{if .Password}}
        <div class="result">
            <strong>Generated Password:</strong> {{.Password}}
        </div>
        {{end}}
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
    </div>
</body>
</html>
`
