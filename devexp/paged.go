package devexp

import (
	"strings"

	"github.com/drahoslavzan/srvutils/devexp/options"
	"github.com/drahoslavzan/srvutils/log"

	pagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

type Paged struct {
	Col    *mongo.Collection
	Opts   *options.LoadOptions
	Filter bson.M
	Logger log.Logger
	Locale *string
}

func (m *Paged) Find(decode interface{}) (total int64) {
	page := m.pageNo()
	filter := m.makeFilter()
	pg := pagination.New(m.Col)

	m.Logger.Debugw("find",
		"page", page,
		"take", m.Opts.Take,
		"sort", m.Opts.Sort,
		"filter", filter,
	)

	m.fillSort(pg, m.Opts)

	paged, err := pg.Limit(m.Opts.Take).Page(page).Filter(filter).Decode(decode).Find()
	if err != nil {
		m.Logger.Panic(err)
	}

	total = paged.Pagination.Total

	return
}

func (m *Paged) GroupBy(group []*options.Group, decode interface{}) (data []bson.Raw, total int64) {
	page := m.pageNo()
	filter := m.makeFilter()
	pg := pagination.New(m.Col)

	pipeline := []interface{}{bson.M{
		"$match": filter,
	}}

	projID := "_id"
	grpProj := bson.M{}
	var grpSel bson.D
	var grpSelInd bson.D
	for _, g := range group {
		sel := getDbKeyName(g.Selector)
		grpSel = append(grpSel, primitive.E{sel, "$" + g.Selector})
		grpSelInd = append(grpSelInd, primitive.E{sel, "$_id." + sel})
		grpProj[projID] = 0
		projID = "items." + projID
	}

	grpLast := len(grpSel) - 1

	sortBy := bson.M{}
	for _, s := range m.Opts.Sort {
		sortBy[s.GetField(m.Opts)] = s.GetOrder()
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

	m.Logger.Debugw("find",
		"page", page,
		"take", m.Opts.Take,
		"sort", m.Opts.Sort,
		"filter", filter,
		"group", group,
		"pipeline", pipeline,
	)

	ag, err := pg.Limit(m.Opts.Take).Page(page).Aggregate(pipeline...)
	if err != nil {
		m.Logger.Panic(err)
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

	if m.Opts.Search != nil {
		filter["$text"] = bson.M{
			"$search": *m.Opts.Search,
		}
	}

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

func (m *Paged) fillSort(pg pagination.PagingQuery, opts *options.LoadOptions) {
	if len(opts.Sort) < 1 {
		pg.Sort("_id", 1)
		return
	}

	locale := "en"
	if m.Locale != nil {
		locale = *m.Locale
	}

	pg.SetCollation(&mopts.Collation{Locale: locale})
	for _, s := range opts.Sort {
		pg.Sort(s.GetField(opts), s.GetOrder())
	}
}

func getDbKeyName(selector string) string {
	// MongoDB cannot use '.' as a field name in aggregation
	return strings.ReplaceAll(selector, ".", "_")
}
