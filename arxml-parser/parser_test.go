package arxml

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"log"
	"testing"
)

func loadArxml(filePath string) *xmlquery.Node {
	doc, err := parseXml("resources/BE-31.arxml")
	if err != nil {
		log.Fatalln(err)
	}
	return doc
}

func TestNetwork(t *testing.T) {
	doc := loadArxml("resources/BE-31.arxml")
	networks := getNetwork(doc)
	fmt.Println(networks)
}

func TestISignal(t *testing.T) {
	doc := loadArxml("resources/BE-33.4.arxml")
	isignals := getISignal(doc)
	fmt.Println(isignals)
}

func TestCompuMethod(t *testing.T) {
	doc := loadArxml("resources/BE-31.arxml")
	datatypes := getDataTypes(doc)
	fmt.Println(datatypes)
}