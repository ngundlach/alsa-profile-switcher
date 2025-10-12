package pactl

type Card struct {
	Name          string             `json:"name"`
	Driver        string             `json:"driver"`
	Properties    Properties         `json:"properties"`
	Profiles      map[string]Profile `json:"profiles"`
	ActiveProfile string             `json:"active_profile"`
	Ports         map[string]Port    `json:"ports"`
}

type Properties struct {
	DeviceDescription string `json:"device.description"`
	DeviceNick        string `json:"device.nick"`
	DeviceProductName string `json:"device.product.name"`
}
type Profile struct {
	Description string `json:"description"`
	Available   bool   `json:"available"`
}
type Port struct {
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Properties  PortProperties `json:"properties"`
	Profiles    []string       `json:"profiles"`
}
type PortProperties struct {
	DeviceProductName string `json:"device.product.name"`
}
