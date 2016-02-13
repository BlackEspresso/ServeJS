package goquery

import (
	"bytes"

	"./../modules"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func InitPlugin() *modules.Plugin {

	p1 := modules.Plugin{
		Name: "goquery",
		Init: registerVM,
	}

	return &p1
}

type GoQueryDoc struct {
	doc *goquery.Document
}

func (g *GoQueryDoc) ExtractLinks() []string {
	links := []string{}
	g.doc.Find("a").Each(func(i int, s *goquery.Selection) {
		attrValue, ok := s.Attr("href")
		if ok {
			links = append(links, attrValue)
		}
	})
	return links
}

func (g *GoQueryDoc) ExtractAttributes(tagName string) []map[string]string {
	list := []map[string]string{}
	g.doc.Find(tagName).Each(func(i int, s *goquery.Selection) {
		n := s.Get(0)
		attrs := map[string]string{}
		for _, attr := range n.Attr {
			attrs[attr.Key] = attr.Val
		}
		list = append(list, attrs)
	})
	return list
}

func registerVM(vm *otto.Otto) otto.Value {
	obj, _ := vm.Object("({})")

	obj.Set("newDocument", func(c otto.FunctionCall) otto.Value {
		str, _ := c.Argument(0).ToString()
		b := bytes.NewBufferString(str)
		doc, _ := goquery.NewDocumentFromReader(b)
		gDoc := GoQueryDoc{}
		gDoc.doc = doc
		val, _ := c.Otto.ToValue(&gDoc)
		return val
	})

	return obj.Value()
}
