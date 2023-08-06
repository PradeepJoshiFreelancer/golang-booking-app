package forms

type errors map[string][]string

// Add error to the error list
func (e errors) AddError(feild, msg string) {
	e[feild] = append(e[feild], msg)
}

// Get the first error
func (e errors) GetError(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
