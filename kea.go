package main

type KeaDHCP4SubnetConfig struct {
	Subnet       string                            `json:"subnet"`
	ID           uint32                            `json:"id,omitempty"`
	Pools        []KeaDHCP4SubnetPoolConfig        `json:"pools"`
	OptionData   []KeaDHCP4SubnetOptionConfig      `json:"option-data,omitempty"`
	Reservations []KeaDHCP4SubnetReservationConfig `json:"reservations,omitempty"`
}

type KeaDHCP4SubnetPoolConfig struct {
	Pool string `json:"pool"` // Range Ex: 192.0.0.10 - 192.0.0.100
}

type KeaDHCP4SubnetOptionConfig struct {
	Name      KeaDHCP4Option `json:"name,omitempty"`
	Code      int            `json:"code,omitempty"`
	Space     string         `json:"space,omitempty"`
	CSVFormat bool           `json:"csv-format,omitempty"`
	Data      string         `json:"data,omitempty"`
}

type KeaDHCP4SubnetReservationConfig struct {
	HWAddress string `json:"hw-address,omitempty"`
	DUID      string `json:"duid,omitempty"`
	CircuitID string `json:"circuit-id,omitempty"`
	IPAddress string `json:"ip-address"`
}

type KeaDHCP4Option string

const (
	TIME_OFFSET            KeaDHCP4Option = "time-offset"            // DHCP Code 2
	ROUTERS                KeaDHCP4Option = "routers"                // DHCP Code 3
	TIME_SERVERS           KeaDHCP4Option = "time-servers"           // DHCP Code 4
	NAME_SERVERS           KeaDHCP4Option = "name-servers"           // DHCP Code 5
	DOMAIN_NAME_SERVERS    KeaDHCP4Option = "domain-name-servers"    // DHCP Code 6
	LOG_SERVERS            KeaDHCP4Option = "log-servers"            // DHCP Code 7
	COOKIE_SERVERS         KeaDHCP4Option = "cookie-servers"         // DHCP Code 8
	LPR_SERVERS            KeaDHCP4Option = "lpr-servers"            // DHCP Code 9
	IMPRESS_SERVERS        KeaDHCP4Option = "impress-servers"        // DHCP Code 10
	DOMAIN_NAME            KeaDHCP4Option = "domain-name"            // DHCP Code 15
	MAX_DGRAM_REASSEMBLY   KeaDHCP4Option = "max-dgram-reassembly"   // DHCP Code 22
	DEFAULT_IP_TTL         KeaDHCP4Option = "default-ip-ttl"         // DHCP Code 23
	PATH_MTU_AGING_TIMEOUT KeaDHCP4Option = "path-mtu-aging-timeout" // DHCP Code 24
	PATH_MTU_PLATEAU_TABLE KeaDHCP4Option = "path-mtu-plateau-table" // DHCP Code 25
	INTERFACE_MTU          KeaDHCP4Option = "interface-mtu"          // DHCP Code 26
	NIS_DOMAIN             KeaDHCP4Option = "nis-domain"             // DHCP Code 40
	NIS_SERVERS            KeaDHCP4Option = "nis-servers"            // DHCP Code 41
	NTP_SERVERS            KeaDHCP4Option = "ntp-servers"            // DHCP Code 42
)
