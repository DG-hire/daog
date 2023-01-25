package main

import (
	"encoding/json"
	"fmt"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/example/entities"
	dbtime "github.com/rolandhe/daog/time"
	txrequest "github.com/rolandhe/daog/tx"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

var datasource daog.Datasource

func init() {
	conf := &daog.DbConf{
		DbUrl:  "root:12345678@tcp(localhost:3306)/daog?parseTime=true&timeout=1s&readTimeout=2s&writeTimeout=2s",
		LogSQL: true,
	}
	var err error
	datasource, err = daog.NewDatasource(conf)
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {
	defer datasource.Shutdown()

	//create()
	//query()
	queryByIds()
	//queryByMatcher()
	//update()

	//deleteById()
}

func query() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.GetById(1, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))
	fmt.Println(g)
}

func deleteById() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.DeleteById(2, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("delete", g)
}

func queryByIds() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	gs, err := daog.GetByIds([]int64{1, 2}, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIds", string(j))
	fmt.Println(gs)
}

func queryByMatcher() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer func() {
		tc.Complete(err)
	}()
	matcher := daog.NewMatcher().Like(entities.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(entities.GroupInfoFields.Id, 4)
	gs, err := daog.QueryListMatcher(matcher, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcher", string(j))
	fmt.Println(gs)
}

func create() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	amount, err := decimal.NewFromString("128.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &entities.GroupInfo{
		Name:        "roland",
		MainData:    `{"a":102}`,
		CreateAt:    dbtime.NormalDatetime(time.Now()),
		TotalAmount: amount,
	}
	affect, err := daog.Insert(t, entities.GroupInfoMeta, tc)
	fmt.Println(affect, t.Id, err)

	t.Name = "roland he"
	af, err := daog.Update(t, entities.GroupInfoMeta, tc)
	fmt.Println(af, err)
}

func update() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.GetById(4, entities.GroupInfoMeta, tc)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))

	g.Name = "Eric"
	af, err := daog.Update(g, entities.GroupInfoMeta, tc)
	fmt.Println(af, err)

}
