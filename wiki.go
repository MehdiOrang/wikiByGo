package main

import (
    "fmt"
    "os"
)

type Page struct{
    title string
    Body  []byte
}

func (p *Page) save() error{
    source := p.title
    return os.WriteFile(source, p.body, 0600)
}


