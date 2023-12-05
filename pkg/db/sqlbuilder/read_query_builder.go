package sqlbuilder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/form"
	"github.com/exgamer/go-rest-sdk/pkg/logger"
	paginator "github.com/exgamer/go-rest-sdk/pkg/pager"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewReadQueryBuilder[E interface{}]() *ReadQueryBuilder[E] {
	return &ReadQueryBuilder[E]{
		timeout: 30,
		DbType:  "mysql",
	}
}

func NewMysqlReadQueryBuilder[E interface{}]() *ReadQueryBuilder[E] {
	qb := &ReadQueryBuilder[E]{
		timeout: 30,
		DbType:  "mysql",
	}

	qb.QueryData.DbType = qb.DbType

	return qb
}

func NewPostgresReadQueryBuilder[E interface{}]() *ReadQueryBuilder[E] {
	qb := &ReadQueryBuilder[E]{
		timeout: 30,
		DbType:  "postgres",
	}

	qb.QueryData.DbType = qb.DbType

	return qb
}

// QueryBuilder - query builder
type ReadQueryBuilder[E interface{}] struct {
	Db               *sql.DB
	timeout          time.Duration
	ValuePlaceholder string
	DbType           string
	Entities         []E
	QueryData        QueryData
}

type QueryData struct {
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
	Where string
}

// SetDb - set sql.DB connection (for make query)
func (queryBuilder *ReadQueryBuilder[E]) SetDb(db *sql.DB) *ReadQueryBuilder[E] {
	queryBuilder.Db = db

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) SetDbType(dbType string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.DbType = dbType

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) getPlaceholder() string {
	if queryBuilder.QueryData.DbType == "postgres" {
		index := len(queryBuilder.QueryData.Params) + 1

		return "$" + strconv.Itoa(index)
	}

	return "?"
}

// SetData - set data to query
func (queryBuilder *ReadQueryBuilder[E]) SetData(data map[string]any) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.Data = data
	fmt.Println(data)

	return queryBuilder
}

// SetFormData - set data to query (for insert/update)
func (queryBuilder *ReadQueryBuilder[E]) SetFormData(form form.Form) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.Data = form.AsMap()

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) SetQueryTimeout(timeout time.Duration) *ReadQueryBuilder[E] {
	queryBuilder.timeout = timeout

	return queryBuilder
}

// Table - set table name
func (queryBuilder *ReadQueryBuilder[E]) Table(table string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.TableName = table

	return queryBuilder
}

// Select - set select fields
func (queryBuilder *ReadQueryBuilder[E]) Select(fields ...string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.SelectArray = fields

	return queryBuilder
}

// TableAlias - set from table alias
func (queryBuilder *ReadQueryBuilder[E]) TableAlias(tableAlias string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.TableAliasName = tableAlias

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) CalculateRows() *ReadQueryBuilder[E] {
	queryBuilder.QueryData.CalcRows = true

	return queryBuilder
}

// SetLimitOffsetByPage - set limit and offset to query by page
func (queryBuilder *ReadQueryBuilder[E]) SetLimitOffsetByPage(page int, perPage int) *ReadQueryBuilder[E] {
	var limit = 30
	var offset = 0

	if perPage > 0 {
		limit = perPage
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	queryBuilder.QueryData.LimitCount = limit
	queryBuilder.QueryData.OffsetCount = offset

	return queryBuilder
}

// Limit - set limit to query
func (queryBuilder *ReadQueryBuilder[E]) Limit(limit int) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.LimitCount = limit

	return queryBuilder
}

// Offset - set offset to query
func (queryBuilder *ReadQueryBuilder[E]) Offset(offset int) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.OffsetCount = offset

	return queryBuilder
}

// OrderByAsc - add order by asc
func (queryBuilder *ReadQueryBuilder[E]) OrderByAsc(orderBy string) *ReadQueryBuilder[E] {

	return queryBuilder.orderBy(orderBy, "ASC")
}

