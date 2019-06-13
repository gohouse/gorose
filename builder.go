package gorose

type Builder struct {
	//driver string
}

func NewBuilder(driver string) IBuilder {
	return NewBuilderDriver().Getter(driver)
}
