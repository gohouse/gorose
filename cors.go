package gorose

import "github.com/gohouse/gorose/cors"

func BootLogger() func(*Connection) {
	return func(conn *Connection) {
		conn.Logger = cors.NewDefaultLogger()
	}
}

func NewLogger() func(*Connection) {
	return func(conn *Connection) {
		conn.Logger = cors.NewDefaultLogger()
	}
}
