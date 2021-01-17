package model

// Peer structure
type HostConfig struct {
	MeshName string
	MeshId   string
	Hosts    []Host
}

// Host structure
type Message struct {
	Id     string       `json:"id"       bson:"id"`
	Config []HostConfig `json:"config"   bson:"config"`
}
