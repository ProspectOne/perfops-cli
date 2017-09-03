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
		ID         int        `json:"id"`
		Name       string     `json:"name"`
		ISO        string     `json:"iso"`
		ISONumeric string     `json:"iso_numeric"`
		Continent  *Continent `json:"continent,omitempty"`
	}

	// Node contains informatin about a test node.
	Node struct {
		ID        int      `json:"id"`
		Latitude  float64  `json:"latitude"`
		Longitude float64  `json:"longitude"`
		City      string   `json:"city"`
		SubRegion string   `json:"sub_region"`
		Country   *Country `json:"country,omitempty"`
	}
)
