package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jiandahao/jsonplus"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DocDetails struct {
	Name string `json:"name" yaml:"name"`
	Url  string `json:"url" yaml:"url"`
	//ExternalDocs struct {
	//	Description string `json:"description" yaml:"description"`
	//	Url         string `json:"url" yaml:"url"`
	//} `json:"external_docs" yaml:"external_docs"`
}

//"externalDocs":{"description":"Find out more about our store","url":"http://swagger.io"}
var (
	templ *template.Template
	docs  = new(Documents)
	port = 8088
)

type Documents struct {
	DocList []*DocDetails
	mux     sync.Mutex
}

func (d *Documents) uniqueAppend(document *DocDetails) {
	d.mux.Lock()
	defer d.mux.Unlock()
	for _, doc := range d.DocList {
		if doc.Name == document.Name && doc.Url == document.Url {
			return
		}
	}
	d.DocList = append(d.DocList, document)
}

type Filter func(filename string) bool

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, filter Filter) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth+PthSep+fi.Name(), filter)
		} else {
			// 过滤指定格式
			//ok := strings.HasSuffix(fi.Name(), ".txt")
			if filter != nil {
				if ok := filter(fi.Name()); ok {
					files = append(files, dirPth+PthSep+fi.Name())
				}
			} else {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table, filter)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

const DocumentDirPath = "docs"

func init() {
	var err error
	templ, err = template.ParseFiles("template/index.tmpl")
	if err != nil {
		panic(err)
	}
	// TODO： create folder if docs is not existed

	//docs.DocList = []*DocDetails{
	//	{
	//		Name: "test1",
	//		Url:  "/docs/test_swagger.swagger.json",
	//	},
	//	{
	//		Name: "test2",
	//		Url:  "https://petstore.swagger.io/v2/swagger.json",
	//	},
	//}

	if !IsExist(DocumentDirPath) {
		if err := os.Mkdir(DocumentDirPath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	filePaths, err := GetAllFiles(DocumentDirPath, func(filename string) bool {
		return strings.HasSuffix(filename, ".swagger.json")
	})

	if err != nil {
		return
	}

	var j *jsonplus.JsonPlus
	for _, filePath := range filePaths {
		fmt.Println(filePath)
		j, err = jsonplus.New(filePath)
		if err != nil {
			panic(err)
		}

		title, err := j.Get("info.title")
		if err != nil {
			panic(title)
		}

		docs.DocList = append(docs.DocList, &DocDetails{
			Name: fmt.Sprintf("%s", title),
			Url:  fmt.Sprintf("/%s", filePath),
		})
	}
}

func main() {
	router := gin.Default()

	router.GET("/", GetIndexPage)
	router.POST("/new", AddDocuments)
	router.Static("/docs", "./docs")
	router.Static("/dist", "./dist")
	router.Run(fmt.Sprintf(":%v", port))
	//http.ListenAndServe(":8088", NewHandler(router))
}

func GetIndexPage(c *gin.Context) {
	// render index.html
	if err := templ.Execute(c.Writer, docs); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Failed to render template with error:%s",err.Error()))
	}
}

func AddDocuments(c *gin.Context) {
	// reader file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if !strings.HasSuffix(header.Filename, ".json") {
		// TODO：如果是yaml则转换成json， 如果为proto文件则需要转换成openapi 3.0 json
		c.String(http.StatusUnsupportedMediaType, fmt.Sprintf("Unsupported file type"))
		return
	}

	// save file
	out, err := os.Create("docs/" + filepath.Base(header.Filename))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err.Error()))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)

	j, err := jsonplus.New(out.Name())
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	title, err := j.Get("info.title")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	go docs.uniqueAppend(&DocDetails{
		Name: fmt.Sprintf("%s",title),
		Url:  fmt.Sprintf("/docs/%s", header.Filename),
	})

	c.String(http.StatusCreated, "upload successful \n")
	//否则如果为.proto文件，先转换成swagger json，然后转换成openapi 3.0,最后保存
}
