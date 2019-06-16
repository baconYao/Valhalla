package game

import (
	"log"

	"github.com/Hucaru/Valhalla/game/entity"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

func (server *Channel) playerChangeChannel(conn mnet.Client, reader mpacket.Reader) {
	id := reader.ReadByte()

	server.migrating[conn] = id
	server.sessions[conn].Save(server.db)

	if int(id) < len(server.channels) {
		if server.channels[id].Port == 0 {
			conn.Send(entity.PacketMessageDialogueBox("Cannot change channel"))
		} else {
			conn.Send(entity.PacketChangeChannel(server.channels[id].IP, server.channels[id].Port))
		}
	}
}

func (server *Channel) playerConnect(conn mnet.Client, reader mpacket.Reader) {
	charID := reader.ReadInt32()

	var accountID int32
	err := server.db.QueryRow("SELECT accountID FROM characters WHERE id=?", charID).Scan(&accountID)

	if err != nil {
		log.Println(err)
	}

	conn.SetAccountID(accountID)

	// check migration

	char := entity.Character{}
	char.LoadFromID(server.db, charID)

	var adminLevel int
	err = server.db.QueryRow("SELECT adminLevel FROM accounts WHERE accountID=?", conn.GetAccountID()).Scan(&adminLevel)

	if err != nil {
		log.Println(err)
	}

	conn.SetAdminLevel(adminLevel)

	server.sessions[conn] = &char

	conn.Send(entity.PacketPlayerEnterGame(char, int32(server.id)))
	conn.Send(entity.PacketMessageScrollingHeader("Valhalla Archival Project"))
}