// OrderByDesc - add order by desc
func (queryBuilder *ReadQueryBuilder[E]) OrderByDesc(orderBy string) *ReadQueryBuilder[E] {

	return queryBuilder.orderBy(orderBy, "DESC")
}

// orderBy - add order by
func (queryBuilder *ReadQueryBuilder[E]) orderBy(orderBy string, sort string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.Order = append(queryBuilder.QueryData.Order, orderBy+" "+sort)

	return queryBuilder
}

// GroupBy - add group by
func (queryBuilder *ReadQueryBuilder[E]) GroupBy(groupBy string) *ReadQueryBuilder[E] {
	queryBuilder.QueryData.Group = groupBy

	return queryBuilder
}

// GetParams - returns query params
func (queryBuilder *ReadQueryBuilder[E]) GetParams() []any {

	return queryBuilder.QueryData.Params
}

func (queryBuilder *ReadQueryBuilder[E]) MakeFoundRowsSql() string {
	return "SELECT FOUND_ROWS()"
}

// Join - add join condition
func (queryBuilder *ReadQueryBuilder[E]) Join(table string, first string, second string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.join(table, first+"="+second, "JOIN", params...)

	return queryBuilder
}

// OuterJoin - add outer join condition
func (queryBuilder *ReadQueryBuilder[E]) OuterJoin(table string, first string, second string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.join(table, first+"="+second, "OUTER JOIN", params...)

	return queryBuilder
}

// InnerJoin - add inner join condition
func (queryBuilder *ReadQueryBuilder[E]) InnerJoin(table string, first string, second string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.join(table, first+"="+second, "INNER JOIN", params...)

	return queryBuilder
}

// RightJoin - add right join condition
func (queryBuilder *ReadQueryBuilder[E]) RightJoin(table string, first string, second string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.join(table, first+"="+second, "RIGHT JOIN", params...)

	return queryBuilder
}

// LeftJoin - add left join condition
func (queryBuilder *ReadQueryBuilder[E]) LeftJoin(table string, first string, second string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.join(table, first+"="+second, "LEFT JOIN", params...)

	return queryBuilder
}

// join - add join condition by string condition
func (queryBuilder *ReadQueryBuilder[E]) join(table string, on string, joinType string, params ...any) {
	join := JoinCondition{Table: table, On: on, Type: joinType}
	queryBuilder.QueryData.JoinCondition = append(queryBuilder.QueryData.JoinCondition, join)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}
}

func (queryBuilder *ReadQueryBuilder[E]) OuterJoinWithCondition(table string, first string, second string, condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.joinWithCondition(table, first+"="+second, "OUTER JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) InnerJoinWithCondition(table string, first string, second string, condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.joinWithCondition(table, first+"="+second, "INNER JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) RightJoinWithCondition(table string, first string, second string, condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.joinWithCondition(table, first+"="+second, "RIGHT JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) LeftJoinWithCondition(table string, first string, second string, condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.joinWithCondition(table, first+"="+second, "LEFT JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ReadQueryBuilder[E]) joinWithCondition(table string, on string, joinType string, condition string, params ...any) {
	join := JoinCondition{Table: table, On: on, Type: joinType, Where: condition}
	queryBuilder.QueryData.JoinCondition = append(queryBuilder.QueryData.JoinCondition, join)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}
}

// OrWhere - adds or where condition
func (queryBuilder *ReadQueryBuilder[E]) OrWhere(field string, value any) *ReadQueryBuilder[E] {
	queryBuilder.OrWhereByCondition(field+"="+queryBuilder.getPlaceholder(), value)

	return queryBuilder
}

// OrWhereByCondition - adds or where condition  by string
func (queryBuilder *ReadQueryBuilder[E]) OrWhereByCondition(condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.WhereByCondition(condition, "OR", params...)

	return queryBuilder
}

// AndWhere - adds and where condition
func (queryBuilder *ReadQueryBuilder[E]) AndWhere(field string, value any) *ReadQueryBuilder[E] {
	queryBuilder.AndWhereByCondition(field+"="+queryBuilder.getPlaceholder(), value)

	return queryBuilder
}

// AndWhereByCondition - adds where condition by string
func (queryBuilder *ReadQueryBuilder[E]) AndWhereByCondition(condition string, params ...any) *ReadQueryBuilder[E] {
	queryBuilder.WhereByCondition(condition, "AND", params...)

	return queryBuilder
}

// AndWhereIn - adds and in condition
func (queryBuilder *ReadQueryBuilder[E]) AndWhereIn(field string, params []any) *ReadQueryBuilder[E] {
	return queryBuilder.WhereIn(field, params, "AND")
}

// WhereIn - adds in condition
func (queryBuilder *ReadQueryBuilder[E]) WhereIn(field string, params []any, operator string) *ReadQueryBuilder[E] {
	sParams := make([]string, len(params))

	for i := 0; i < len(params); i++ {
		sParams[i] = queryBuilder.getPlaceholder()
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, params[i])
	}

	condition := field + " IN (" + strings.Join(sParams, ",") + ")"
	where := WhereCondition{Condition: condition, Operator: operator}
	queryBuilder.QueryData.WhereCondition = append(queryBuilder.QueryData.WhereCondition, where)

	return queryBuilder
}

