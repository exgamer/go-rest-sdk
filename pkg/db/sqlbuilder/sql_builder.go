package sqlbuilder

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func MakeFoundRowsSql() string {
	return "SELECT FOUND_ROWS()"
}

// MakeCountSelectSql - returns full count select sql string
func MakeCountSelectSql(queryData *QueryData) string {
	return makeSelectSql(queryData, true, false)
}

// MakeSelectSql - returns full select sql string
func MakeSelectSql(queryData *QueryData) string {
	return makeSelectSql(queryData, false, true)
}

// makeSelectSql - returns full select sql string
func makeSelectSql(queryData *QueryData, countSelect bool, addOrder bool) string {
	builder := strings.Builder{}
	builder.WriteString("SELECT ")

	if queryData.CalcRows == true {
		builder.WriteString(" SQL_CALC_FOUND_ROWS ")
	}

	if countSelect == false {
		if len(queryData.SelectArray) > 0 {
			builder.WriteString(strings.Join(queryData.SelectArray, ", "))
		} else {
			builder.WriteString("*")
		}
	} else {
		builder.WriteString("count(*)")
	}

	builder.WriteString(" FROM " + queryData.TableName)

	if len(queryData.TableAliasName) > 0 {
		builder.WriteString(" " + queryData.TableAliasName)
	}

	builder.WriteString(" " + makeJoinSql(queryData))
	builder.WriteString(" " + makeWhereSql(queryData))

	if len(queryData.Group) > 0 && countSelect == false {
		builder.WriteString(" GROUP BY " + queryData.Group)
	}

	if len(queryData.Order) > 0 && addOrder == true {
		builder.WriteString(" ORDER BY " + strings.Join(queryData.Order, ","))
	}

	if queryData.LimitCount > 0 {
		builder.WriteString(" LIMIT " + strconv.Itoa(queryData.LimitCount))
	}

	if queryData.OffsetCount > 0 {
		builder.WriteString(" OFFSET " + strconv.Itoa(queryData.OffsetCount))
	}

	sqlString := builder.String()
	log.Println("MakeSelectSql: " + sqlString)
	fmt.Print(queryData.Params)

	return sqlString
}

// makeJoinSql - returns join part of sql string
func makeJoinSql(queryData *QueryData) string {
	sqlString := ""

	if len(queryData.JoinCondition) == 0 {
		return sqlString
	}

	joinArray := make([]string, len(queryData.JoinCondition))

	for i := 0; i < len(queryData.JoinCondition); i++ {
		joinArray[i] = queryData.JoinCondition[i].Type + " " + queryData.JoinCondition[i].Table + " ON " + queryData.JoinCondition[i].On + " " + queryData.JoinCondition[i].Where
	}

	return strings.Join(joinArray, " ")
}

// makeWhereSql - returns where part of sql string
func makeWhereSql(queryData *QueryData) string {
	sqlString := ""

	if len(queryData.WhereCondition) == 0 {
		return sqlString
	}

	whereArray := make([]string, len(queryData.WhereCondition))

	for i := 0; i < len(queryData.WhereCondition); i++ {
		if i == 0 {
			fmt.Printf("%v", whereArray)

			whereArray[i] = queryData.WhereCondition[i].Condition

			continue
		}

		whereArray[i] = queryData.WhereCondition[i].Operator + " " + queryData.WhereCondition[i].Condition
	}

	return " WHERE " + strings.Join(whereArray, " ")
}

// MakeInsertSql - returns insert sql string
func MakeInsertSql(queryData *QueryData) string {
	cols := make([]string, len(queryData.Data))
	vals := make([]string, len(queryData.Data))

	i := 0

	for key, element := range queryData.Data {
		cols[i] = key

		if reflect.TypeOf(element).Name() == "string" {
			vals[i] = "'" + fmt.Sprint(element) + "'"
		} else {
			vals[i] = fmt.Sprint(element)
		}

		i++
	}

	sqlString := "INSERT INTO " + queryData.TableName + " (" + strings.Join(cols, ",") + ") VALUES" + " (" + strings.Join(vals, ",") + ")"
	log.Println("MakeInsertSql: " + sqlString)

	return sqlString
}

// MakeUpdateSql - returns update sql string
func MakeUpdateSql(queryData *QueryData) string {
	cols := make([]string, len(queryData.Data))

	i := 0

	for key, element := range queryData.Data {
		if reflect.TypeOf(element).Name() == "string" {
			cols[i] = key + "='" + fmt.Sprint(element) + "'"
		} else {
			cols[i] = key + "=" + fmt.Sprint(element)
		}

		i++
	}

	builder := strings.Builder{}
	builder.WriteString("UPDATE " + queryData.TableName + " SET " + strings.Join(cols, ",") + " ")
	builder.WriteString(makeWhereSql(queryData))

	sqlString := builder.String()

	log.Println("MakeUpdateSql: " + sqlString)

	return sqlString
}

// MakeDeleteSql - returns delete sql string
func MakeDeleteSql(queryData *QueryData) string {
	builder := strings.Builder{}
	builder.WriteString("DELETE FROM " + queryData.TableName + " ")
	builder.WriteString(makeWhereSql(queryData))

	sqlString := builder.String()

	log.Println("MakeDeleteSql: " + sqlString)

	return sqlString
}
