package factory

type Factory interface {
	GetData() interface{} //harusnya return interface of models
	RunFactory(count int) error
}
