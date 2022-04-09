package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

func getBibs(c *gin.Context) {
	formatter, _ := c.GetQuery("f")

	for i := 0; i < len(bibs); i++ {
		bibs[i].Format(formatter, i+1)
	}

	c.IndentedJSON(http.StatusOK, bibs)
}

func postBib(c *gin.Context) {
	var newBib Bibliography

	if err := c.Bind(&newBib); err != nil {
		return
	}

	bibs = append(bibs, newBib)
	c.IndentedJSON(http.StatusCreated, newBib)
}

func getBibByID(c *gin.Context) {
	id := c.Param("id")
	formatter, _ := c.GetQuery("f")

	for i, b := range bibs {
		if b.ID == id {
			b.Format(formatter, i+1)
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bibliography not found"})
}

func main() {
	router := gin.Default()

	router.GET("/bibs", getBibs)
	router.POST("/bibs", postBib)
	router.GET("/bibs/:id", getBibByID)

	router.Run("localhost:8080")
}
