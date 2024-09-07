package mysql

import (
	"fmt"
	"reflect"
	"strings"
)

type JoinOptions struct {
	JoinType string `json:"joinType,omitempty"`
	On       string `json:"on,omitempty"`
	Table    string `json:"table,omitempty"`
}

type QueryBuilder struct {
	columns      []string
	values       []interface{}
	setClauses   []string
	conditions   []string
	conjunctions []string
	args         []interface{}
	joins        []JoinOptions
	limit        int
	offset       int
	orderBy      string
	setArgs      []interface{}
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (b *QueryBuilder) Select(columns []string) *QueryBuilder {
	b.columns = columns
	return b
}

func (b *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	if len(b.conditions) > 0 {
		b.conjunctions = append(b.conjunctions, "AND")
	}
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *QueryBuilder) WhereWithConjunction(conjunction, condition string, args ...interface{}) *QueryBuilder {
	if len(b.conditions) > 0 {
		b.conjunctions = append(b.conjunctions, conjunction)
	}
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *QueryBuilder) Join(joinType, table, on string) *QueryBuilder {
	b.joins = append(b.joins, JoinOptions{
		JoinType: joinType,
		Table:    table,
		On:       on,
	})
	return b
}

func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
	b.limit = limit
	return b
}

func (b *QueryBuilder) Offset(offset int) *QueryBuilder {
	b.offset = offset
	return b
}

func (b *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	b.orderBy = orderBy
	return b
}

func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := strings.Repeat("?,", len(values)-1) + "?"
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IN (%s)", column, placeholders))
	qb.args = append(qb.args, values...)
	return qb
}

func (qb *QueryBuilder) Set(column string, value interface{}) *QueryBuilder {
	qb.setClauses = append(qb.setClauses, fmt.Sprintf("%s = ?", column))
	qb.setArgs = append(qb.setArgs, value)

	return qb
}

func (b *QueryBuilder) AddColumnValue(column string, value interface{}) *QueryBuilder {
	b.columns = append(b.columns, column)
	b.values = append(b.values, value)
	return b
}

// func (qb *QueryBuilder) Search(conditions map[string][]string) *QueryBuilder {
// 	conditionStrings := make([]string, 0)
// 	args := []interface{}{}

// 	for column, values := range conditions {
// 		for _, value := range values {
// 			conditionStrings = append(conditionStrings, fmt.Sprintf("%s LIKE ?", column))
// 			args = append(args, "%"+value+"%")
// 		}
// 	}

// 	if len(conditionStrings) > 0 {
// 		qb.conditions = append(qb.conditions, strings.Join(conditionStrings, " OR "))
// 		qb.args = append(qb.args, args...)
// 	}

// 	return qb
// }

func (qb *QueryBuilder) Search(columns []string, searchText string) *QueryBuilder {
	if len(columns) == 0 || searchText == "" {
		return qb
	}

	conditionStrings := make([]string, 0, len(columns))
	args := make([]interface{}, 0, len(columns))

	for _, column := range columns {
		conditionStrings = append(conditionStrings, fmt.Sprintf("%s LIKE ?", column))
		args = append(args, "%"+searchText+"%")
	}

	if len(conditionStrings) > 0 {
		qb.conditions = append(qb.conditions, strings.Join(conditionStrings, " OR "))
		qb.args = append(qb.args, args...)
	}

	return qb
}

func (b *QueryBuilder) BuildSelectQuery(tableName string) (string, []interface{}) {
	columns := joinColumns(b.columns)
	joinClause := buildJoinClause(b.joins)
	whereClause := buildWhereClause(b.conditions)

	query := fmt.Sprintf("SELECT %s FROM %s %s %s", columns, tableName, joinClause, whereClause)
	return query, b.args
}

func (b *QueryBuilder) BuildSelectManyQuery(tableName string) (string, []interface{}) {
	columns := joinColumns(b.columns)
	joinClause := buildJoinClause(b.joins)
	whereClause := buildWhereClause(b.conditions)
	paginationClause := buildPaginationClause(b.limit, b.offset, b.orderBy)

	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s", columns, tableName, joinClause, whereClause, paginationClause)
	return query, b.args
}

func (b *QueryBuilder) BuildUpdateQuery(tableName string) (string, []interface{}) {
	setClause := strings.Join(b.setClauses, ", ")
	whereClause := buildWhereClause(b.conditions)

	query := fmt.Sprintf("UPDATE %s SET %s %s", tableName, setClause, whereClause)
	return query, b.args
}

func (qb *QueryBuilder) BuildUpdateManyQuery(tableName string) (string, []interface{}) {
	setClause := strings.Join(qb.setClauses, ", ")
	whereClause := strings.Join(qb.conditions, " AND ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClause, whereClause)
	args := append(qb.setArgs, qb.args...)
	return query, args
}

func (b *QueryBuilder) BuildInsertQuery(tableName string) (string, []interface{}) {
	columns := strings.Join(b.columns, ", ")
	placeholders := strings.Repeat("?, ", len(b.values)-1) + "?"
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
	return query, b.values
}

func (b *QueryBuilder) BuildDeleteQuery(tableName string) (string, []interface{}) {
	whereClause := buildWhereClause(b.conditions)
	query := fmt.Sprintf("DELETE FROM %s %s", tableName, whereClause)
	return query, b.args
}

func (qb *QueryBuilder) BuildCountQuery(tableName string) (string, []interface{}) {
	baseQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	whereClause := buildWhereClause(qb.conditions)
	if whereClause != "" {
		baseQuery = fmt.Sprintf("%s %s", baseQuery, whereClause)
	}

	return baseQuery, qb.args
}

func (qb *QueryBuilder) ExtractFieldsForInsert(data interface{}) error {
	// Ensure the data is a pointer to a struct
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	// Dereference the pointer to access the struct
	val = val.Elem()
	typ := val.Type()

	// Iterate over all fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		fieldValue := reflect.ValueOf(data).Elem().Field(i)

		// Get the `db` tag to use as the column name
		columnName := fieldType.Tag.Get("db")
		if columnName == "" || fieldValue.IsZero() || columnName == "created_at" || columnName == "updated_at" || strings.Contains(columnName, "status") || columnName == "id" || (field.Kind() == reflect.Ptr && field.IsNil()) {
			continue // Skip fields without a tag or that are nil
		}

		// Add the column and value to the query builder
		qb.AddColumnValue(fmt.Sprintf("`%s`", columnName), field.Interface())
	}

	return nil
}

func joinColumns(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}
	return strings.Join(columns, ", ")
}

func buildJoinClause(joins []JoinOptions) string {
	if len(joins) == 0 {
		return ""
	}

	joinClauses := make([]string, len(joins))
	for i, join := range joins {
		joinClauses[i] = fmt.Sprintf("%s JOIN %s ON %s", join.JoinType, join.Table, join.On)
	}
	return strings.Join(joinClauses, " ")
}

func buildPaginationClause(limit, offset int, orderBy string) string {
	clauses := []string{}

	if orderBy != "" {
		clauses = append(clauses, fmt.Sprintf("ORDER BY %s", orderBy))
	}

	if limit > 0 {
		clauses = append(clauses, fmt.Sprintf("LIMIT %d", limit))
	}

	if offset > 0 {
		clauses = append(clauses, fmt.Sprintf("OFFSET %d", (offset-1)*(limit)))
	}

	return strings.Join(clauses, " ")
}

func buildWhereClause(conditions []string) string {
	if len(conditions) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
}
