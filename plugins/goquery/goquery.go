package goquery

import (
	"bytes"

	"golang.org/x/net/html"

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

func (g *GoQueryDoc) ExtractForms() []map[string]interface{} {
	list := []map[string]interface{}{}
	g.doc.Find("form").Each(func(i int, s *goquery.Selection) {
		node := s.Get(0)
		attrs := getAttrs(node)
		cAttrs := upCastList(attrs)
		list = append(list, cAttrs)
		inputList := []map[string]interface{}{}
		s.Find("input").Each(func(i int, s *goquery.Selection) {
			inputNode := s.Get(0)
			attrsInput := getAttrs(inputNode)
			inputList = append(inputList, upCastList(attrsInput))
		})
		cAttrs["inputFields"] = inputList
	})
	return list
}

func upCastList(list map[string]string) map[string]interface{} {
	castAttrs := map[string]interface{}{}
	for k, v := range list {
		castAttrs[k] = v
	}
	return castAttrs
}

func getAttrs(node *html.Node) map[string]string {
	attrs := map[string]string{}
	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}
	return attrs
}

func (g *GoQueryDoc) ExtractAttributes(tagName string) []map[string]string {
	list := []map[string]string{}
	g.doc.Find(tagName).Each(func(i int, s *goquery.Selection) {
		node := s.Get(0)
		attrs := getAttrs(node)
		list = append(list, attrs)
	})
	return list
}

func registerVM(vm *modules.JsVm) otto.Value {
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
