package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Response struct {
	Message string
}

type Request struct {
	Name string
}

type HelloRequest struct{}

func (h *HelloRequest) Hello(reqa Request, res *Response) (err error) {

	res.Message = "Exentsion call executed"
	return
	// Sends the email using the elasitmail
	// uses the elastic mail configurations for it
	mailClient := http.Client{
		Timeout: time.Minute * 1,
	}

	color.Green("elastic mail api key", "key-2fcff6a1835d0be445f2abb3f1c9adf5")

	//?apikey="+Config.Mail.APIKey
	myurl := url.URL{}
	myurl.Scheme = "https"
	myurl.Host = "api.elasticemail.com"
	myurl.Path = "v2/email/send"

	q := myurl.Query()
	q.Add("apikey", "key-2fcff6a1835d0be445f2abb3f1c9adf5")
	q.Add("to", "m.champaneri.20@gmail.com")
	q.Add("subject", "Test mail")
	q.Add("bodyHtml", "<h1> Test is Test</h1>")
	q.Add("charset", "UTF-8")
	q.Add("from", "support@puberstreet.com")
	q.Add("fromName", "Cobra stack")

	myurl.RawQuery = q.Encode()
	req, err := http.NewRequest("POST", myurl.String(), strings.NewReader(string("a")))

	if err != nil {
		color.Red("mail sent with elastic mail")
		return
	}
	_, errDO := mailClient.Do(req)
	if errDO != nil {
		color.Red("error on do request ", errDO.Error())
		return
	}

	color.Green("mail sent with elastic mail")
	return
	//fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
