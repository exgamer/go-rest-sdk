package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/exception"
	"github.com/exgamer/go-rest-sdk/pkg/form"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewModifyQueryBuilder() *ModifyQueryBuilder {
	return &ModifyQueryBuilder{
		timeout: 30,
		DbType:  "mysql",
	}
}

func NewMysqlModifyQueryBuilder() *ModifyQueryBuilder {
	return &ModifyQueryBuilder{
		timeout: 30,
		DbType:  "postgres",
	}
}

func NewPostgresModifyQueryBuilder() *ModifyQueryBuilder {
	return &ModifyQueryBuilder{
		timeout: 30,
		DbType:  "postgres",
	}
}

// QueryBuilder - query builder
type ModifyQueryBuilder struct {
	Db               *sql.DB
	timeout          time.Duration
	ValuePlaceholder string
	DbType           string
	QueryData        QueryData
}

// SetDb - set sql.DB connection (for make query)
func (queryBuilder *ModifyQueryBuilder) SetDb(db *sql.DB) *ModifyQueryBuilder {
	queryBuilder.Db = db

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) SetDbType(dbType string) *ModifyQueryBuilder {
	queryBuilder.QueryData.DbType = dbType

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) getPlaceholder() string {
	if queryBuilder.QueryData.DbType == "postgres" {
		index := len(queryBuilder.QueryData.Params) + 1

		return "$" + strconv.Itoa(index)
	}

	return "?"
}

// SetData - set data to query
func (queryBuilder *ModifyQueryBuilder) SetData(data map[string]any) *ModifyQueryBuilder {
	queryBuilder.QueryData.Data = data
	fmt.Println(data)

	return queryBuilder
}

// SetFormData - set data to query (for insert/update)
func (queryBuilder *ModifyQueryBuilder) SetFormData(form form.Form) *ModifyQueryBuilder {
	queryBuilder.QueryData.Data = form.AsMap()

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) SetQueryTimeout(timeout time.Duration) *ModifyQueryBuilder {
	queryBuilder.timeout = timeout

	return queryBuilder
}

// Table - set table name
func (queryBuilder *ModifyQueryBuilder) Table(table string) *ModifyQueryBuilder {
	queryBuilder.QueryData.TableName = table

	return queryBuilder
}

// Select - set select fields
func (queryBuilder *ModifyQueryBuilder) Select(fields ...string) *ModifyQueryBuilder {
	queryBuilder.QueryData.SelectArray = fields

	return queryBuilder
}

// TableAlias - set from table alias
func (queryBuilder *ModifyQueryBuilder) TableAlias(tableAlias string) *ModifyQueryBuilder {
	queryBuilder.QueryData.TableAliasName = tableAlias

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) CalculateRows() *ModifyQueryBuilder {
	queryBuilder.QueryData.CalcRows = true

	return queryBuilder
}

// SetLimitOffsetByPage - set limit and offset to query by page
func (queryBuilder *ModifyQueryBuilder) SetLimitOffsetByPage(page int, perPage int) *ModifyQueryBuilder {
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
func (queryBuilder *ModifyQueryBuilder) Limit(limit int) *ModifyQueryBuilder {
	queryBuilder.QueryData.LimitCount = limit

	return queryBuilder
}

// Offset - set offset to query
func (queryBuilder *ModifyQueryBuilder) Offset(offset int) *ModifyQueryBuilder {
	queryBuilder.QueryData.OffsetCount = offset

	return queryBuilder
}

// OrderByAsc - add order by asc
func (queryBuilder *ModifyQueryBuilder) OrderByAsc(orderBy string) *ModifyQueryBuilder {

	return queryBuilder.orderBy(orderBy, "ASC")
}

// OrderByDesc - add order by desc
func (queryBuilder *ModifyQueryBuilder) OrderByDesc(orderBy string) *ModifyQueryBuilder {

	return queryBuilder.orderBy(orderBy, "DESC")
}

// orderBy - add order by
func (queryBuilder *ModifyQueryBuilder) orderBy(orderBy string, sort string) *ModifyQueryBuilder {
	queryBuilder.QueryData.Order = append(queryBuilder.QueryData.Order, orderBy+" "+sort)

	return queryBuilder
}

// GroupBy - add group by
func (queryBuilder *ModifyQueryBuilder) GroupBy(groupBy string) *ModifyQueryBuilder {
	queryBuilder.QueryData.Group = groupBy

	return queryBuilder
}

// GetParams - returns query params
func (queryBuilder *ModifyQueryBuilder) GetParams() []any {

	return queryBuilder.QueryData.Params
}

func (queryBuilder *ModifyQueryBuilder) MakeFoundRowsSql() string {
	return "SELECT FOUND_ROWS()"
}

// Join - add join condition
func (queryBuilder *ModifyQueryBuilder) Join(table string, first string, second string, params ...any) *ModifyQueryBuilder {
	queryBuilder.join(table, first+"="+second, "JOIN", params...)

	return queryBuilder
}

// OuterJoin - add outer join condition
func (queryBuilder *ModifyQueryBuilder) OuterJoin(table string, first string, second string, params ...any) *ModifyQueryBuilder {
	queryBuilder.join(table, first+"="+second, "OUTER JOIN", params...)

	return queryBuilder
}

// InnerJoin - add inner join condition
func (queryBuilder *ModifyQueryBuilder) InnerJoin(table string, first string, second string, params ...any) *ModifyQueryBuilder {
	queryBuilder.join(table, first+"="+second, "INNER JOIN", params...)

	return queryBuilder
}

// RightJoin - add right join condition
func (queryBuilder *ModifyQueryBuilder) RightJoin(table string, first string, second string, params ...any) *ModifyQueryBuilder {
	queryBuilder.join(table, first+"="+second, "RIGHT JOIN", params...)

	return queryBuilder
}

// LeftJoin - add left join condition
func (queryBuilder *ModifyQueryBuilder) LeftJoin(table string, first string, second string, params ...any) *ModifyQueryBuilder {
	queryBuilder.join(table, first+"="+second, "LEFT JOIN", params...)

	return queryBuilder
}

// join - add join condition by string condition
func (queryBuilder *ModifyQueryBuilder) join(table string, on string, joinType string, params ...any) {
	join := JoinCondition{Table: table, On: on, Type: joinType}
	queryBuilder.QueryData.JoinCondition = append(queryBuilder.QueryData.JoinCondition, join)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}
}

