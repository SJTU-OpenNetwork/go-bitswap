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
	AddTickets(ticket []Ticket) error
    RemoveTickets(pid peer.ID, cids []cid.Cid) error // remove a set of tickets from the `sended ticket list`
    AlreadySent(pid peer.ID, cid cid.Cid) bool

    PrepareSending(acks []TicketAck) error // called if don't have corresponding block when receiving an ACK, put the entry on a list
	RemoveSendingTasks(pid peer.ID, cids []cid.Cid) error // remove a specific task from the `prepared sending task list`, called after receive a reject
	PopSendingTasks(cids []cid.Cid) ([]TicketAck, error) // pop all tasks for a specific cid in `prepared sending task list`, called when a block is received
	PredictTime() int64 // Get the predicted time used to complete all the tickets sent. Return time in millsecond

    // Interfaces for receving tickets - Add by Jerry
    StoreReceivedTickets(tickets []Ticket) error
    GetReceivedTicket(cids []cid.Cid) (map[cid.Cid] Ticket, error)

    // Interfaces for sending tickets and ticketAcks - Add by Jerry
    SendTickets(tickets []Ticket, pid peer.ID)
    SendTicketAcks(acks []TicketAck, pid peer.ID)

    Loggable

    //Unused
	//RemoveTicketEqualsTo(ticket Ticket)
	//RemoveCanceled() int
	//PopTickets() *TicketTask
	//TicketNumber() int
	//TicketSize() int64
	//GetTicketsByCids(cids []cid.Cid) (map[cid.Cid] []Ticket, error)
	//StoreType()
	//Clean()
	RemoveTicket(pid peer.ID, cid cid.Cid) error // remove a specific ticket from the `sended ticket list`, called after receive a reject
	//GetTicket(cid cid.Cid, pid peer.ID) (Ticket, error)
	//GetTickets(cid cid.Cid) ([]Ticket, error)
}
