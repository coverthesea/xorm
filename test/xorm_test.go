package xorm

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/xormplus/xorm"

	_ "github.com/lib/pq"
)

type Article struct {
	Id             int       `xorm:"not null pk autoincr unique INTEGER"`
	Content        string    `xorm:"not null TEXT"`
	Title          string    `xorm:"not null VARCHAR(255)"`
	Categorysubid  int       `xorm:"not null INTEGER"`
	Remark         string    `xorm:"not null VARCHAR(2555)"`
	Userid         int       `xorm:"not null INTEGER"`
	Viewcount      int       `xorm:"not null default 0 INTEGER"`
	Replycount     int       `xorm:"not null default 0 INTEGER"`
	Tags           string    `xorm:"not null VARCHAR(300)"`
	Createdatetime JSONTime  `xorm:"not null default 'now()' DATETIME"`
	Isdraft        int       `xorm:"SMALLINT"`
	Lastupdatetime time.Time `xorm:"not null default 'now()' DATETIME"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006/01/08 15:04:05"))
	return []byte(stamp), nil
}

var db *xorm.Engine

func Test_InitDB(t *testing.T) {
	var err error
	db, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/mblog?sslmode=disable")

	if err != nil {
		t.Fatal(err)
	}

	err = db.InitSqlMap()
	if err != nil {
		t.Fatal(err)
	}
	err = db.InitSqlTemplate()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Get_Struct(t *testing.T) {
	var article Article
	has, err := db.Id(3).Get(&article)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Get_Struct]->rows: not exist\n")
	}

	t.Log("[Test_Get_Struct]->rows:\n" , article)
}

func Test_GetFirst_Json(t *testing.T) {
	var article Article
	has, rows, err := db.Id(2).GetFirst(&article).Json()
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_Json]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_Json]->rows:\n" + rows)
}

func Test_GetFirst_Xml(t *testing.T) {
	var article Article
	has, rows, err := db.Where("userid =?", 3).GetFirst(&article).Xml()
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_Xml]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_Xml]->rows:\n" + rows)
}

func Test_GetFirst_XmlIndent(t *testing.T) {
	var article Article
	has, rows, err := db.Where("userid =?", 3).GetFirst(&article).XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_XmlIndent]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_XmlIndent]->rows:\n" + rows)
}

func Test_FindAll_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_Json]->rows:\n" + rows)
}

func Test_FindAll_ID(t *testing.T) {
	rows := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).Query()
	if rows.Error != nil {
		t.Fatal(rows.Error)
	}

	t.Log("[Test_FindAll_Json]->rows[0][\"id\"]:\n", rows.Result[0]["id"])
	t.Log("[Test_FindAll_Json]->reflect.TypeOf(rows.Result[0][\"id\"]):\n", reflect.TypeOf(rows.Result[0]["id"]))
	t.Log("[Test_FindAll_Json]->rows[0][\"title\"]:\n", rows.Result[0]["title"])
	t.Log("[Test_FindAll_Json]->reflect.TypeOf(rows.Result[0][\"title\"]):\n", reflect.TypeOf(rows.Result[0]["title"]))
	t.Log("[Test_FindAll_Json]->rows[0][\"createdatetime\"]:\n", rows.Result[0]["createdatetime"])
	t.Log("[Test_FindAll_Json]->reflect.TypeOf(rows.Result[0][\"createdatetime\"]):\n", reflect.TypeOf(rows.Result[0]["createdatetime"]))
}

func Test_FindAll_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_Xml]->rows:\n" + rows)
}

func Test_FindAll_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).QueryWithDateFormat("20060102").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_Json]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).QueryWithDateFormat("20060102").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id in (?,?)", 2, 5).QueryWithDateFormat("20060102").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 4, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).QueryByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).QueryByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).QueryByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 5, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).QueryByParamMapWithDateFormat("2006/01/02").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMapWithDateFormat("2006-01-02 15:04").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMapWithDateFormat("2006-01-02 15:04").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMapWithDateFormat("2006-01-02 15:04").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMapWithDateFormat("01/02/2006").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMapWithDateFormat("01/02/2006").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMapWithDateFormat("01/02/2006").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_Find_Structs_Json(t *testing.T) {
	articles := make([]Article, 0)
	json,err := db.Where("id=?", 6).Find(&articles).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_Json]->rows:\n" + json)
}

func Test_Find_Structs_Xml(t *testing.T) {
	articles := make([]Article, 0)
	xml,err := db.Where("id=?", 6).Find(&articles).Xml()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_Xml]->rows:\n" + xml)
}

func Test_Find_Structs_XmlIndent(t *testing.T) {
	articles := make([]Article, 0)
	xml,err := db.Where("id=?", 6).Find(&articles).XmlIndent("","  ","Article")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_XmlIndent]->rows:\n" + xml)
}