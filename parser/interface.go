package parser

type IParser interface {
	Parse(d string, dbConfCluster interface{}) error
}
