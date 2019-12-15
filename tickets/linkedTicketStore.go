package tickets

import (
	"container/list"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	logging "github.com/ipfs/go-log"
	"sync"
)

var log = logging.Logger("hon.linkedTicketStore")

type NotInitializeError struct{}
func (e *NotInitializeError) Error() string{
	return "linkedTicketStore not initialized!"
}

type InvalidPublisherError struct{
	creater peer.ID
	publisher peer.ID
}
func (e *InvalidPublisherError) Error() string{
	return fmt.Sprintf("Invalid ticket.\n" +
		"Ticket Store Creater: %s\n" +
		"Ticket Publisher: %s", e.creater.String(), e.publisher.String());
}

// linked Ticket Store implement the ticketStore interface.
// It uses linked list as the based data structure and further builds tracker on it.
// Add, Remove, Select, Modify will be done within O(1) time complexity
// TODO: Make linkedTicketStore support both send and recv ticket
type linkedTicketStore struct{
	mutex sync.Mutex
	creater peer.ID
	dataStore map[cid.Cid] *list.List
	dataTracker map[cid.Cid]map[peer.ID] *list.Element
	storeType int32
}

func NewLinkedSendTicketStore(creater peer.ID) *linkedTicketStore{
	return &linkedTicketStore{
		creater: creater,
		dataStore: make(map[cid.Cid]*list.List),
		dataTracker: make(map[cid.Cid]map[peer.ID]*list.Element),
		storeType: STORE_SEND,
	}
}

func NewLinkedRecvTicketStore(creater peer.ID) *linkedTicketStore{
	return &linkedTicketStore{
		creater: creater,
		dataStore: make(map[cid.Cid]*list.List),
		dataTracker: make(map[cid.Cid]map[peer.ID]*list.Element),
		storeType: STORE_RECV,
	}
}

func (s *linkedTicketStore) varify(ticket Ticket) error{
	if(s.creater != ticket.Publisher()){
		return &InvalidPublisherError{
			creater:   s.creater,
			publisher: ticket.Publisher(),
		}
	}

	return nil
}

func (s *linkedTicketStore) AddTicket(ticket Ticket) error{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// TicketStore only store the ticket published from self
	err := s.varify(ticket)
	if(err != nil){
		return err
	}

	// Judge whether this ticket is already exists
	// Remove the old one if so
	ok := s.TicketExists(ticket.SendTo(), ticket.Cid())
	if(ok){
		log.Warningf("Ticket of %s sent to %s already exists.\n" +
			"It would be replaced by the new one.\n" +
			"It is better to remove the old one manually before add the new one", ticket.Cid().String(), ticket.SendTo().String())
		s.RemoveTicket(ticket.SendTo(), ticket.Cid())
	}

	// Add ticket to dataStore
	tmplist, ok := s.dataStore[ticket.Cid()]
	var tmpElm *list.Element
	if(ok){
		tmpElm = tmplist.PushBack(ticket)
	}else{
		s.dataStore[ticket.Cid()] = list.New()
		tmpElm = s.dataStore[ticket.Cid()].PushBack(ticket)
	}

	// Add element to dataTracker
	tmpmap, ok := s.dataTracker[ticket.Cid()]
	if(ok){//Sub map already created
		tmpmap[ticket.SendTo()] = tmpElm
	}else{//Create sub map first
		s.dataTracker[ticket.Cid()] = make(map[peer.ID]*list.Element)
		s.dataTracker[ticket.Cid()][ticket.SendTo()] = tmpElm
	}

	return nil
}

func (s *linkedTicketStore) TicketExists(pid peer.ID, cid cid.Cid) bool {
	tmpmap, ok := s.dataTracker[cid]
	if(!ok){
		return false
	}else{
		_, ok := tmpmap[pid]
		return ok
	}
}

func (s *linkedTicketStore) GetTickets(cid cid.Cid) ([]Ticket, error){
	return nil, nil
}

func (s *linkedTicketStore) RemoveTicket(pid peer.ID, cid cid.Cid) error{
	return nil
}

func (s *linkedTicketStore) RemoveTicketEqualsTo(ticket Ticket){

}

func (s *linkedTicketStore) PopTickets() *TicketTask{
	return nil
}

func (s *linkedTicketStore) Clean(){

}

func (s *linkedTicketStore) RemoveCanceled() int {
	return 0
}

func (s *linkedTicketStore) TicketNumber() int{
	return 0
}

func (s *linkedTicketStore) TicketSize() int64 {
	return 0
}

func RemoveTickets(pid peer.ID, cid []cid.Cid) error {
    return nil
}

func PrepareSending(acks []TicketAck) error {
    return nil
}

func RemoveSendingTask(pid peer.ID, cid []cid.Cid) error {
    return nil
}

func PopSendingTasks(cid []cid.Cid) ([]TicketAck, error) {
    return nil, nil
}
