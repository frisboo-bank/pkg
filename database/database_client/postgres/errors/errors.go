package errors

import "fmt"

const errorPrefix = "Postgres:"

/**
 *
 * Postgres failed to connect error
 *
 */
type PostgresFailedToConnectError struct {
	Err error
}

func (e *PostgresFailedToConnectError) Error() string {
	return fmt.Sprintf("%s failed to connect to the database with error: %v", errorPrefix, e.Err)
}

/**
 *
 * Postgres failed to ping error
 *
 */
type PostgresFailedToPingError struct {
	Err error
}

func (e *PostgresFailedToPingError) Error() string {
	return fmt.Sprintf("%s failed to ping the database and received the error: %v", errorPrefix, e.Err)
}

/**
 *
 * Postgres failed to close DB connection error
 *
 */
type PostgresFailedToCloseDBConnectionError struct {
	Err error
}

func (e *PostgresFailedToCloseDBConnectionError) Error() string {
	return fmt.Sprintf("%s failed to close the DB connection with error: %v", errorPrefix, e.Err)
}

/**
 *
 * Postgres failed to close Sqlx DB connection error
 *
 */
type PostgresFailedToCloseSqlxDBConnectionError struct {
	Err error
}

func (e *PostgresFailedToCloseSqlxDBConnectionError) Error() string {
	return fmt.Sprintf("%s failed to close the Sqlx DB connection with error: %v", errorPrefix, e.Err)
}

/**
 *
 * Postgres failed to execute transaction error
 *
 */
type PostgresFailedToExecuteTransactionError struct {
	Err error
}

func (e *PostgresFailedToExecuteTransactionError) Error() string {
	return fmt.Sprintf("%s failed to execute transaction with error: %v", errorPrefix, e.Err)
}

/**
 *
 * Postgres failed to roolback transaction error
 *
 */
type PostgresFailedToRollbackTransactionError struct {
	Err         error
	OriginalErr error
}

func (e *PostgresFailedToRollbackTransactionError) Error() string {
	return fmt.Sprintf(
		"%s failed to execute transaction with error: %v and the rollback also failed with error: %v",
		errorPrefix,
		e.OriginalErr,
		e.Err,
	)
}

/**
 *
 * Postgres failed to commit transaction error
 *
 */
type PostgresFailedToCommitTransactionError struct {
	Err error
}

func (e *PostgresFailedToCommitTransactionError) Error() string {
	return fmt.Sprintf("%s failed to commit transaction with error: %v", errorPrefix, e.Err)
}
