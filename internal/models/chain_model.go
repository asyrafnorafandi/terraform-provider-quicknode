package models

type ChainModel struct {
	Slug     string         `json:"slug"`
	Networks []NetworkModel `json:"networks"`
}

type NetworkModel struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
