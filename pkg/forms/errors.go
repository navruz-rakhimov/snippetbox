package forms

// Type of map to hold validation error messages for form fields.
// Map that maps a field name to the slice of error messages.
// There might be multiple errors for a single field: length limit, blank, etc
type errors map[string][]string

// Add method to add error message for a given field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get method to retrieve the first error message for a given field.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
