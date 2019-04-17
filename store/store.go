package store

import (
	"mercury/packet"
)

type Store interface {
	Commit(*packet.WebsocketContext) error
	Rebuild() error
	Start()
	Receive(*packet.WebsocketContext) error
}
