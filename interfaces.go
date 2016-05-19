package pressmark

type BeforeCreateCallback interface {
	BeforeCreate() error
}

type BeforeSaveCallback interface {
	BeforeSave() error
}

type BeforeUpdateCallback interface {
	BeforeUpdate() error
}

type Initiator interface {
	Init() 
}
