package controllers

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PageController struct{
	tpl *template.Template
}

func CreateNewPageController(tpl *template.Template) *PageController {
	return &PageController{tpl}
}

func (p *PageController) Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p.tpl.ExecuteTemplate(w, "index.gohtml", nil)
}