package log

import (
	"log/slog"
)

// SQLSchema returns a standardized group of attributes
// to log a Wide Event involving a SQL query.
func SQLSchema(
	queryText string,
	table string,
	operationName string,
	procedureName string,
	returnedRows int,
) slog.Attr {
	return slog.GroupAttrs("sql",
		slog.String("query", queryText),
		slog.String("affected_table", table),
		slog.String("operation", operationName),
		slog.String("stored_procedure", procedureName),
		slog.Int("returned_rows", returnedRows),
	)
}
