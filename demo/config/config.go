package config

//var Config = make(map[string]map[string]string)
var Configs = map[string]map[string]string {
	"mysql_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		"charset":  "utf8",
		"protocol": "tcp",
		"driver":	"mysql",
	},
	"postgres_dev": {
		"host":     "localhost",
		"username": "root",
		"password": "",
		"port":     "3306",
		"database": "test",
		//"charset":  "utf8",
		//"protocol": "tcp",
		"driver":	"postgres",
	},
	"oracle_dev": {
		//"host":     "localhost",
		"username": "root",
		"password": "",
		//"port":     "3306",
		"database": "test",
		//"charset":  "utf8",
		//"protocol": "tcp",
		"driver":	"oracle",
	},
	"sqlite_dev": {
		//"host":     "localhost",
		//"username": "root",
		//"password": "",
		//"port":     "3306",
		"database": "./foo.db",
		//"charset":  "utf8",
		//"protocol": "tcp",
		"driver":	"sqlite",
	},
}
