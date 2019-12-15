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
	AddTicket(ticket Ticket) error // called after sending a ticket
	GetTickets(cid cid.Cid) ([]Ticket, error)
	RemoveTicket(pid peer.ID, cid cid.Cid) error // remove a specific ticket from the `sended ticket list`, called after receive a reject
	RemoveTicketEqualsTo(ticket Ticket)
	Clean()
	RemoveCanceled() int
	//Size() int

	PopTickets() *TicketTask

	TicketNumber() int
	TicketSize() int64
	//StoreType()

	// Add by Jerry
    RemoveTickets(pid peer.ID, cid []cid.Cid) error // remove a set of tickets from the `sended ticket list`
	PrepareSending(acks []TicketAck) error // called if don't have corresponding block when receiving an ACK, put the entry on a list
	RemoveSendingTask(pid peer.ID, cid []cid.Cid) error // remove a specific task from the `prepared sending task list`, called after receive a reject
	PopSendingTasks(cid []cid.Cid) ([]TicketAck, error) // pop all tasks for a specific cid in `prepared sending task list`, called when a block is received
}
