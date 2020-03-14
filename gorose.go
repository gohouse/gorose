package gorose

// GOROSE_IMG ...
const GOROSE_IMG = `

  ,ad8888ba,                88888888ba
 d8"'    '"8b               88      "8b
d8'                         88      ,8P
88              ,adPPYba,   88aaaaaa8P'  ,adPPYba,   ,adPPYba,   ,adPPYba,
88      88888  a8"     "8a  88""""88'   a8"     "8a  I8[    ""  a8P_____88
Y8,        88  8b       d8  88    '8b   8b       d8   '"Y8ba,   8PP"""""""
 Y8a.    .a88  "8a,   ,a8"  88     '8b  "8a,   ,a8"  aa    ]8I  "8b,   ,aa
  '"Y88888P"    '"YbbdP"'   88      '8b  '"YbbdP"'   '"YbbdP"'   '"Ybbd8"'

`

const (
	// VERSION_TEXT ...
	VERSION_TEXT = "\ngolang orm of gorose's version : "
	// VERSION_NO ...
	VERSION_NO = "v2.1.5"
	// VERSION ...
	VERSION = VERSION_TEXT + VERSION_NO + GOROSE_IMG
)

// Open ...
func Open(conf ...interface{}) (engin *Engin, err error) {
	// 驱动engin
	engin, err = NewEngin(conf...)
	if err != nil {
		if engin.GetLogger().EnableErrorLog() {
			engin.GetLogger().Error(err.Error())
		}
		return
	}

	return
}
