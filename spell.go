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
	"strings"
	"time"

	"github.com/kyokomi/emoji"
)

var emoji_array []string
var spell_type []string
var spell_preposition []string
var spell_subject []string

func main() {
	emoji_array = LoadJson("emoji.json")
	spell_type = LoadJson("type.json")
	spell_preposition = LoadJson("preposition.json")
	spell_subject = LoadJson("subject.json")

	rand.Seed(time.Now().UTC().UnixNano())

	//fmt.Println(GenerateSpell())
	Server()
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
	//return spell.String()
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

func RandomArg(subject []string) string {
	return subject[rand.Intn(len(subject)-1)]
}
