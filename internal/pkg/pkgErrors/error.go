package pkgErrors

import "errors"

var (
	ErrNotFound         = errors.New("Not Found")
	ErrDbConfigNotFound = errors.New("one or more database configuration parameters are empty")
	ErrInvalidName      = errors.New("Invalid Data")
	ErrForeignKey       = errors.New("ERROR: insert or update on table \"class_info\" violates foreign key constraint \"fk_student\" (SQLSTATE 23503)")
	ErrParse            = errors.New("Could not get DB_PORT:")
)
