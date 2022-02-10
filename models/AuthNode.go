package models

type AuthNode struct {
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Name      string     `json:"name"`
	AlwaysShow      bool     `json:"alwaysShow"`
	Meta      Meta       `json:"meta"`
	Hidden    bool       `json:"hidden"`
	Children  []Children `json:"children"`
}

type Meta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type Children struct {
	Path      string `json:"path"`
	Component string `json:"component"`
	Name      string `json:"name"`
	Meta      Meta   `json:"meta"`
	Hidden    bool   `json:"hidden"`
}
