package object

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/phper-go/frame/func/conv"
)

var DefaultTypes = []string{
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"unit",
	"unit8",
	"unit16",
	"uni32",
	"unit64",
	"float32",
	"float64",
	"bool",
	"string",
	"[]bype",
}

func New(class interface{}) interface{} {
	typeOf := reflect.Indirect(reflect.ValueOf(class)).Type()
	return reflect.New(typeOf).Interface()
}

func LoadFile(obj interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	dirconfig, _ := filepath.Abs(filepath.Dir(file))
	var cmaps = make(map[string]interface{})
	var objmaps = Values(obj)
	var filemaps = make(map[string]string)
	json.Unmarshal(data, &cmaps)
	// fmt.Println(cmaps)
	for key, val := range cmaps {
		conf_file := conv.String(val)
		index := strings.Index(conf_file, ".")
		key = strings.ToLower(key)
		if -1 != index && reflect.TypeOf(objmaps[key]).String() != "string" {
			conf_file = dirconfig + "/" + conf_file
			if _, err := os.Stat(conf_file); err == nil {
				filedata, err := ioutil.ReadFile(conf_file)
				if err == nil {
					filemaps["\""+conv.String(val)+"\""] = string(filedata)
				}
			}
		}
	}
	content := string(data)
	for k, v := range filemaps {
		content = strings.Replace(content, k, v, -1)
	}
	return json.Unmarshal([]byte(content), obj)
}

func MethodMaps(class interface{}) map[string]string {

	classTypeOf := reflect.TypeOf(class)
	methodNum := classTypeOf.NumMethod()
	if methodNum == 0 {
		return nil
	}
	methodMaps := make(map[string]string)
	for i := 0; i < methodNum; i++ {
		name := classTypeOf.Method(i).Name
		methodMaps[strings.ToLower(name)] = name
	}
	return methodMaps
}

func Methods(class interface{}) []string {
	var methods []string
	maps := MethodMaps(class)
	for _, name := range maps {
		methods = append(methods, name)
	}
	return methods
}

func FindMethod(class interface{}, name string) string {
	maps := MethodMaps(class)
	method := maps[strings.ToLower(name)]
	return method
}

func Set(obj interface{}, data interface{}) error {

	objValues := Values(obj)
	values := Values(data)
	for key, val := range values {
		if v, ok := objValues[strings.ToLower(key)]; ok {
			objType := reflect.Indirect(reflect.ValueOf(v)).Kind().String()
			objValues[strings.ToLower(key)] = conv.Format(val, objType)
		}
	}
	var jsonData []byte
	var err error
	if jsonData, err = json.Marshal(objValues); err != nil {
		return err
	}
	return json.Unmarshal(jsonData, obj)

	// filedTypes := FieldTypes(obj)
	// for key, val := range objData {

	// 	if objType, ok := filedTypes[strings.ToLower(key)]; ok {
	// 		if val == nil {
	// 			delete(objData, key)
	// 		} else {
	// 			objData[key] = conv.Format(val, objType)
	// 		}
	// 	}
	// }
	// if jsonData, err = json.Marshal(objData); err != nil {
	// 	return err
	// }
	// return json.Unmarshal(jsonData, obj)
}

func SetField(obj interface{}, key string, value interface{}) error {
	var maps = map[string]interface{}{key: value}
	return Set(obj, maps)
}

func Fields(object interface{}) []string {

	var fields []string
	typeOf := reflect.Indirect(reflect.ValueOf(object)).Type()
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Anonymous == false {
			name := typeOf.Field(i).Name
			if unicode.IsUpper(rune(name[0])) {
				fields = append(fields, name)
			}
		}
	}
	return fields
}

func FieldTypes(obj interface{}) map[string]string {

	typeOf := reflect.Indirect(reflect.ValueOf(obj)).Type()
	types := fieldTypes(typeOf)
	return types
}

func fieldTypes(typeOf reflect.Type) map[string]string {

	types := make(map[string]string)

	for i := 0; i < typeOf.NumField(); i++ {
		name := typeOf.Field(i).Name

		if unicode.IsUpper(rune(name[0])) {
			fieldType := typeOf.Field(i).Type
			if typeOf.Field(i).Anonymous {
				maps := fieldTypes(fieldType)
				if len(maps) > 0 {
					for key, val := range maps {
						types[key] = val
					}
				}
			} else {
				types[strings.ToLower(name)] = fieldType.String()
			}
		}
	}
	return types
}

func Inputs(form interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	values := Values(form)
	for key, val := range values {
		if val != nil {
			fields[key] = val
		}
	}
	return fields
}

func Values(obj interface{}) map[string]interface{} {

	valMaps := make(map[string]interface{})
	valueOf := reflect.Indirect(reflect.ValueOf(obj))
	typeOf := valueOf.Type()

	switch valueOf.Kind().String() {
	case "map":
		if jsonData, err := json.Marshal(obj); err == nil {
			json.Unmarshal(jsonData, &valMaps)
		}
	case "struct":
		for i := 0; i < typeOf.NumField(); i++ {
			if unicode.IsUpper(rune(typeOf.Field(i).Name[0])) {
				if typeOf.Field(i).Anonymous == false {
					valMaps[strings.ToLower(typeOf.Field(i).Name)] = valueOf.Field(i).Interface()
				} else {
					vals := Values(valueOf.Field(i).Interface())
					for k, v := range vals {
						valMaps[k] = v
					}
				}
			}
		}
	}
	return valMaps
}

func ToPtr(obj interface{}) interface{} {

	objType := reflect.TypeOf(obj)
	if objType.Kind().String() == "ptr" {
		return obj
	} else {
		temp := reflect.New(objType)
		temp.Elem().Set(reflect.ValueOf(obj))
		return temp.Interface()
	}
}

func Name(obj interface{}) string {
	return reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
}

func Tags(obj interface{}, tagname string) map[string]string {

	tags := make(map[string]string)
	typeOf := reflect.Indirect(reflect.ValueOf(obj)).Type()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tags[strings.ToLower(field.Name)] = field.Tag.Get(tagname)
	}
	return tags
}
