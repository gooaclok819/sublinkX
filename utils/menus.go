package utils

import "encoding/json"

type Meta struct {
	Title     string   `json:"title"`
	Icon      string   `json:"icon"`
	Hidden    bool     `json:"hidden"`
	Roles     []string `json:"roles"`
	KeepAlive bool     `json:"keepAlive,omitempty"`
}

type Child struct {
	Path      string `json:"path"`
	Component string `json:"component"`
	Name      string `json:"name"`
	Meta      Meta   `json:"meta"`
}

type Menu struct {
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Redirect  string  `json:"redirect"`
	Name      string  `json:"name"`
	Meta      Meta    `json:"meta"`
	Children  []Child `json:"children"`
}

func GetMenus() []byte {
	menus := []Menu{
		{
			Path:      "/system",
			Component: "Layout",
			Redirect:  "/system/user",
			Name:      "/system",
			Meta: Meta{
				Title:  "系统管理",
				Icon:   "system",
				Hidden: false,
				Roles:  []string{"admin"},
			},
			Children: []Child{
				{
					Path:      "user",
					Component: "system/user/index",
					Name:      "User",
					Meta: Meta{
						Title:     "用户管理",
						Icon:      "user",
						Hidden:    false,
						Roles:     []string{"admin"},
						KeepAlive: true,
					},
				},
			},
		},
	}
	jsonMenus, _ := json.Marshal(menus)
	return jsonMenus
}
