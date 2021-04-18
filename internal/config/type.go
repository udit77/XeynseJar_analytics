package config

type Config struct {
	Server   Server   `json:"server"`
	DataBase DataBase `json:"db"`
	MQTT     MQTT     `json:"mqtt"`
	Certs    Certs    `json:"certs"`
}

type Server struct {
	Address string `json:"address"`
}

type DataBase struct {
	XeynsJar          XeynseJar           `json:"xeynse_jar"`
	HomeConfiguration HomeConfigurationDB `json:"homeConfiguration"`
	Products          ProductsDB          `json:"products"`
}

type HomeConfigurationDB struct {
	Table string `json:"table"`
}

type XeynseJar struct {
	JarStatus JarStatus `json:"jar_status"`
}

type ProductsDB struct {
	Table string `json:"table"`
}

type JarStatus struct {
	Table    string `json:"table"`
	Host     string `json:"host"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type MQTT struct {
	Broker   string `json:"broker"`
	ClientID string `json:"clientID"`
}

type Certs struct {
	Root    string `json:"root"`
	Cert    string `json:"cert"`
	Private string `json:"private"`
}
