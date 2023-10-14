package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/entity"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/form"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	paginator "github.com/exgamer/go-rest-sdk/pkg/pager"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		timeout: 30,
		DbType:  "mysql",
	}
}

func NewMysqlQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		timeout: 30,
		DbType:  "postgres",
	}
}

func NewPostgresQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		timeout: 30,
		DbType:  "postgres",
	}
}

//QueryBuilder - query builder
type QueryBuilder struct {
	Db               *sql.DB
	TableName        string
	TableAliasName   string
	Data             map[string]any
	WhereCondition   []WhereCondition
	CalcRows         bool
	SelectArray      []string
	JoinCondition    []JoinCondition
	Params           []any
	LimitCount       int
	OffsetCount      int
	Order            []string
	Group            string
	timeout          time.Duration
	ValuePlaceholder string
	DbType           string
}

type WhereCondition struct {
	Operator  string
	Condition string
}

type JoinCondition struct {
	Type  string
	Table string
	On    string
}

//SetDb - set sql.DB connection (for make query)
func (queryBuilder *QueryBuilder) SetDb(db *sql.DB) *QueryBuilder {
	queryBuilder.Db = db

	return queryBuilder
}

func (queryBuilder *QueryBuilder) SetDbType(dbType string) *QueryBuilder {
	queryBuilder.DbType = dbType

	return queryBuilder
}

func (queryBuilder *QueryBuilder) getPlaceholder(i int) string {
	if queryBuilder.DbType == "postgres" {
		index := len(queryBuilder.Params) + 1 + i

		return "$" + strconv.Itoa(index)
	}

	return "?"
}

//SetEntity - set entity to query
func (queryBuilder *QueryBuilder) SetEntity(entity entity.Entity) *QueryBuilder {
	queryBuilder.TableName = entity.GetTable()

	return queryBuilder
}

//SetData - set data to query
func (queryBuilder *QueryBuilder) SetData(data map[string]any) *QueryBuilder {
	queryBuilder.Data = data
	fmt.Println(data)

	return queryBuilder
}

//SetFormData - set data to query (for insert/update)
func (queryBuilder *QueryBuilder) SetFormData(form form.Form) *QueryBuilder {
	queryBuilder.Data = form.AsMap()

	return queryBuilder
}

func (queryBuilder *QueryBuilder) SetQueryTimeout(timeout time.Duration) *QueryBuilder {
	queryBuilder.timeout = timeout

	return queryBuilder
}

//Table - set table name
func (queryBuilder *QueryBuilder) Table(table string) *QueryBuilder {
	queryBuilder.TableName = table

	return queryBuilder
}

//Select - set select fields
func (queryBuilder *QueryBuilder) Select(fields ...string) *QueryBuilder {
	queryBuilder.SelectArray = fields

	return queryBuilder
}

//TableAlias - set from table alias
func (queryBuilder *QueryBuilder) TableAlias(tableAlias string) *QueryBuilder {
	queryBuilder.TableAliasName = tableAlias

	return queryBuilder
}

func (queryBuilder *QueryBuilder) CalculateRows() *QueryBuilder {
	queryBuilder.CalcRows = true

	return queryBuilder
}

