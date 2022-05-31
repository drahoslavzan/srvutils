package options

import (
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter any

// index for a field has to be created
var dateFieldRegex = regexp.MustCompile(`^([^.]+)\.(year|quarter|month|dayofweek|day)$`)

func (m *LoadOptions) ParseFilter() bson.M {
	if m.Field == nil {
		m.Field = map[string]Field{}
	}
	return m.parseFilter(m.Filter)
}

func (m *LoadOptions) parseFilter(filter Filter) bson.M {
	if filter == nil {
		return bson.M{}
	}

	switch v := filter.(type) {
	case []any:
		return m.parseFilterList(v)
	case string:
		return bson.M{m.parseField(v).Name: bson.M{"$eq": true}}
	default:
		panic(fmt.Errorf("invalid filter: %v", v))
	}
}

func (m *LoadOptions) parseFilterList(fl []any) bson.M {
	sz := len(fl)

	if sz < 1 {
		return nil
	}

	if sz == 1 {
		return m.parseFilter(fl[0])
	}

	// unary operator "!"
	if sz == 2 {
		if fl[0] != "!" {
			panic(fmt.Errorf("invalid unary operator: %v", fl[0]))
		}
		return bson.M{"$nor": []bson.M{m.parseFilter(fl[1])}}
	}

	if sz == 3 {
		if _, ok := fl[0].([]any); ok {
			return m.parseFilterListChain(fl)
		}
		field := m.parseField(fl[0])
		operand := fl[2]
		if field.Serialize != nil {
			operand = field.Serialize(operand)
		}
		return bson.M{field.Name: parseExpression(fl[1], operand)}
	}

	if sz%2 == 0 {
		panic(fmt.Errorf("chain of binary operators expected, provided even number of elements: %v", sz))
	}

	return m.parseFilterListChain(fl)
}

func (m *LoadOptions) parseFilterListChain(fl []any) bson.M {
	opds := []bson.M{}
	expOperator := parseChainOperator(fl[1])
	for i, elem := range fl {
		if i%2 == 1 {
			op := parseChainOperator(elem)
			if op != expOperator {
				panic(fmt.Errorf("invalid operator in the chain: %v", op))
			}
		} else {
			opds = append(opds, m.parseFilter(elem))
		}
	}

	return bson.M{expOperator: opds}
}

func (m *LoadOptions) parseField(val any) *Field {
	switch v := val.(type) {
	case string:
		f, ok := m.Field[v]
		if !ok {
			f = Field{
				Name: v,
			}
		} else {
			if len(f.Name) < 1 {
				f.Name = v
			}
		}
		m := dateFieldRegex.FindStringSubmatch(f.Name)
		if m != nil {
			f.Name = fmt.Sprintf("__%s_%s", m[0], m[1])
		}
		return &f
	default:
		panic(fmt.Errorf("invalid field: %v", v))
	}
}

func parseChainOperator(val any) string {
	op, ok := val.(string)
	if !ok || (op != "and" && op != "or") {
		panic(fmt.Errorf("invalid logical operator provided: %v", op))
	}

	return "$" + op
}

func parseExpression(operator, operand any) bson.M {
	op, ok := operator.(string)
	if !ok {
		panic(fmt.Errorf("invalid operator provided: %v", op))
	}

	switch op {
	case "=":
		return parseOperand("$eq", operand)
	case "<>":
		return parseOperand("$ne", operand)
	case ">":
		return parseOperand("$gt", operand)
	case ">=":
		return parseOperand("$gte", operand)
	case "<":
		return parseOperand("$lt", operand)
	case "<=":
		return parseOperand("$lte", operand)
	case "startswith":
		return parseRegex("^%s", operand)
	case "endswith":
		return parseRegex("%s$", operand)
	case "contains":
		return parseRegex("%s", operand)
	case "notcontains":
		return parseRegex("^((?!%s).)*$", operand)
	default:
		panic(fmt.Errorf("invalid operator provided: %v", op))
	}
}

func parseOperand(op string, val any) bson.M {
	switch v := val.(type) {
	case nil:
	case bool:
	case int:
	case float64:
	case string:
	case time.Time:
	case primitive.ObjectID:
	default:
		panic(fmt.Errorf("invalid operand: %v", v))
	}

	return bson.M{op: val}
}

func parseRegex(format string, val any) bson.M {
	v, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("invalid regex: %v", v))
	}

	return bson.M{"$regex": fmt.Sprintf(format, regexp.QuoteMeta(v)), "$options": "i"}
}
