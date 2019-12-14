package tickets

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

type TicketTask struct{
	Target peer.ID
	Tickets []Ticket
}

/**
 * ticketStore is used to store the tickets published.
 * Please do not use it to store the ticket received!!!
 * It is used in following cases:
 *		- Compute the size of tickets waiting for deliving
 * 		- Record the tickets so that we can verify the ticket ack.
 * 		- Control how many tickets to publish
 */
type TicketStore interface{
	AddTicket(ticket Ticket) error
	GetTickets(cid cid.Cid) ([]Ticket, error)
	RemoveTicket(pid peer.ID, cid cid.Cid) error
	RemoveTicketEqualsTo(ticket Ticket)
	Clean()
	RemoveCanceled() int
	//Size() int

	PopTickets() *TicketTask

	TicketNumber() int
	TicketSize() int64
	//StoreType()
}
