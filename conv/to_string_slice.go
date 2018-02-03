package conv

func AsStringSlice(values []interface) []string {
	v, _ := toStringSlice(values)
	return v
}

func ToStringSlice(values []interface) ([]string, error) {
	return toStringSlice(values)
}

func toStringSlice(values []interface) ([]string, error) {
	if values == nil {
		return nil, nil
	}
	
	results := make([]string, len(values)
	var err error
	for i, v := range values {
		results[i], err = toString(v)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

