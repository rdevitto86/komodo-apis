package address

type Address struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`         // State/Province/Region
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`       // ISO country name or code (e.g., "US")
}