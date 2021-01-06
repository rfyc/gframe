package db

import (
	"context"
	"strconv"
	"strings"

	"github.com/phper-go/frame/func/conv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

type DBCommand struct {
	pdo       PDO
	context   context.Context
	dbmodel   interfaces.DBModel
	table     string
	field     string
	condition string
	limit     int
	offset    int
	order     string
	group     string
	join      string
	leftjoin  string
	rightjoin string
	having    string
	where     string
	args      []interface{}
}

func (this *DBCommand) SetPDO(pdo PDO) *DBCommand {
	this.pdo = pdo
	return this
}
func (this *DBCommand) SetContext(ctx context.Context) *DBCommand {
	this.context = ctx
	return this
}
func (this *DBCommand) Table(table string) *DBCommand {
	this.table = table
	return this
}

func (this *DBCommand) Select(field string) *DBCommand {
	this.field = field
	return this
}

func (this *DBCommand) Order(order string) *DBCommand {
	this.order = order
	return this
}

func (this *DBCommand) Group(group string) *DBCommand {
	this.group = group
	return this
}

func (this *DBCommand) Having(having string) *DBCommand {
	this.having = having
	return this
}

func (this *DBCommand) Limit(limit int) *DBCommand {
	this.limit = limit
	return this
}

func (this *DBCommand) Offset(offset int) *DBCommand {
	this.offset = offset
	return this
}

func (this *DBCommand) Join(join string, args ...interface{}) *DBCommand {

	return this
}

func (this *DBCommand) LeftJoin(leftjoin string, args ...interface{}) *DBCommand {

	return this
}

func (this *DBCommand) RightJoin(rightjoin string, args ...interface{}) *DBCommand {

	return this
}

func (this *DBCommand) Where(where string, args ...interface{}) *DBCommand {
	if len(this.where) == 0 {
		this.where = where
	} else {
		this.where += " and  " + where
	}
	count := len(args)
	if count > 0 {
		for i := 0; i < count; i++ {
			this.args = append(this.args, args[i])
		}
	}
	return this
}

func (this *DBCommand) Sql() string {
	sql := ""
	sql += "select "
	if len(this.field) > 0 {
		sql += this.field
	} else {
		sql += "*"
	}
	sql += " from " + this.table
	if len(this.join) > 0 {
		sql += " join " + this.join
	}
	if len(this.leftjoin) > 0 {
		sql += " left join " + this.leftjoin
	}
	if len(this.rightjoin) > 0 {
		sql += " right join " + this.rightjoin
	}

	if len(this.where) > 0 {
		sql += " where " + this.where
	}

	if len(this.group) > 0 {
		sql += " group by " + this.group
	}
	if len(this.having) > 0 {
		sql += " having " + this.having
	}
	if len(this.order) > 0 {
		sql += " order by " + this.order
	}
	if this.limit > 0 {
		sql += " limit " + strconv.Itoa(this.limit)
	}
	if this.offset > 0 {
		sql += " offset " + strconv.Itoa(this.offset)
	}
	return sql
}

func (this *DBCommand) Clone() *DBCommand {

	cmd := &DBCommand{}
	cmd.pdo = this.pdo
	cmd.table = this.table
	cmd.field = this.field
	cmd.condition = this.condition
	cmd.limit = this.limit
	cmd.offset = this.offset
	cmd.order = this.order
	cmd.group = this.group
	cmd.join = this.join
	cmd.leftjoin = this.leftjoin
	cmd.rightjoin = this.rightjoin
	cmd.having = this.having
	cmd.where = this.where
	cmd.args = this.args
	return cmd
}

func (this *DBCommand) QueryRow() (map[string]string, error) {

	var row = make(map[string]string)
	limit := this.limit
	this.Limit(1)
	result, err := this.QueryRows()
	if err != nil {
		return row, err
	}
	this.Limit(limit)
	if len(result) > 0 {
		row = result[0]
	}
	return row, err
}

func (this *DBCommand) QueryBind(obj interface{}) error {

	row, err := this.QueryRow()
	if err == nil {
		object.Set(obj, row)
	}
	return err
}

func (this *DBCommand) QueryBinds(objs []interface{}) error {
	rows, err := this.QueryRows()
	if err != nil {
		return err
	}
	object.Set(objs, rows)
	return err
}

func (this *DBCommand) QueryRows() ([]map[string]string, error) {

	var result []map[string]string
	rows, err := this.pdo.Query(this.Sql(), this.args...)
	if err != nil {
		return result, err
	}
	columns, err := rows.Columns()
	values := make([]string, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}

	for rows.Next() {
		_ = rows.Scan(scans...)
		each := make(map[string]string)
		for i, col := range values {
			each[columns[i]] = col
		}
		result = append(result, each)
	}
	rows.Close()
	return result, err
}

func (this *DBCommand) QueryCount() (int, error) {

	field := this.field
	this.Select("count(1) as cnt")
	result, err := this.QueryRow()
	if err != nil {
		return 0, err
	}
	this.Select(field)
	cnt, _ := strconv.Atoi(result["cnt"])
	return cnt, err
}

func (this *DBCommand) QueryList() ([]map[string]interface{}, error) {

	var result []map[string]interface{}
	rows, err := this.pdo.Query(this.Sql(), this.args...)
	if err != nil {
		return result, err
	}
	columns, err := rows.Columns()
	values := make([]string, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}

	for rows.Next() {
		_ = rows.Scan(scans...)
		each := make(map[string]interface{})
		for i, col := range values {
			each[columns[i]] = col
		}
		result = append(result, each)
	}
	rows.Close()
	return result, err
}

func (this *DBCommand) Update(fields map[string]interface{}, condition string, args ...interface{}) (int64, error) {

	var field = ""
	var params []interface{}
	var index = 0

	params = make([]interface{}, len(fields)+len(args))
	for key, val := range fields {
		field += key + "=?,"
		params[index] = val
		index++
	}

	field = strings.Trim(field, ",")
	for _, arg := range args {
		params[index] = arg
		index++
	}

	sql := "update " + this.table + " set " + field
	if len(condition) > 0 {
		sql += " where " + condition
	}

	result, err := this.pdo.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *DBCommand) Insert(record map[string]interface{}) (int64, error) {

	var field = ""
	var params = make([]interface{}, len(record))
	var index = 0
	for key, val := range record {
		field += key + "=?,"
		params[index] = val
		index++
	}
	field = strings.Trim(field, ",")
	sql := "insert into " + this.table + " set " + field
	result, err := this.pdo.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (this *DBCommand) InsertRows() {

}

func (this *DBCommand) Execute() {

}

func (this *DBCommand) Begin() {

}

func (this *DBCommand) Commit() {

}

func (this *DBCommand) Rollback() {

}

func (this *DBCommand) Find(model interfaces.DBModel, pkValue interface{}) error {
	return this.Where(model.PrimaryKey()+"=?", pkValue).QueryBind(model)
}

func (this *DBCommand) Save(pkID string, fields map[string]interface{}) error {

	var err error
	var pkValue, pkIndex string
	for key, val := range fields {
		if strings.ToLower(key) == strings.ToLower(pkID) {
			pkIndex = key
			pkValue = conv.String(val)
			break
		}
	}
	if pkValue == "" || pkValue == "0" {
		delete(fields, pkIndex)
		_, err = this.Insert(fields)
	} else {
		_, err = this.Update(fields, pkID+"=?", pkValue)
	}
	return err
}
