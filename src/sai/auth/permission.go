package auth

type Param struct {
	Name  string   `json:"name"`
	Rules []string `json:"rules"`
	//Required bool     `json:"required"`
}

type Rights struct {
	Read  int `json:"read"`
	Write int `json:"write"`
}

type Permission struct {
	URL    string  `json:"url"`
	Rights Rights  `json:"rights"`
	Params []Param `json:"params"`
}
