package main

import "fmt"

func (e *ECOM) Install(r Request, w *Response) (err error) {
	// To be defined
	return
}

func (e *Admin) HookInMenu(r Request, w *Response) (err error) {
	w.Output = fmt.Sprint(`
	<a class="navbar-item" href="/admin/posts/list">
	Products
	</a>
	`)
	w.Status = "success"
	w.Type = "HTML"
	return
}
