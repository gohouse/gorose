package gorose

func NewBuilder(driver string) IBuilder {
	return NewBuilderDriver().Getter(driver)
}
