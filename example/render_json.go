package main

import (
	"log"
	"os"

	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
)

type Object struct {
	Name string
	Date string
}

func main() {
	newObject := []Object{
		{Name: "nikhil", Date: "01-01-2024"},
		{Name: "john", Date: "01-02-2024"},
	}

	logger := logrus.New()
	render := renderer.GetRenderer(os.Stdout, logger, false, true, false, false)

	if err := render.Render(newObject); err != nil {
		log.Fatal(err)
	}
}

/*
The above code should generate below json
[
     {
          "Name": "nikhil",
          "Date": "01-01-2024"
     },
     {
          "Name": "john",
          "Date": "01-02-2024"
     }
]
*/
