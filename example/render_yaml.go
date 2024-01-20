package main //nolint:typecheck

import (
	"log"
	"os"

	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
)

func main() {
	newObject := []Object{
		{Name: "nikhil", Date: "01-01-2024"},
		{Name: "john", Date: "01-02-2024"},
	}

	logger := logrus.New()
	render := renderer.GetRenderer(os.Stdout, logger, true, false, false, false)

	if err := render.Render(newObject); err != nil {
		log.Fatal(err)
	}
}

/*
The above code should generate below json
---
- Date: 01-01-2024
  Name: nikhil
- Date: 01-02-2024
  Name: john
*/
