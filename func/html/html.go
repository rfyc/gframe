package html

import (
	"html"
)

func Encode(s string) string {
	return html.EscapeString(s)
}

func Decode(s string) string {
	return html.UnescapeString(s)
}

func Tag(tagname, text string, isclosed bool, htmlOptions ...map[string]string) string {

	var options map[string]string
	if len(htmlOptions) > 0 {
		options = htmlOptions[0]
	} else {
		options = make(map[string]string)
	}

	content := "<"
	content += tagname
	for key, value := range options {
		content += " " + key + "=\"" + value + "\""
	}
	content += ">"
	if len(text) > 0 {
		content += text
	}
	if isclosed {
		content += "</" + tagname + ">"
	} else {
		content += " />"
	}

	return content
}

// func ConvList(list []map[string]string) []map[string]template.HTML {
// 	var volist []map[string]template.HTML
// 	for _, row := range list {
// 		one := make(map[string]template.HTML)
// 		for key, value := range row {
// 			one[key] = template.HTML(value)
// 		}
// 		volist = append(volist, one)
// 	}
// 	return volist
// }

func Link(text, url string, htmlOptions ...map[string]string) string {

	var options map[string]string
	if len(htmlOptions) > 0 {
		options = htmlOptions[0]
	} else {
		options = make(map[string]string)
	}
	options["href"] = url

	return Tag("a", text, true, options)
}

func Br() string {
	return Tag("br", "", false)
}
