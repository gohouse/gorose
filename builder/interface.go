package builder

type IBuilder interface {
	BuildQuery() (string, error)
	BuildExecute() (string, error)
}
