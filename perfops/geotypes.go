package perfops

type (
	// Continent contains information about a continent.
	Continent struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		ISO  string `json:"iso"`
	}

	// Country contains information about a country.
	Country struct {
		ID         int         `json:"id"`
		Name       string      `json:"name"`
		ISO        string      `json:"iso"`
		ISONumeric interface{} `json:"isoNumeric"` // TODO: Must fix API to return same type for tests and for list endpoints
		Continent  *Continent  `json:"continent,omitempty"`
	}

	// City contains information about a city.
	City struct {
		Name    string `json:"name"`
		Country *struct {
			Name string `json:"name"`
		} `json:"country,omitempty"`
		Continent *struct {
			Name string `json:"name"`
		} `json:"continent,omitempty"`
	}
)
