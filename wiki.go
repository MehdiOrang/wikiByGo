package main

import (
    "fmt"
    "html/template"
    "net/http"
    "os"
    "log"
)

type Page struct{
    Title string
    Body  []byte
}
var templates = template.Must(template.ParseFiles("view.html"))
func (p *Page) save() error{
    source := p.Title
    return os.WriteFile(source, p.Body, 0600)
}

func load(name string) (*Page, error) {
    source := name
    body, err := os.ReadFile(source)

    if err != nil {
	return nil, err
    }

    return &Page{Title : name, Body: body}, nil
}


func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	page, err  := load(title)
	if err != nil{
	    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	    return
	}
	renderTemplate(w, "view", page)
}

func editHandler(w http.ResponseWriter, r * http.Request){
    source := r.URL.Path[len("/edit/"):]
    p, err := load(source)
    if err != nil {
	p = &Page{Title: source}
    }

    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
          "<textarea name=\"body\">%s</textarea><br>"+
          "<input type=\"submit\" value=\"save\">"+
        "</form>",
        p.Title, p.Title, p.Body)

}

func saveHandler(w http.ResponseWriter, r *http.Request){
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
     err := templates.ExecuteTemplate (w,tmpl + ".html", p)
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main(){
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/save/", saveHandler)
    http.HandleFunc("/edit/", editHandler)
    log.Fatal(http.ListenAndServe(":8080",nil))
}