// WhereByCondition - adds string condition
func (queryBuilder *ReadQueryBuilder[E]) WhereByCondition(condition string, operator string, params ...any) *ReadQueryBuilder[E] {
	where := WhereCondition{Condition: condition, Operator: operator}
	queryBuilder.QueryData.WhereCondition = append(queryBuilder.QueryData.WhereCondition, where)
	fmt.Print(params)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}

	return queryBuilder
}

// All - returns all query result
func (queryBuilder *ReadQueryBuilder[E]) All() (*[]E, *exception.AppException) {
	res, _, appException := queryBuilder.paginate(false, 0, 1000)

	return res, appException
}

// One - returns one query result
func (queryBuilder *ReadQueryBuilder[E]) One() (*E, *exception.AppException) {
	res, _, appException := queryBuilder.paginate(false, 0, 1)

	if len(*res) == 0 {
		return nil, nil
	}

	return &(*res)[0], appException
}

// Paginate - returns paginated query result
func (queryBuilder *ReadQueryBuilder[E]) Paginate(page int, perPage int) (*[]E, *paginator.Pager, *exception.AppException) {
	return queryBuilder.paginate(true, page, perPage)
}

func (queryBuilder *ReadQueryBuilder[E]) paginate(paginate bool, page int, perPage int) (*[]E, *paginator.Pager, *exception.AppException) {
	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	res := make([]map[string]interface{}, 0)
	pager := paginator.Pager{
		Page: page,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	if paginate == true {
		err := queryBuilder.Db.QueryRowContext(ctx, MakeCountSelectSql(&queryBuilder.QueryData), queryBuilder.GetParams()...).Scan(&pager.TotalItems)
		if err != nil {
			logger.LogError(err)

			return nil, &pager, exception.NewAppException(http.StatusInternalServerError, err, nil)
		}

		if pager.TotalItems == 0 {
			return &[]E{}, &pager, nil
		}
	}

	queryBuilder.SetLimitOffsetByPage(page, perPage)
	// Execute the query
	rows, err := queryBuilder.Db.QueryContext(ctx, MakeSelectSql(&queryBuilder.QueryData), queryBuilder.GetParams()...)

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

	pager.ItemsPerPage = queryBuilder.QueryData.LimitCount
	pager.CurrentPageItems = len(res)

	jsonRes, err := json.Marshal(res)

	if err != nil {
		logger.LogError(err)

		return nil, &pager, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	resCount := len(res)
	entities := make([]E, resCount)
	json.Unmarshal(jsonRes, &entities)

	return &entities, &pager, nil
}
