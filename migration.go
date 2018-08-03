package main

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const filesDirectory = "json/"
const authorsSongPath = "authors.song.json"
const authorsTangPath = "authors.tang.json"

type Author struct {
	ID         uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Desc       string `gorm:"type:text" json:"desc"`
	Author     string `json:"name"`
	AuthorType string
}

type Poetry struct {
	ID         uint     `gorm:"primary_key;AUTO_INCREMENT"`
	Author     string   `json:"author"`
	Paragraphs []string `json:"paragraphs"`
	Strains    []string `json:"strains"`
	Title      string   `json:"title"`
	PoetryType string
}

func check(e error) {
	if nil != e {
		panic(e)
	}
}

func ReadAuthor(filePath string) []*Author {
	file, err := os.Open(filePath)
	check(err)

	defer file.Close()

	var authors []*Author
	err = json.NewDecoder(file).Decode(&authors)
	check(err)

	return authors
}

func ReadPoetry(filePath string) []*Poetry {
	file, err := os.Open(filePath)
	check(err)

	defer file.Close()

	var poetrys []*Poetry
	err = json.NewDecoder(file).Decode(&poetrys)
	check(err)

	return poetrys
}

func CreateAuthor(db *gorm.DB) {
	if !db.HasTable(&Author{}) {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(&Author{}).Error
		check(err)
	}
}

func main() {
	// 连接数据库
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:9001)/peotry?charset=utf8mb4")
	check(err)
	defer db.Close()

	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	CreateAuthor(db)

	authors := ReadAuthor(filesDirectory + authorsSongPath)

	for _, author := range authors {
		author.AuthorType = "SONG"
		db.Create(author)
		fmt.Println(author.ID)
	}
}
