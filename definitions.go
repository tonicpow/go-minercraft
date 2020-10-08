package minercraft

// Miner is a configuration per miner, including connection url, auth token, etc
type Miner struct {
	MinerID string `json:"miner_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Token   string `json:"token,omitempty"`
	URL     string `json:"url"`
}

// UpdateToken will update an auth token for a given miner
func (m *Miner) UpdateToken(token string) {
	m.Token = token
}
