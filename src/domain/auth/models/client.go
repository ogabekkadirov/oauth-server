package models

type Client struct {
	ID           	string `json:"id" db:"id"`
	Secret       	string `json:"secret" db:"secret"`
	RedirectURIs 	[]string `json:"redirect_uris" db:"redirect_uris"`
	GrantTypes   	[]string `json:"grant_types" db:"grant_types"`
	Scopes        	[]string `json:"scopes" db:"scopes"`
}
