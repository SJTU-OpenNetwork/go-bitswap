// (Overview)

package tickets

import (
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type Ticket interface {
	Pid() peer.ID
	Cid() cid.Cid
	TimeStamp() int64
	//TODO: Add digital signature to ticket. - Riften
	//Signature String
}

type BasicTicket struct{
	peerId peer.ID
	contentId cid.Cid
	timeStamp int64
}

func (t *BasicTicket) Pid() peer.ID{
	return t.peerId
}

func (t *BasicTicket) Cid() cid.Cid{
	return t.contentId
}

func (t* BasicTicket) TimeStamp() int64{
	return t.timeStamp
}
