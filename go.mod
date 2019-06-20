module gorose

go 1.12

require (
	github.com/gohouse/gocar v0.0.0-20190620102130-2b3fd7556989
	github.com/gohouse/t v0.0.4
	github.com/mattn/go-sqlite3 v1.10.0
)

replace (
	github.com/gohouse/gocar => ../gocar
	github.com/gohouse/t => ../t
)
