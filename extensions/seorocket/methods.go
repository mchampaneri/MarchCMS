package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func (e *SEORocket) Install(r Request, w *Response) (err error) {
	// To be defined
	return
}

func (e *SEORocket) ContactUS(r Request, w *Response) (err error) {
	w.Output = fmt.Sprintf(`
		<div>
			<form action="/extension" method="post">
				<label> Email Address </label>
				<input type="text" name="email"/>
				<input type="hidden" name="extname" value="seorocket">
				<input type="hidden" name="extmethod" value="SEORocket.SendMail">
				<input type="hidden" name="redirectURL" value="%s">
				<input type="submit" value="Send"/>
			</form>
		</div>
	`, r.Input["stringInput"])
	w.Status = "success"
	w.Type = "HTML"
	return
}

func (e *SEORocket) SendMail(r Request, w *Response) (err error) {
	formValues := r.Input["form"].(url.Values)
	log.Println(formValues)
	w.Output = formValues.Get("email")
	w.Status = "success"
	w.Type = "HTML"
	return
}

func (e *SEORocket) SEOMeta(r Request, w *Response) (err error) {
	// Reading Input keywords
	splits := strings.Split(r.Input["stringInput"].(string), ",")
	keywords := splits[0]
	description := splits[1]
	SEOMetaTags := fmt.Sprintf(`<meta name="keywords" value="%s" /><meta name="description" value="%s" />`, keywords, description)
	w.Output = SEOMetaTags
	w.Status = "success"
	w.Type = "HTML"
	return
}
