package form

import (
	"github.com/phper-go/frame/db"
)

type Search struct {
	PageSize  int
	PageIndex int
}

func (this *Search) Title() []string {
	return []string{}
}

func (this *Search) DBCmd() *db.DBCommand {
	return &db.DBCommand{}
}

func (this *Search) DBPageCmd(dbCmd *db.DBCommand) *db.DBCommand {

	cmd := dbCmd.Clone()
	if this.PageIndex < 1 {
		this.PageIndex = 1
	}
	if this.PageSize < 1 {
		this.PageSize = 20
	}
	cmd.Offset((this.PageIndex - 1) * this.PageSize)
	cmd.Limit(this.PageSize)
	return cmd
}

func (this *Search) List(dbCmd db.DBCommand) []map[string]interface{} {

	return []map[string]interface{}{}
}
