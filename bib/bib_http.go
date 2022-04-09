package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Bibliography 参考文献
type Bibliography struct {
	ID        string
	Author    string
	Title     string
	Type      string
	Publisher string
	Reference string

	Comment   string
	Formatted string
}

// Format 用指定格式格式化参考文献，将结果填充到 Formatted 字段
func (b *Bibliography) Format(formatter string, idx int) {
	switch formatter {
	case "gbt":
		b.Formatted = fmt.Sprintf("[%v] %v. %v[%v]. %v:%v.",
			idx, b.Author, b.Title, b.Type, b.Publisher, b.Reference)
	default:
		b.Formatted = ""
	}
}

// bibs 是初始化的文献
var bibs = []Bibliography{
	{
		ID:        "csapp",
		Author:    "Randal E. Bryant and David R. O'Hallaron",
		Title:     "Computer Systems: A Programmer's Perspective, 3rd Edition",
		Type:      "M",
		Publisher: "Pearson,2016",
		Reference: "19-22",
		Comment:   "圣经：深入理解计算机系",
	},
	{
		ID:        "sicp",
		Author:    "Hal Abelson and Jerry Sussman and Julie Sussman",
		Title:     "Structure and Interpretation of Computer Programs",
		Type:      "M",
		Publisher: "MIT Press, 1984",
		Reference: "23-30",
		Comment:   "魔法书：计算机程序的构造和解释",
	},
	{
		ID:        "go",
		Author:    "左书祺(@Draven)",
		Title:     "Go语言设计与实现",
		Type:      "M",
		Publisher: "北京:人民邮电出版社,2021",
		Reference: "181-184",
		Comment:   "我还没读完",
	},
}

// GET /bibs[?f=gbt] 获取所有参考文献，若指定 f=gbt 则格式化为 GB/T 7714 格式
func getBibs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())

	formatter := r.URL.Query().Get("f")

	for i := 0; i < len(bibs); i++ {
		bibs[i].Format(formatter, i+1)
	}

	j, _ := json.MarshalIndent(bibs, "", "    ")
	w.Write(j)
}

// POST /bibs 创建新的参考文献，添加到 bibs
func postBib(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.String())

	var newBib Bibliography

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &newBib); err != nil {
		return
	}

	bibs = append(bibs, newBib)

	j, _ := json.MarshalIndent(newBib, "", "    ")
	w.Write(j)
}

// GET /bibs/:id 获取指定 id 的一篇参考文献，若指定 f=gbt 则格式化为 GB/T 7714 格式
func getBibByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/bibs/"):]
	formatter := r.URL.Query().Get("f")

	log.Println(r.Method, r.URL.String(), fmt.Sprintf("id=%v", id))

	for i, b := range bibs {
		if b.ID == id {
			b.Format(formatter, i+1)
			j, _ := json.MarshalIndent(b, "", "    ")
			w.Write(j)
			return
		}
	}

	j, _ := json.MarshalIndent(
		map[string]string{"message": "bibliography not found"},
		"", "    ")
	w.Write(j)
}

func handleBibs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBibs(w, r)
	case http.MethodPost:
		postBib(w, r)
	}
}

func main() {
	http.HandleFunc("/bibs", handleBibs)
	http.HandleFunc("/bibs/", getBibByID)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
