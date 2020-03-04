package kube

import "errors"

// mergeMaps returns the union of two maps
func mergeMaps(m1, m2 map[string]string) (map[string]string, error) {
	for k, v := range m1 {
		if _, ok := m2[k]; ok {
			return m2, errors.New("duplicate key found")
		}
		m2[k] = v
	}
	return m2, nil
}

// removeEmpty removes empty entries from a map
func removeEmpty(m1 map[string]string) map[string]string {
	for k, v := range m1 {
		if v == "" {
			delete(m1, k)
		}
	}
	return m1
}
