package view

import (
	"io/ioutil"
	"log"
	"net/url"

	"github.com/zserge/lorca"
)

var UI lorca.UI

func InitUI() {
	//Inicializo
	UI, _ = lorca.New("", "", 800, 700)
	//Cargo la primera vista
	ChangeView("www/login.html")
}

func ChangeView(nombreVista string) {
	content, err := ioutil.ReadFile(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	myFileContents := string(content)

	loadableContents := "data:text/html," + url.PathEscape(myFileContents)
	UI.Load(loadableContents)
}
