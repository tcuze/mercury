package registry

type Service struct {
	Namespace string
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Nodes     []*Node           `json:"nodes"`
	Metadata  map[string]string `json:"metadata"`
}

type Node struct {
	Id       string            `json:"id"`
	Host     string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
}
