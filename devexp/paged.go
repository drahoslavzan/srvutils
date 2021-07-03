package devexp

import (
	"encoding/json"

	"github.com/drahoslavzan/srvutils/devexp/options"
	"github.com/drahoslavzan/srvutils/env"
	"github.com/drahoslavzan/srvutils/log"

	pagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Paged struct {
	LogFunc string
	Col     *mongo.Collection
	Opts    *options.LoadOptions
	Filter  bson.M
}

func (m *Paged) Find(decode interface{}) (total int64) {
	logger := log.GetLogger(log.LoggerOpts{FuncName: m.LogFunc})
	page := m.pageNo()
	filter := m.makeFilter()
	pg := pagination.New(m.Col)

	if env.IsDevelopment() {
		f, _ := json.MarshalIndent(filter, "", "  ")
		logger.Debugf("page: %v, take: %v, sort: %+v", page, m.Opts.Take, m.Opts.Sort)
		logger.Debugf("  - filter: %s", f)
	}

	fillSort(pg, m.Opts.Sort)

	paged, err := pg.Limit(m.Opts.Take).Page(page).Filter(filter).Decode(decode).Find()
	if err != nil {
		panic(err)
	}

	total = paged.Pagination.Total

	return
}

func (m *Paged) GroupBy(group []*options.Group, decode interface{}) (data []bson.Raw, total int64) {
	logger := log.GetLogger(log.LoggerOpts{FuncName: m.LogFunc})
	page := m.pageNo()
	filter := m.makeFilter()
	pg := pagination.New(m.Col)

	if env.IsDevelopment() {
		f, _ := json.MarshalIndent(filter, "", "  ")
		logger.Debugf("page: %v, take: %v, sort: %+v, group: %+v", page, m.Opts.Take, m.Opts.Sort, group)
		logger.Debugf("  - filter: %s", f)
	}

	pipeline := []interface{}{bson.M{
		"$match": filter,
	}}

	projID := "_id"
	grpProj := bson.M{}
	var grpSel bson.D
	var grpSelInd bson.D
	for _, g := range group {
		grpSel = append(grpSel, primitive.E{g.Selector, "$" + g.Selector})
		grpSelInd = append(grpSelInd, primitive.E{g.Selector, "$_id." + g.Selector})
		grpProj[projID] = 0
		projID = "items." + projID
	}

	grpLast := len(grpSel) - 1

	sortBy := bson.M{}
	for _, s := range m.Opts.Sort {
		sortBy[s.GetField()] = s.GetOrder()
	}

	pipeline = append(pipeline, bson.M{"$sort": sortBy})

	pipeline = append(pipeline, bson.M{"$group": bson.M{
		"_id":   grpSel,
		"key":   bson.M{"$first": grpSel[grpLast].Value},
		"items": bson.M{"$push": bson.M{"item": "$$ROOT"}},
	}})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"key": group[grpLast].GetOrder(),
	}})

	for i := grpLast - 1; i >= 0; i-- {
		pipeline = append(pipeline, bson.M{"$group": bson.M{
			"_id":   grpSelInd[:i+1],
			"key":   bson.M{"$first": grpSelInd[i].Value},
			"items": bson.M{"$push": "$$ROOT"},
		}})

		pipeline = append(pipeline, bson.M{"$sort": bson.M{
			"key": group[i].GetOrder(),
		}})
	}

	pipeline = append(pipeline, bson.M{"$project": grpProj})

	ag, err := pg.Limit(m.Opts.Take).Page(page).Aggregate(pipeline...)
	if err != nil {
		panic(err)
	}

	data = ag.Data
	total = ag.Pagination.Total

	return
}

func (m *Paged) pageNo() int64 {
	return m.Opts.Skip/m.Opts.Take + 1
}

func (m *Paged) makeFilter() bson.M {
	filter := m.Opts.ParseFilter()
	if m.Filter != nil {
		for k, v := range m.Filter {
			if pv, ok := filter[k]; ok {
				delete(filter, k)
				filter["$and"] = []interface{}{
					bson.M{k: v},
					bson.M{k: pv},
				}
				continue
			}
			filter[k] = v
		}
	}

	return filter
}

func fillSort(pg pagination.PagingQuery, sorts []options.Sort) {
	if len(sorts) < 1 {
		pg.Sort("_id", 1)
		return
	}

	for _, s := range sorts {
		pg.Sort(s.GetField(), s.GetOrder())
	}
}
