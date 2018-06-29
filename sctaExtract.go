package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type sctaCollection struct {
	Collection string `json:"collection"`
	URL        string `json:"url"`
}

type sctaItems struct {
	Item string `json:"item"`
	URL  string `json:"url"`
}

type work struct {
	ID  string
	CSV string
}

type passage struct {
	ID    string
	Text  string
	Label string
}

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func main() {
	outputFile := "scta.csv"
	link := "http://scta.info/csv/scta"
	switch len(os.Args) {
	case 1:
		fmt.Println("No parameters given. Using default:")
		fmt.Println("Endpoint:", link)
		fmt.Println("Output-File:", outputFile)
	case 2:
		fmt.Println("No parameters given. Using default:")
		fmt.Println("Endpoint:", link)
	case 3:
		outputFile = os.Args[1]
		link = os.Args[2]
		fmt.Println("Endpoint:", link)
		fmt.Println("Output-File:", outputFile)
	default:
		fmt.Println("Usage: SCTAExtract [output-filename] [optionally: endpoint]")
		os.Exit(3)
	}

	data, _ := getContent(link)
	var inventory []sctaCollection
	json.Unmarshal(data, &inventory)
	var allItems []sctaItems
	for i := range inventory {
		var items []sctaItems
		link = inventory[i].URL
		data, _ := getContent(link)
		json.Unmarshal(data, &items)
		percent := 100 * (i + 1) / (len(inventory) - 1)
		allItems = append(allItems, items...)
		fmt.Print("\rmain references fetched: ", percent, "%")
	}
	var works []work
	var passages []passage
	fmt.Println()

	for i := range allItems {
		link = allItems[i].URL
		data, _ := getContent(link)
		csvdata := string(data)
		if csvdata != "" {
			works = append(works, work{ID: allItems[i].Item, CSV: csvdata})
			reader := csv.NewReader(strings.NewReader(csvdata))
			reader.Comma = ','
			reader.LazyQuotes = true
			reader.FieldsPerRecord = 3

			for {
				line, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}
				passages = append(passages, passage{ID: line[0], Text: line[1], Label: line[2]})
			}
		}
		percent := 100 * (i + 1) / (len(allItems) - 1)
		fmt.Print("\rsub references fetched: ", percent, "%")
	}
	fmt.Println(len(passages), "passages extracted")
	fmt.Println("writing to file...")
	fmt.Println("")

	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()

	f.WriteString("identifier")
	f.WriteString("\t")
	f.WriteString("text")
	f.WriteString("\t")
	f.WriteString("label")
	f.WriteString("\n")

	greekWordRegExp := regexp.MustCompile(`\p{Greek}+`)
	latinWordRegExp := regexp.MustCompile(`\p{Latin}+`)
	arabicWordRegExp := regexp.MustCompile(`\p{Arabic}+`)
	greekwords := 0
	latinwords := 0
	arabicwords := 0

	for i := range passages {
		text := passages[i].Text
		var re = regexp.MustCompile(`<orig>[^>]+>`)
		text = re.ReplaceAllString(text, "")
		re = regexp.MustCompile(`<[^>]+>`)
		text = re.ReplaceAllString(text, "")
		re = regexp.MustCompile(`[\s\p{Zs}]{2,}`)
		text = strings.Replace(text, "\n", " ", -1)
		text = re.ReplaceAllString(text, " ")

		greekword := greekWordRegExp.FindAllString(text, -1)
		latinword := latinWordRegExp.FindAllString(text, -1)
		arabicword := arabicWordRegExp.FindAllString(text, -1)

		title := passages[i].ID
		label := passages[i].Label

		f.WriteString(title)
		f.WriteString("\t")
		f.WriteString(text)
		f.WriteString("\t")
		f.WriteString(label)
		f.WriteString("\n")

		greekwords = greekwords + len(greekword)
		latinwords = latinwords + len(latinword)
		arabicwords = arabicwords + len(arabicword)

		percent := 100 * (i + 1) / (len(passages) - 1)
		fmt.Print("\rsub references fetched: ", percent, "%")
	}
	fmt.Println("wrote", latinwords, "Latin words from", len(passages), "passages to scta.csv")
	fmt.Println("also", greekwords, "Greek words")
	fmt.Println("wrote", arabicwords, "Arabic words")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
