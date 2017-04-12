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

	"github.com/Sigafoos/jsonrand"
	"github.com/kyokomi/emoji"
)

var emoji_array jsonrand.String
var spell_type jsonrand.String
var spell_preposition jsonrand.String
var spell_subject jsonrand.String
var options map[string]string

func main() {
	emoji_array.Load("emoji.json")
	spell_type.Load("type.json")
	spell_preposition.Load("preposition.json")
	spell_subject.Load("subject.json")
	options = LoadOptions("options.json")

	spell := GenerateSpell()
	Post(spell)
}

func GenerateSpell() string {
	var spell bytes.Buffer

	spell.WriteString("[")
	spell.WriteString(spell_type.Element())
	spell.WriteString(" ")
	spell.WriteString(spell_preposition.Element())
	spell.WriteString(" ")
	spell.WriteString(spell_subject.Element())

	spell.WriteString("]\n\n")

	for i := 0; i < rand.Intn(3)+3; i++ {
		spell.WriteString(emoji.Sprint(emoji_array.Element()))
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
	defer resp.Body.Close()
	fmt.Println(spell)
	//io.Copy(os.Stdout, resp.Body)
}
