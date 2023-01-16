package factory

type Factory interface {
	GetData() interface{} //harusnya return interface of models
}
