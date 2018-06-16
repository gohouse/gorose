package drivers

func GetDsnByDriverName(dbObj map[string]string) (driver string, dsn string) {
	switch dbObj["driver"] {
	case "mysql":
		driver, dsn = MySQL(dbObj)
	case "sqlite3":
		driver, dsn = Sqlite3(dbObj)
	case "postgres":
		driver, dsn = Postgres(dbObj)
	case "oracle":
		driver, dsn = Oracle(dbObj)
	case "mssql":
		driver, dsn = MsSQL(dbObj)
	}
	return
}

