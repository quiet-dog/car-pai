package response

type Select struct {
	Label    string    `json:"label"`
	Value    uint      `json:"value"`
	Type     string    `json:"type"`
	Children []*Select `json:"children"`
	Disabled bool      `json:"disabled"`
}