//SetLimitOffsetByPage - set limit and offset to query by page
func (queryBuilder *QueryBuilder) SetLimitOffsetByPage(page int, perPage int) *QueryBuilder {
	var limit = 30
	var offset = 0

	if perPage > 0 {
		limit = perPage
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	queryBuilder.LimitCount = limit
	queryBuilder.OffsetCount = offset

	return queryBuilder
}

//Limit - set limit to query
func (queryBuilder *QueryBuilder) Limit(limit int) *QueryBuilder {
	queryBuilder.LimitCount = limit

	return queryBuilder
}

//Offset - set offset to query
func (queryBuilder *QueryBuilder) Offset(offset int) *QueryBuilder {
	queryBuilder.OffsetCount = offset

	return queryBuilder
}

//OrderByAsc - add order by asc
func (queryBuilder *QueryBuilder) OrderByAsc(orderBy string) *QueryBuilder {

	return queryBuilder.orderBy(orderBy, "ASC")
}

//OrderByDesc - add order by desc
func (queryBuilder *QueryBuilder) OrderByDesc(orderBy string) *QueryBuilder {

	return queryBuilder.orderBy(orderBy, "DESC")
}

//orderBy - add order by
func (queryBuilder *QueryBuilder) orderBy(orderBy string, sort string) *QueryBuilder {
	queryBuilder.Order = append(queryBuilder.Order, orderBy+" "+sort)

	return queryBuilder
}

//GroupBy - add group by
func (queryBuilder *QueryBuilder) GroupBy(groupBy string) *QueryBuilder {
	queryBuilder.Group = groupBy

	return queryBuilder
}

//GetParams - returns query params
func (queryBuilder *QueryBuilder) GetParams() []any {

	return queryBuilder.Params
}

func (queryBuilder *QueryBuilder) MakeFoundRowsSql() string {
	return "SELECT FOUND_ROWS()"
}

//MakeCountSelectSql - returns full count select sql string
func (queryBuilder *QueryBuilder) MakeCountSelectSql() string {
	return queryBuilder.makeSelectSql(true, false)
}

//MakeSelectSql - returns full select sql string
func (queryBuilder *QueryBuilder) MakeSelectSql() string {
	return queryBuilder.makeSelectSql(false, true)
}

//makeSelectSql - returns full select sql string
func (queryBuilder *QueryBuilder) makeSelectSql(countSelect bool, addOrder bool) string {
	sqlString := "SELECT "

	if queryBuilder.CalcRows == true {
		sqlString += " SQL_CALC_FOUND_ROWS "
	}

	if countSelect == false {
		if len(queryBuilder.SelectArray) > 0 {
			sqlString += strings.Join(queryBuilder.SelectArray, ", ")
		} else {
			sqlString += "*"
		}
	} else {
		sqlString += "count(*)"
	}

	sqlString += " FROM " + queryBuilder.TableName

	if len(queryBuilder.TableAliasName) > 0 {
		sqlString += " " + queryBuilder.TableAliasName
	}

	sqlString += " " + queryBuilder.makeJoinSql()
	sqlString += " " + queryBuilder.makeWhereSql()

	if len(queryBuilder.Group) > 0 {
		sqlString += " GROUP BY " + queryBuilder.Group
	}

	if len(queryBuilder.Order) > 0 && addOrder == true {
		sqlString += " ORDER BY " + strings.Join(queryBuilder.Order, ",")
	}

	if queryBuilder.LimitCount > 0 {
		sqlString += " LIMIT " + strconv.Itoa(queryBuilder.LimitCount)
	}

	if queryBuilder.OffsetCount > 0 {
		sqlString += " OFFSET " + strconv.Itoa(queryBuilder.OffsetCount)
	}

	log.Println("MakeSelectSql: " + sqlString)
	fmt.Print(queryBuilder.Params)

	return sqlString
}

//Join - add join condition
func (queryBuilder *QueryBuilder) Join(table string, first string, second string, params ...string) *QueryBuilder {
	queryBuilder.join(table, first+"="+second, "JOIN", params...)

	return queryBuilder
}

//OuterJoin - add outer join condition
func (queryBuilder *QueryBuilder) OuterJoin(table string, first string, second string, params ...string) *QueryBuilder {
	queryBuilder.join(table, first+"="+second, "OUTER JOIN", params...)

	return queryBuilder
}

//InnerJoin - add inner join condition
func (queryBuilder *QueryBuilder) InnerJoin(table string, first string, second string, params ...string) *QueryBuilder {
	queryBuilder.join(table, first+"="+second, "INNER JOIN", params...)

	return queryBuilder
}

//RightJoin - add right join condition
func (queryBuilder *QueryBuilder) RightJoin(table string, first string, second string, params ...string) *QueryBuilder {
	queryBuilder.join(table, first+"="+second, "RIGHT JOIN", params...)

	return queryBuilder
}

//LeftJoin - add left join condition
func (queryBuilder *QueryBuilder) LeftJoin(table string, first string, second string, params ...string) *QueryBuilder {
	queryBuilder.join(table, first+"="+second, "LEFT JOIN", params...)

	return queryBuilder
}

//join - add join condition by string condition
func (queryBuilder *QueryBuilder) join(table string, on string, joinType string, params ...string) {
	join := JoinCondition{Table: table, On: on, Type: joinType}
	queryBuilder.JoinCondition = append(queryBuilder.JoinCondition, join)

	for _, p := range params {
		queryBuilder.Params = append(queryBuilder.Params, p)
	}
}

//makeJoinSql - returns join part of sql string
func (queryBuilder *QueryBuilder) makeJoinSql() string {
	sqlString := ""

	if len(queryBuilder.JoinCondition) == 0 {
		return sqlString
	}

	joinArray := make([]string, len(queryBuilder.JoinCondition))

	for i := 0; i < len(queryBuilder.JoinCondition); i++ {
		joinArray[i] = queryBuilder.JoinCondition[i].Type + " " + queryBuilder.JoinCondition[i].Table + " ON " + queryBuilder.JoinCondition[i].On
	}

	return strings.Join(joinArray, " ")
}

//OrWhere - adds or where condition
func (queryBuilder *QueryBuilder) OrWhere(field string, value string) *QueryBuilder {
	queryBuilder.OrWhereByCondition(field+"="+queryBuilder.getPlaceholder(0), value)

	return queryBuilder
}

//OrWhereByCondition - adds or where condition  by string
func (queryBuilder *QueryBuilder) OrWhereByCondition(condition string, params ...string) *QueryBuilder {
	queryBuilder.WhereByCondition(condition, "OR", params...)

	return queryBuilder
}

//AndWhere - adds and where condition
func (queryBuilder *QueryBuilder) AndWhere(field string, value string) *QueryBuilder {
	queryBuilder.AndWhereByCondition(field+"="+queryBuilder.getPlaceholder(0), value)

	return queryBuilder
}

//AndWhereByCondition - adds where condition by string
func (queryBuilder *QueryBuilder) AndWhereByCondition(condition string, params ...string) *QueryBuilder {
	queryBuilder.WhereByCondition(condition, "AND", params...)

	return queryBuilder
}

//AndWhereIn - adds and in condition
func (queryBuilder *QueryBuilder) AndWhereIn(field string, params []string) *QueryBuilder {
	return queryBuilder.WhereIn(field, params, "AND")
}

//WhereIn - adds in condition
func (queryBuilder *QueryBuilder) WhereIn(field string, params []string, operator string) *QueryBuilder {
	sParams := make([]string, len(params))

	for i := 0; i < len(params); i++ {
		sParams[i] = queryBuilder.getPlaceholder(i)
		queryBuilder.Params = append(queryBuilder.Params, params[i])
	}

	condition := field + " IN (" + strings.Join(sParams, ",") + ")"
	where := WhereCondition{Condition: condition, Operator: operator}
	queryBuilder.WhereCondition = append(queryBuilder.WhereCondition, where)

	return queryBuilder
}

//WhereByCondition - adds string condition
func (queryBuilder *QueryBuilder) WhereByCondition(condition string, operator string, params ...string) *QueryBuilder {
	where := WhereCondition{Condition: condition, Operator: operator}
	queryBuilder.WhereCondition = append(queryBuilder.WhereCondition, where)
	fmt.Print(params)

	for _, p := range params {
		queryBuilder.Params = append(queryBuilder.Params, p)
	}

	return queryBuilder
}

//makeWhereSql - returns where part of sql string
func (queryBuilder *QueryBuilder) makeWhereSql() string {
	sqlString := ""

	if len(queryBuilder.WhereCondition) == 0 {
		return sqlString
	}

	whereArray := make([]string, len(queryBuilder.WhereCondition))

	for i := 0; i < len(queryBuilder.WhereCondition); i++ {
		if i == 0 {
			fmt.Printf("%v", whereArray)

			whereArray[i] = queryBuilder.WhereCondition[i].Condition

			continue
		}

		whereArray[i] = queryBuilder.WhereCondition[i].Operator + " " + queryBuilder.WhereCondition[i].Condition
	}

	return " WHERE " + strings.Join(whereArray, " ")
}

//MakeInsertSql - returns insert sql string
func (queryBuilder *QueryBuilder) MakeInsertSql() string {
	cols := make([]string, len(queryBuilder.Data))
	vals := make([]string, len(queryBuilder.Data))

	i := 0

	for key, element := range queryBuilder.Data {
		cols[i] = key

		if reflect.TypeOf(element).Name() == "string" {
			vals[i] = "'" + fmt.Sprint(element) + "'"
		} else {
			vals[i] = fmt.Sprint(element)
		}

		i++
	}

	sqlString := "INSERT INTO " + queryBuilder.TableName + " (" + strings.Join(cols, ",") + ") VALUES" + " (" + strings.Join(vals, ",") + ")"
	log.Println("MakeInsertSql: " + sqlString)

	return sqlString
}

//MakeUpdateSql - returns update sql string
func (queryBuilder *QueryBuilder) MakeUpdateSql() string {
	cols := make([]string, len(queryBuilder.Data))

	i := 0

	for key, element := range queryBuilder.Data {
		if reflect.TypeOf(element).Name() == "string" {
			cols[i] = key + "='" + fmt.Sprint(element) + "'"
		} else {
			cols[i] = key + "=" + fmt.Sprint(element)
		}

		i++
	}

	sqlString := "UPDATE " + queryBuilder.TableName + " SET " + strings.Join(cols, ",") + " "
	sqlString += queryBuilder.makeWhereSql()
	log.Println("MakeUpdateSql: " + sqlString)

	return sqlString
}

//MakeDeleteSql - returns delete sql string
func (queryBuilder *QueryBuilder) MakeDeleteSql() string {
	sqlString := "DELETE FROM " + queryBuilder.TableName + " "
	sqlString += queryBuilder.makeWhereSql()
	log.Println("MakeDeleteSql: " + sqlString)

	return sqlString
}

//All - returns all query result
func (queryBuilder *QueryBuilder) All() (*[]map[string]interface{}, *exception.AppException) {
	res, _, appException := queryBuilder.paginate(false, 0, 1000)

	return res, appException
}

//One - returns one query result
func (queryBuilder *QueryBuilder) One() (*[]map[string]interface{}, *exception.AppException) {
	res, _, appException := queryBuilder.paginate(false, 0, 1)

	return res, appException
}

//Paginate - returns paginated query result
func (queryBuilder *QueryBuilder) Paginate(page int, perPage int) (*[]map[string]interface{}, *paginator.Pager, *exception.AppException) {
	return queryBuilder.paginate(true, page, perPage)
}

func (queryBuilder *QueryBuilder) paginate(paginate bool, page int, perPage int) (*[]map[string]interface{}, *paginator.Pager, *exception.AppException) {
	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	res := make([]map[string]interface{}, 0)
	pager := paginator.Pager{
		Page: page,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	if paginate == true {
		err := queryBuilder.Db.QueryRowContext(ctx, queryBuilder.MakeCountSelectSql(), queryBuilder.GetParams()...).Scan(&pager.TotalItems)
		if err != nil {
			logger.LogError(err)

			return nil, &pager, exception.NewAppException(http.StatusInternalServerError, err, nil)
		}

		if pager.TotalItems == 0 {
			return &res, &pager, nil
		}

	}

	queryBuilder.SetLimitOffsetByPage(page, perPage)
	// Execute the query
	rows, err := queryBuilder.Db.QueryContext(ctx, queryBuilder.MakeSelectSql(), queryBuilder.GetParams()...)

	if err != nil {
		logger.LogError(err)

		return nil, &pager, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	defer rows.Close()
	//  Data columns
	columns, err := rows.Columns()

	if err != nil {
		logger.LogError(err)

		return nil, &pager, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	//  Number of columns
	count := len(columns)
	//  The value of each column of a piece of data （ You need to specify the number of columns with length , To get the address ）
	values := make([]interface{}, count)
	//  The address of the value of each column of a piece of data
	valPointers := make([]interface{}, count)

	for rows.Next() {
		//  Get the address of the value of each column
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}

		//  Get the value of each column , Put it in the corresponding address
		rows.Scan(valPointers...)
		//  A piece of data Map ( Key value pairs for column names and values )
		entry := make(map[string]interface{})

		// Map  assignment
		for i, col := range columns {
			var v interface{}

			//  Value copied to val( therefore Scan The address specified during can be reused )
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				//  Character slice to string
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}

		res = append(res, entry)
	}

	pager.ItemsPerPage = queryBuilder.LimitCount
	pager.CurrentPageItems = len(res)

	return &res, &pager, nil
}

func (queryBuilder *QueryBuilder) Insert() (int64, *exception.AppException) {
	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	res, err := queryBuilder.Db.ExecContext(ctx, queryBuilder.MakeInsertSql())

	if err != nil {
		return 0, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)

	return lastId, nil
}
