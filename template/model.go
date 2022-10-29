package template

type Model map[string]interface{}

func (m Model) AddAttribute(name string, attribute any) {
	m[name] = attribute
}

func (m Model) GetAttribute(name string) any {
	return m[name]
}

func (m Model) ContainsAttribute(name string) bool {
	_, ok := m[name]
	return ok
}

func (m Model) AddAllAttributes(attributes map[string]any) {
	if attributes == nil {
		return
	}
	for key, value := range attributes {
		m[key] = value
	}
}

// MergeAttributes
// Copy all attributes in the supplied {@code Map} into this {@code Map},
// with existing objects of the same name taking precedence (i.e. not getting replaced).
func (m Model) MergeAttributes(attributes map[string]any) {
	for name, value := range attributes {
		if !m.ContainsAttribute(name) {
			m.AddAttribute(name, value)
		}
	}
}
