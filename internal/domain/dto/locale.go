package dto

type IPLocale struct {
	IpAddress        string  `json:"ip_address,omitempty"`
	City             string  `json:"city,omitempty"`
	CityGeonameId    int     `json:"city_geoname_id,omitempty"`
	Region           string  `json:"region,omitempty"`
	RegionIsoCode    string  `json:"region_iso_code,omitempty"`
	RegionGeonameId  int     `json:"region_geoname_id,omitempty"`
	PostalCode       string  `json:"postal_code,omitempty"`
	Country          string  `json:"country,omitempty"`
	CountryCode      string  `json:"country_code,omitempty"`
	CountryGeonameId int     `json:"country_geoname_id,omitempty"`
	Continent        string  `json:"continent,omitempty"`
	ContinentCode    string  `json:"continent_code,omitempty"`
	Longitude        float64 `json:"longitude,omitempty"`
	Latitude         float64 `json:"latitude,omitempty"`
	Security         struct {
		IsVpn bool `json:"is_vpn,omitempty"`
	} `json:"security,omitempty"`
	Connection struct {
		AutonomousSystemNumber       int    `json:"autonomous_system_number,omitempty"`
		AutonomousSystemOrganization string `json:"autonomous_system_organization,omitempty"`
		ConnectionType               string `json:"connection_type,omitempty"`
		IspName                      string `json:"isp_name,omitempty"`
		OrganizationName             string `json:"organization_name,omitempty"`
	} `json:"connection,omitempty"`
}
