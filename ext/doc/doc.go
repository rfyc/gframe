package doc

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/web/core"
	"github.com/phper-go/frame/web/route"
)

type controller struct {
	core.Controller
}

func (this *controller) IndexAction() {

	var keys []string
	for uri, obj := range route.Controllers {
		if uri != "/doc" {
			if _, ok := obj.(interfaces.Action); ok {
				keys = append(keys, uri)
			} else if _, ok := obj.(interfaces.Controller); ok {
				for _, method := range object.Methods(obj) {
					if len(method) > 6 && method[len(method)-6:] == "Action" {
						action := uri + "/" + strings.ToLower(method[:len(method)-6])
						keys = append(keys, action)
					}
				}
			}
		}
	}
	sort.Strings(keys)
	var skeys = make(map[string][]string)
	for i := len(keys) - 1; i >= 0; i-- {
		m := strings.Split(strings.TrimRight(keys[i], "/"), "/")
		k := strings.Join(m[:len(m)-1], "/")
		skeys[k] = append(skeys[k], keys[i])
	}
	fmt.Println(skeys)
	tpl, err := template.New("index").Parse(html_index)
	fmt.Println(err)
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, map[string]interface{}{"keys": skeys})
	fmt.Println(err)
	this.Output().Content = buf.Bytes()
}

func (this *controller) InfoAction() {

}

func (this *controller) DoAction() {

}

func init() {

	core.RegisterController("/doc", &controller{})
}
