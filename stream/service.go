package stream

type Service struct {
	Namespace string
	Name      string
	Version   string
	Metadata  map[string]string
	Host      string
	Port      int
}
