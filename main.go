package main

/**
What has been done

*/
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

func handleHome(req badnet.Request) badnet.Response {

	data := []byte("hello world")

	return badnet.Response{
		ResponseMsg:  "OK",
		ResponseCode: 200,
		Version:      badnet.V1_1,
		Headers: map[string]string{
			badnet.ContentType:   "text/plain",
			badnet.ContentLength: fmt.Sprintf("%d", len(data)),
		},
		Data: data,
	}
}

func main() {
	logger := badlogger.DefaultLogger()
	badnet.GET.RegisterPath("/", handleHome)
	config := badnet.ServerConfiguration{Port: ":8080", Logger: &logger}
	badnet.StartServer(config)
}
