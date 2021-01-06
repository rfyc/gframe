package std

import (
	"encoding/json"
	"fmt"

	"github.com/phper-go/frame/func/conv"
)

type T struct {
	value interface{}
	err   error
	array *Array
	maps  *Maps
}

func (this *T) String() string {
	return conv.String(this.value)
}

func (this *T) Int() int {
	return conv.Int(this.value)
}

func (this *T) Int8() int8 {
	return conv.Int8(this.value)
}

func (this *T) Int16() int16 {
	return conv.Int16(this.value)
}

func (this *T) Int32() int32 {
	return conv.Int32(this.value)
}

func (this *T) Int64() int64 {
	return conv.Int64(this.value)
}

func (this *T) Uint() uint {
	return conv.Uint(this.value)
}

func (this *T) Uint8() uint8 {
	return conv.Uint8(this.value)
}

func (this *T) Uint16() uint16 {
	return conv.Uint16(this.value)
}

func (this *T) Uint32() uint32 {
	return conv.Uint32(this.value)
}

func (this *T) Uint64() uint64 {
	return conv.Uint64(this.value)
}

func (this *T) Float32() float32 {
	return conv.Float32(this.value)
}

func (this *T) Float64() float64 {
	return conv.Float64(this.value)
}

func (this *T) Bool() bool {
	return conv.Bool(this.value)
}

func (this *T) Bytes() []byte {
	return conv.Bytes(this.value)
}

func (this *T) Strings() []string {
	return conv.Strings(this.value)
}

func (this *T) Array() *Array {
	if this.array == nil {
		this.array = NewArray(this.Bytes())
	}
	return this.array
}

func (this *T) Maps() *Maps {
	if this.maps == nil {
		this.maps = NewMaps(this.Bytes())
	}
	return this.maps
}

func (this *T) Index(index int) *T {
	return this.Array().Index(index)
}

func (this *T) Item(key string) *T {
	return this.Maps().Item(key)
}

func (this *T) IsNil() bool {
	return this.value == nil
}

func (this *T) IsEmpty() bool {

	if this.IsNil() {
		return true
	}
	value := this.String()
	if value == "" || value == "0" || value == "false" {
		return true
	}
	return false
}

func (this *T) Error() error {
	return this.err
}

func (this *T) Println() {
	fmt.Println(this.String())
}

type Array struct {
	array []interface{}
	len   int
	err   error
}

func (this *Array) Index(index int) *T {

	if this.Len() > index {
		item := this.array[index]
		return NewT(item)
	}
	var value interface{}
	return NewT(value)
}

func (this *Array) Len() int {
	if this.len > 0 {
		return this.len
	}
	this.len = len(this.array)
	return this.len
}

func (this *Array) All() []*T {
	var items []*T
	for _, value := range this.array {
		items = append(items, NewT(value))
	}
	return items
}
func (this *Array) Error() error {
	return this.err
}

func (this *Array) Println() {
	fmt.Println("error:", this.err)
	for index, one := range this.All() {
		fmt.Println("index", index, ":", one)
	}
}

type Maps struct {
	maps map[string]interface{}
	err  error
}

func (this *Maps) Item(key string) *T {
	if item, ok := this.maps[key]; ok {
		return NewT(item)
	}
	var value interface{}
	return NewT(value)
}

func (this *Maps) Len() int {
	return len(this.maps)
}

func (this *Maps) All() map[string]*T {
	var items map[string]*T
	for key, value := range this.maps {
		items[key] = NewT(value)
	}
	return items
}

func (this *Maps) Error() error {
	return this.err
}

func (this *Maps) Println() {
	fmt.Println("error:", this.err)
	for key, value := range this.All() {
		fmt.Println("key", key, ":", value)
	}
}

func NewArray(bytes []byte) *Array {
	array := &Array{}
	array.err = json.Unmarshal(bytes, &array.array)
	return array
}

func NewMaps(bytes []byte) *Maps {
	maps := &Maps{}
	maps.err = json.Unmarshal(bytes, &maps.maps)
	return maps
}

func NewT(value interface{}, err ...error) *T {
	t := &T{
		value: value,
	}
	if len(err) > 0 {
		t.err = err[0]
	}
	return t
}
