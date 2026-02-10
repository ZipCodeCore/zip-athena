package database

import (
	"fmt"
	"strings"
)

// QueryBuilder helps construct SQL queries in a type-safe and composable way
type QueryBuilder struct {
	table      string
	columns    []string
	conditions []string
	args       []interface{}
	orderBy    string
	limit      int
	offset     int
}

// NewQuery creates a new QueryBuilder for the specified table
func NewQuery(table string) *QueryBuilder {
	return &QueryBuilder{
		table:      table,
		columns:    []string{"*"},
		conditions: []string{},
		args:       []interface{}{},
	}
}

// Select specifies which columns to select
func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
	q.columns = columns
	return q
}

// Where adds a WHERE condition with arguments
func (q *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	q.conditions = append(q.conditions, condition)
	q.args = append(q.args, args...)
	return q
}

// OrderBy sets the ORDER BY clause
func (q *QueryBuilder) OrderBy(order string) *QueryBuilder {
	q.orderBy = order
	return q
}

// Limit sets the LIMIT clause
func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	q.limit = n
	return q
}

// Offset sets the OFFSET clause
func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	q.offset = n
	return q
}

// Build constructs the final SQL query and returns it with the arguments
func (q *QueryBuilder) Build() (string, []interface{}) {
	// Build SELECT clause
	sql := "SELECT " + strings.Join(q.columns, ", ") + " FROM " + q.table

	// Add WHERE clause if conditions exist
	if len(q.conditions) > 0 {
		sql += " WHERE " + strings.Join(q.conditions, " AND ")
	}

	// Add ORDER BY clause
	if q.orderBy != "" {
		sql += " ORDER BY " + q.orderBy
	}

	// Add LIMIT clause
	if q.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	// Add OFFSET clause
	if q.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	return sql, q.args
}

// BuildCount constructs a COUNT query
func (q *QueryBuilder) BuildCount() (string, []interface{}) {
	sql := "SELECT COUNT(*) FROM " + q.table

	if len(q.conditions) > 0 {
		sql += " WHERE " + strings.Join(q.conditions, " AND ")
	}

	return sql, q.args
}

// UpdateBuilder helps construct SQL UPDATE queries
type UpdateBuilder struct {
	table      string
	sets       []string
	args       []interface{}
	conditions []string
	whereArgs  []interface{}
}

// NewUpdate creates a new UpdateBuilder for the specified table
func NewUpdate(table string) *UpdateBuilder {
	return &UpdateBuilder{
		table:      table,
		sets:       []string{},
		args:       []interface{}{},
		conditions: []string{},
		whereArgs:  []interface{}{},
	}
}

// Set adds a column to update
func (u *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	u.sets = append(u.sets, column+" = ?")
	u.args = append(u.args, value)
	return u
}

// Where adds a WHERE condition
func (u *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	u.conditions = append(u.conditions, condition)
	u.whereArgs = append(u.whereArgs, args...)
	return u
}

// Build constructs the final UPDATE query
func (u *UpdateBuilder) Build() (string, []interface{}) {
	sql := "UPDATE " + u.table + " SET " + strings.Join(u.sets, ", ")

	if len(u.conditions) > 0 {
		sql += " WHERE " + strings.Join(u.conditions, " AND ")
	}

	// Combine SET args and WHERE args
	allArgs := append(u.args, u.whereArgs...)

	return sql, allArgs
}

// DeleteBuilder helps construct SQL DELETE queries
type DeleteBuilder struct {
	table      string
	conditions []string
	args       []interface{}
}

// NewDelete creates a new DeleteBuilder for the specified table
func NewDelete(table string) *DeleteBuilder {
	return &DeleteBuilder{
		table:      table,
		conditions: []string{},
		args:       []interface{}{},
	}
}

// Where adds a WHERE condition
func (d *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	d.conditions = append(d.conditions, condition)
	d.args = append(d.args, args...)
	return d
}

// Build constructs the final DELETE query
func (d *DeleteBuilder) Build() (string, []interface{}) {
	sql := "DELETE FROM " + d.table

	if len(d.conditions) > 0 {
		sql += " WHERE " + strings.Join(d.conditions, " AND ")
	}

	return sql, d.args
}
