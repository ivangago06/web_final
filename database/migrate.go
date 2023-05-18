package database

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"web_final/app"

	"github.com/lib/pq"
)

func Migrate(app *app.App, anyStruct interface{}) error {
	valueOfStruct := reflect.ValueOf(anyStruct)
	typeOfStruct := valueOfStruct.Type()

	tableName := typeOfStruct.Name()
	err := createTable(app, tableName)
	if err != nil {
		return err
	}

	for i := 0; i < valueOfStruct.NumField(); i++ {
		fieldType := typeOfStruct.Field(i)
		fieldName := fieldType.Name
		if fieldName != "Id" && fieldName != "id" {
			err := createColumn(app, tableName, fieldName, fieldType.Type.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createTable(app *app.App, tableName string) error {
	// Check to see if the table already exists
	var tableExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE c.relname ~ $1 AND pg_catalog.pg_table_is_visible(c.oid))", "^"+tableName+"$").Scan(&tableExists)
	if err != nil {
		log.Println("Error alidando si la tabla existe: " + tableName)
		return err
	}

	if tableExists {
		log.Println("Table already exists: " + tableName)
		return nil
	} else {
		sanitizedTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\"Id\" serial primary key)", tableName)

		_, err := app.Db.Query(sanitizedTableQuery)
		if err != nil {
			log.Println("Error creating table: " + tableName)
			return err
		}

		log.Println("Tabla creada exitosamente: " + tableName)
		return nil
	}
}

func createColumn(app *app.App, tableName, columnName, columnType string) error {

	var columnExists bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", tableName, columnName).Scan(&columnExists)
	if err != nil {
		log.Println("Error checking if column exists: " + columnName + " in table: " + tableName)
		return err
	}

	if columnExists {
		log.Println("La columna ya existe: " + columnName + " in table: " + tableName)
		return nil
	} else {
		postgresType, err := getPostgresType(columnType)
		if err != nil {
			log.Println("Error creando la columna: " + columnName + " en la tabla: " + tableName + " de tipo: " + postgresType)
			return err
		}

		sanitizedTableName := pq.QuoteIdentifier(tableName)
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS \"%s\" %s", sanitizedTableName, columnName, postgresType)

		_, err = app.Db.Query(query)
		if err != nil {
			log.Println("Error creando la columna: " + columnName + " en la tabla: " + tableName + " de tipo: " + postgresType)
			return err
		}

		log.Println("Columna creada exitosamente:", columnName)

		return nil
	}
}

func getPostgresType(goType string) (string, error) {
	switch goType {
	case "int", "int32", "uint", "uint32":
		return "integer", nil
	case "int64", "uint64":
		return "bigint", nil
	case "int16", "int8", "uint16", "uint8", "byte":
		return "smallint", nil
	case "string":
		return "text", nil
	case "float64":
		return "double precision", nil
	case "bool":
		return "boolean", nil
	case "Time":
		return "timestamp", nil
	case "[]byte":
		return "bytea", nil
	}

	return "", errors.New("Tipo no reconocido: " + goType)
}
