package gorose

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
		"in", "not in", "between", "not between"}
)

type Builder struct {
}

func NewBuilder(driver string) IBuilder {
	return NewBuilderDriver().Getter(driver)
}
