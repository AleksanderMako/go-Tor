package libinterface

type OnionLib interface {
	CreateOnionChain([]string) (string, error)
	HandshakeWithPeers(string) error
	GenerateSymetricKeys(string) error
	BuildP2PCircuit([]byte, string) error
	SendMessage([]byte, string) error
}
