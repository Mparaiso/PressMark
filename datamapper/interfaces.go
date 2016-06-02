package datamapper

type DataMapperMetadataProvider interface {
	DataMapperMetaData() DataMapperMetadata
}

type BeforeCreateCallback interface {
	BeforeCreate() error
}

type BeforeSaveCallback interface {
	BeforeSave() error
}

type BeforeUpdateCallback interface {
	BeforeUpdate() error
}

type QueryBuilder interface {
	AcceptRepository(*Repository) (string, []interface{})
}
