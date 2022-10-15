package utils

const (
	QueryGet = "SELECT %s %s FROM %s %s "

	QueryGetCount = "SELECT COUNT('total') FROM (%s) as total"

	QueryPaginationOrdinal = "(%s > %s) ORDER BY %s LIMIT %s"

	QueryPaginationDefault = "LIMIT %s OFFSET %s"

	QueryInsert = "INSERT INTO %s (%s) VALUES (%s);"

	QueryUpdate = "UPDATE %s SET %s WHERE %s = ?; "

	QueryDelete = "DELETE FROM %s WHERE %s = ?; "
)
