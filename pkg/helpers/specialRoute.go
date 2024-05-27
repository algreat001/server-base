package helpers

import (
	"strings"
)

type SpecialRoute struct {
	routes []string
}

func NewSpecialRoute() *SpecialRoute {
	return &SpecialRoute{
		routes: make([]string, 0),
	}
}

func (a *SpecialRoute) Add(path string) {
	a.routes = append(a.routes, strings.ToLower(path))
	return
}

func (a *SpecialRoute) IsContained(path string) bool {
	lPath := strings.ToLower(path)
	if (strings.Index(lPath, "/api") != 0) && (strings.Index(lPath, "api") != 0) {
		return true
	}
	for _, route := range a.routes {
		if route == lPath || (len(lPath) > len(route) && route == lPath[:len(lPath)-1]) {
			return true
		}
	}
	return false
}
