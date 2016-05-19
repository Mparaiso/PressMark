package pressmark
// Validator validates a model
type Validator interface{
    Validate(model interface{}) map[string][]string
}