func (queryBuilder *ModifyQueryBuilder) OuterJoinWithCondition(table string, first string, second string, condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.joinWithCondition(table, first+"="+second, "OUTER JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) InnerJoinWithCondition(table string, first string, second string, condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.joinWithCondition(table, first+"="+second, "INNER JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) RightJoinWithCondition(table string, first string, second string, condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.joinWithCondition(table, first+"="+second, "RIGHT JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) LeftJoinWithCondition(table string, first string, second string, condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.joinWithCondition(table, first+"="+second, "LEFT JOIN", condition, params...)

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) joinWithCondition(table string, on string, joinType string, condition string, params ...any) {
	join := JoinCondition{Table: table, On: on, Type: joinType, Where: condition}
	queryBuilder.QueryData.JoinCondition = append(queryBuilder.QueryData.JoinCondition, join)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}
}

// OrWhere - adds or where condition
func (queryBuilder *ModifyQueryBuilder) OrWhere(field string, value any) *ModifyQueryBuilder {
	queryBuilder.OrWhereByCondition(field+"="+queryBuilder.getPlaceholder(), value)

	return queryBuilder
}

// OrWhereByCondition - adds or where condition  by string
func (queryBuilder *ModifyQueryBuilder) OrWhereByCondition(condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.WhereByCondition(condition, "OR", params...)

	return queryBuilder
}

// AndWhere - adds and where condition
func (queryBuilder *ModifyQueryBuilder) AndWhere(field string, value any) *ModifyQueryBuilder {
	queryBuilder.AndWhereByCondition(field+"="+queryBuilder.getPlaceholder(), value)

	return queryBuilder
}

// AndWhereByCondition - adds where condition by string
func (queryBuilder *ModifyQueryBuilder) AndWhereByCondition(condition string, params ...any) *ModifyQueryBuilder {
	queryBuilder.WhereByCondition(condition, "AND", params...)

	return queryBuilder
}

// AndWhereIn - adds and in condition
func (queryBuilder *ModifyQueryBuilder) AndWhereIn(field string, params []any) *ModifyQueryBuilder {
	return queryBuilder.WhereIn(field, params, "AND")
}

// WhereIn - adds in condition
func (queryBuilder *ModifyQueryBuilder) WhereIn(field string, params []any, operator string) *ModifyQueryBuilder {
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
func (queryBuilder *ModifyQueryBuilder) WhereByCondition(condition string, operator string, params ...any) *ModifyQueryBuilder {
	where := WhereCondition{Condition: condition, Operator: operator}
	queryBuilder.QueryData.WhereCondition = append(queryBuilder.QueryData.WhereCondition, where)
	fmt.Print(params)

	for _, p := range params {
		queryBuilder.QueryData.Params = append(queryBuilder.QueryData.Params, p)
	}

	return queryBuilder
}

func (queryBuilder *ModifyQueryBuilder) Insert() (int64, *exception.AppException) {
	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	res, err := queryBuilder.Db.ExecContext(ctx, MakeInsertSql(&queryBuilder.QueryData))

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

func (queryBuilder *ModifyQueryBuilder) Update() (bool, *exception.AppException) {
	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	_, err := queryBuilder.Db.ExecContext(ctx, MakeUpdateSql(&queryBuilder.QueryData), queryBuilder.GetParams()...)

	if err != nil {
		return false, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	return true, nil
}

func (queryBuilder *ModifyQueryBuilder) Delete() *exception.AppException {
	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	_, err := queryBuilder.Db.ExecContext(ctx, MakeDeleteSql(&queryBuilder.QueryData), queryBuilder.GetParams()...)

	if err != nil {
		return exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	return nil
}

func (queryBuilder *ModifyQueryBuilder) Exec(sql string, params ...any) (bool, *exception.AppException) {
	ctx, cancel := context.WithTimeout(context.Background(), queryBuilder.timeout*time.Second)
	defer cancel()

	_, err := queryBuilder.Db.ExecContext(ctx, sql, params...)

	if err != nil {
		return false, exception.NewAppException(http.StatusInternalServerError, err, nil)
	}

	return true, nil
}
