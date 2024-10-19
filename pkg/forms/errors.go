package forms

// here using []string slice as the value of the map,
// confusing
type errors map[string][]string

// Add error to the map, use form field as the key
// wondering if happen there are two same name key?
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	// here, is a []string, not string
	es := e[field]

	if len(es) == 0 {
		return ""
	}

	// of course, need return the first value, no idea why doing this.
	return es[0]

}
