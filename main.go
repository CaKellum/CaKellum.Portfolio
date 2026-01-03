package main

/**
TODO:

	* build out templates using HTMX for front end
		- learn HTMX
		- learn how template html in golang

	* What do i want the website to look like? will need to sit down and think about that for a second
	* build out enpoints and the corresponding functionality of them
		- what do the request form htmx look like
*/

import (
	"fmt"
	"os"
	"text/template"

	"com.kellum.portfolio/badlogger"
	"com.kellum.portfolio/badnet"
)

type LanguageInfo struct {
	Color string
	Name  string
}

type ProjectInfo struct {
	Name      string
	Link      string
	Languages []LanguageInfo
}

type BasePageProperties struct {
	Title       string
	Content     string
	CurrentYear string
}

type FileBuffer struct{ Buffer []byte }

func (fb *FileBuffer) Write(p []byte) (n int, err error) {
	if len(fb.Buffer) == 0 {
		fb.Buffer = p
	} else {
		fb.Buffer = append(fb.Buffer, p...)
	}
	return len(fb.Buffer), nil
}

func getBaseTemplate(content string) ([]byte, error) {
	props := BasePageProperties{
		Title:       "CaKellum",
		Content:     content,
		CurrentYear: "2025",
	}
	tmpl := template.Must(template.ParseFiles("templates/base.tmpl"))
	fb := FileBuffer{}
	if err := tmpl.Execute(&fb, props); err != nil {
		return make([]byte, 0), err
	}
	return fb.Buffer, nil
}

func handleHome(req badnet.Request) badnet.Response {
	fileData, fileErr := os.ReadFile("templates/about_content.tmpl")
	if fileErr != nil {
		errStr := fmt.Sprint(fileErr)
		data := []byte(errStr)
		return badnet.Response{
			ResponseMsg:  "Internal Server Error",
			ResponseCode: 500,
			Version:      badnet.V1_1,
			Headers: map[string]string{
				badnet.ContentType:   "text/plain",
				badnet.ContentLength: fmt.Sprintf("%d", len(data)),
			},
			Data: data,
		}

	}
	data, err := getBaseTemplate(string(fileData))
	if err != nil {
		errStr := fmt.Sprint(err)
		data := []byte(errStr)
		return badnet.Response{
			ResponseMsg:  "Internal Server Error",
			ResponseCode: 500,
			Version:      badnet.V1_1,
			Headers: map[string]string{
				badnet.ContentType:   "text/plain",
				badnet.ContentLength: fmt.Sprintf("%d", len(data)),
			},
			Data: data,
		}
	}
	return badnet.Response{
		ResponseMsg:  "OK",
		ResponseCode: 200,
		Version:      badnet.V1_1,
		Headers: map[string]string{
			badnet.ContentType:   "text/html",
			badnet.ContentLength: fmt.Sprintf("%d", len(data)),
		},
		Data: data,
	}
}

func handleAbout(req badnet.Request) badnet.Response {
	data, err := getBaseTemplate("<p>This is the about page</p>")
	if err != nil {
		errStr := fmt.Sprint(err)
		data := []byte(errStr)
		return badnet.Response{
			ResponseMsg:  "Internal Server Error",
			ResponseCode: 500,
			Version:      badnet.V1_1,
			Headers: badnet.HTTPHeaders{
				badnet.ContentType:   "text/plain",
				badnet.ContentLength: fmt.Sprintf("%d", len(data)),
			},
			Data: data,
		}
	}
	return badnet.Response{
		ResponseMsg:  "OK",
		ResponseCode: 200,
		Version:      badnet.V1_1,
		Headers: badnet.HTTPHeaders{
			badnet.ContentType:   "text/html",
			badnet.ContentLength: fmt.Sprintf("%d", len(data)),
		},
		Data: data,
	}
}

func handleCSS(req badnet.Request) badnet.Response {
	file, err := os.ReadFile("templates/static/main.css")
	if err != nil {
		errStr := fmt.Sprint(err)
		data := []byte(errStr)
		return badnet.Response{
			ResponseMsg:  "Internal Server Error",
			ResponseCode: 500,
			Version:      badnet.V1_1,
			Headers: badnet.HTTPHeaders{
				badnet.ContentType:   "text/plain",
				badnet.ContentLength: fmt.Sprintf("%d", len(data)),
			},
			Data: data,
		}
	}
	return badnet.Response{
		ResponseMsg:  "Ok",
		ResponseCode: 200,
		Version:      badnet.V1_1,
		Headers: badnet.HTTPHeaders{
			badnet.ContentType:   "text/css",
			badnet.ContentLength: fmt.Sprintf("%d", len(file)),
		},
		Data: file,
	}
}

func main() {
	logger := badlogger.DefaultLogger()
	badnet.GET.RegisterPath("/", handleHome)
	badnet.GET.RegisterPath("/about", handleAbout)
	badnet.GET.RegisterPath("/static/main.css", handleCSS)
	config := badnet.ServerConfiguration{Network: "tcp", Port: ":8080", Logger: &logger}
	badnet.StartServer(config)
}
