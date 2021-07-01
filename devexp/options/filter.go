package options

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter interface{}

// index for field has to be created
var dateFieldRegex = regexp.MustCompile(`^([^.]+)\.(year|quarter|month|dayofweek|day)$`)

const idField = "_id"

func (m *LoadOptions) ParseFilter() bson.M {
	return parseFilter(m.Filter)
}

func parseFilter(filter interface{}) bson.M {
	if filter == nil {
		return bson.M{}
	}

	switch v := filter.(type) {
	case []interface{}:
		return parseFilterList(v)
	case string:
		return bson.M{parseField(v): bson.M{"$eq": true}}
	default:
		panic(fmt.Errorf("invalid filter value: %v", v))
	}
}

func parseFilterList(fl []interface{}) bson.M {
	sz := len(fl)

	if sz < 1 {
		return nil
	}

	if sz == 1 {
		return parseFilter(fl[0])
	}

	// unary operator "!"
	if sz == 2 {
		if fl[0] != "!" {
			panic(fmt.Errorf("invalid unary operator: %v", fl[0]))
		}
		return bson.M{"$not": parseFilter(fl[1])}
	}

	if sz == 3 {
		if _, ok := fl[0].([]interface{}); ok {
			return parseFilterListChain(fl)
		}
		field := parseField(fl[0])
		operand := fl[2]
		if field == idField && operand != nil {
			var err error
			if operand, err = primitive.ObjectIDFromHex(operand.(string)); err != nil {
				panic(fmt.Errorf("invalid object id provided: %v", operand))
			}
		}
		return bson.M{field: parseExpression(fl[1], operand)}
	}

	if sz%2 == 0 {
		panic(fmt.Errorf("chain of binary operators expected, provided even number of elements: %v", sz))
	}

	return parseFilterListChain(fl)
}

func parseFilterListChain(fl []interface{}) bson.M {
	opds := []bson.M{}
	expOperator := parseChainOperator(fl[1])
	for i, elem := range fl {
		if i%2 == 1 {
			op := parseChainOperator(elem)
			if op != expOperator {
				panic(fmt.Errorf("invalid operator in the chain: %v", op))
			}
		} else {
			opds = append(opds, parseFilter(elem))
		}
	}

	return bson.M{expOperator: opds}
}

func parseChainOperator(val interface{}) string {
	op, ok := val.(string)
	if !ok || (op != "and" && op != "or") {
		panic(fmt.Errorf("invalid logical operator provided: %v", op))
	}

	return "$" + op
}

func parseField(val interface{}) string {
	switch v := val.(type) {
	case string:
		if v == "id" {
			return idField
		}
		m := dateFieldRegex.FindStringSubmatch(v)
		if m == nil {
			return v
		}
		return fmt.Sprintf("__%s_%s", m[0], m[1])
	default:
		panic(fmt.Errorf("invalid field value: %v", v))
	}
}

func parseExpression(operator, operand interface{}) bson.M {
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

func parseOperand(op string, val interface{}) bson.M {
	switch v := val.(type) {
	case nil:
	case int:
	case float64:
	case string:
	case primitive.ObjectID:
	default:
		panic(fmt.Errorf("invalid operand value: %v", v))
	}

	return bson.M{op: val}
}

func parseRegex(format string, val interface{}) bson.M {
	v, ok := val.(string)
	if !ok {
		panic(fmt.Errorf("invalid regex value: %v", v))
	}

	return bson.M{"$regex": fmt.Sprintf(format, regexp.QuoteMeta(v)), "$options": "i"}
}
