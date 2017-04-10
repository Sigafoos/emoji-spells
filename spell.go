package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kyokomi/emoji"
)

var emoji_array []string
var spell_type []string
var spell_preposition []string
var spell_subject []string
var options map[string]string

func main() {
	emoji_array = LoadJson("emoji.json")
	spell_type = LoadJson("type.json")
	spell_preposition = LoadJson("preposition.json")
	spell_subject = LoadJson("subject.json")
	options = LoadOptions("options.json")

	rand.Seed(time.Now().UTC().UnixNano())

	spell := GenerateSpell()
	Post(spell)
}

func GenerateSpell() string {
	var spell bytes.Buffer

	spell.WriteString("[")
	spell.WriteString(RandomArg(spell_type))
	spell.WriteString(" ")
	spell.WriteString(RandomArg(spell_preposition))
	spell.WriteString(" ")
	spell.WriteString(RandomArg(spell_subject))

	spell.WriteString("]\n\n")

	for i := 0; i < rand.Intn(3)+3; i++ {
		spell.WriteString(emoji.Sprint(RandomArg(emoji_array)))
	}
	return strings.ToUpper(spell.String())
}

func Server() {
	port := ":8080"
	http.HandleFunc("/", EmojiList)

	fmt.Printf("starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func EmojiList(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, GenerateSpell())
}

func LoadJson(filename string) []string {
	var parsed []string

	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
	}
	e = json.Unmarshal(file, &parsed)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
	}

	return parsed
}

func LoadOptions(filename string) map[string]string {
	var parsed map[string]string
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
	}
	e = json.Unmarshal(file, &parsed)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
	}

	return parsed
}

func RandomArg(subject []string) string {
	return subject[rand.Intn(len(subject)-1)]
}

func Post(spell string) {
	data := url.Values{}
	data.Set("status", spell)
	resource := "/api/v1/statuses"
	u, _ := url.ParseRequestURI(options["instance"])
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)
	req, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %v", err))
	}
	req.Header.Add("Authorization", "Bearer "+options["key"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("ERROR with req: %v\n", err))
	}
	resp.Body.Close()
}
