// nolint:typecheck
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nikhilsbhat/common/content"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	fileData, err := os.ReadFile("../fixtures/sample.yaml")
	if err != nil {
		log.Fatal(err)
	}

	obj := content.Object(fileData)
	fileType := obj.CheckFileType(logger)

	fmt.Println(fileType)
	// Above would identify the file content as yaml.
}
