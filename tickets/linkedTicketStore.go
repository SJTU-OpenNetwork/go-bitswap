package tickets

import (
	"container/list"
	"fmt"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	"sync"
)

var log = logging.Logger("hon.linkedTicketStore")
var SendingListTypeError = fmt.Errorf("Elements in prapareSendingList is not TicketAck")

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

type TicketNotFound struct{
	pid peer.ID
	cid cid.Cid
}
func (e *TicketNotFound) Error() string{
	return fmt.Sprintf("Ticket not found.\n" +
		"pid: %s\n" +
		"cid: %s", e.pid, e.cid)
}

type PeerHandler interface {
	SendTicketMessage(entries []Ticket, targets []peer.ID)
	SendTicketAckMessage(entries []TicketAck, targets []peer.ID)
}

// linked Ticket Store implement the ticketStore interface.
// It uses linked list as the based data structure and further builds tracker on it.
// Add, Remove, Select, Modify will be done within O(1) time complexity
// TODO: Make linkedTicketStore support both send and recv ticket
type linkedTicketStore struct{
	mutex sync.Mutex
	//creater peer.ID
	storeNumber int
	storeSize int64

	dataStore map[cid.Cid] *list.List
	dataTracker map[cid.Cid]map[peer.ID] *list.Element
	storeType int32

    prepareSendingList *list.List
    receivedTickets map[cid.Cid] Ticket
    peerManager PeerHandler
}

func NewLinkedTicketStore(pm PeerHandler) *linkedTicketStore{
	return &linkedTicketStore{
		dataStore: make(map[cid.Cid]*list.List),
		dataTracker: make(map[cid.Cid]map[peer.ID]*list.Element),
		storeType: STORE_SEND,

        prepareSendingList: list.New(),
        receivedTickets: make(map[cid.Cid] Ticket),
        peerManager: pm,
	}
}

// deprecated
//func NewLinkedTicketStore(creater peer.ID) *linkedTicketStore{
//	return &linkedTicketStore{
//		creater: creater,
//		dataStore: make(map[cid.Cid]*list.List),
//		dataTracker: make(map[cid.Cid]map[peer.ID]*list.Element),
//		storeType: STORE_SEND,
//
//        prepareSendingList: list.New(),
//        receivedTickets: make(map[cid.Cid] Ticket),
//	}
//}

//func NewLinkedRecvTicketStore(creater peer.ID) *linkedTicketStore{
//	return &linkedTicketStore{
//		creater: creater,
//		dataStore: make(map[cid.Cid]*list.List),
//		dataTracker: make(map[cid.Cid]map[peer.ID]*list.Element),
//		storeType: STORE_RECV,
//	}
//}



func (s *linkedTicketStore) AddTicket(ticket Ticket) error{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// TicketStore only store the ticket published from self
	//err := s.varify(ticket)
	//if(err != nil){
	//	return err
	//}


	// TODO: What if we already have the ticket??? - Riften
	// Judge whether this ticket is already exists
	// Remove the old one if so
	//ok := s.TicketExists(ticket.SendTo(), ticket.Cid())
	//if(ok){
	//	log.Warningf("Ticket of %s sent to %s already exists.\n" +
	//		"We would not add a new one or replace the old one.\n" +
	//		"It is better to remove the old one manually before add the new one", ticket.Cid().String(), ticket.SendTo().String())
		// Remove Here will cause a dead lock
		//s.RemoveTicket(ticket.SendTo(), ticket.Cid())
	//}

	// Add ticket to dataStore
	tmplist, ok := s.dataStore[ticket.Cid()]
	var tmpElm *list.Element
	if(ok){
		//_, ok = tmplist[]
		tmpElm = tmplist.PushBack(ticket)
	}else{
		s.dataStore[ticket.Cid()] = list.New()
		tmpElm = s.dataStore[ticket.Cid()].PushBack(ticket)
	}

	// Add element to dataTracker
	tmpmap, ok := s.dataTracker[ticket.Cid()]
	if(ok){//Sub map already created
		//TODO: Fix redundant add if it happends - Riften
		if(tmpmap[ticket.SendTo()]!=nil){
			log.Error("Redundant add a ticket anready exists!\n This may cause inconsistency with store and tracker!")
		}
		tmpmap[ticket.SendTo()] = tmpElm
	}else{//Create sub map first
		s.dataTracker[ticket.Cid()] = make(map[peer.ID]*list.Element)
		s.dataTracker[ticket.Cid()][ticket.SendTo()] = tmpElm
	}

	s.storeNumber ++
	s.storeSize += ticket.GetSize()

	return nil
}

func (s *linkedTicketStore) AddTickets(ts []Ticket) error {
	for _, t := range ts{
		s.AddTicket(t)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Remove tracker:
	tmpmap, ok := s.dataTracker[cid]
	var elem *list.Element
	if(!ok){
        return nil
//		return &TicketNotFound{
//			pid: pid,
//			cid: cid,
//		}
	}else{
		elem, ok = tmpmap[pid]
		if(!ok){
            return nil
//			return &TicketNotFound{
//				pid: pid,
//				cid: cid,
//			}
		}
	}

	delete(tmpmap, pid)
	if(len(tmpmap) <= 0){
		//tmpmap = nil
		delete(s.dataTracker, cid)
	}

	// Remove data
	tmpTicket := s.dataStore[cid].Remove(elem).(Ticket)

	if(s.dataStore[cid].Len() == 0){
		delete(s.dataStore, cid)
	}

	s.storeNumber --
	s.storeSize -= tmpTicket.GetSize()


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
	return s.storeNumber
}

func (s *linkedTicketStore) TicketSize() int64 {
	return s.storeSize
}

func (s *linkedTicketStore) RemoveTickets(pid peer.ID, cid []cid.Cid) error {
    for _, c := range cid{
    	err := s.RemoveTicket(pid, c)
    	if err != nil {
    		return err
		}
	}
	return nil
}

func (s *linkedTicketStore) PrepareSending(acks []TicketAck) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    for _, ack := range acks {
        s.prepareSendingList.PushBack(ack)
    }
    return nil
}

func (s *linkedTicketStore) RemoveSendingTasks(pid peer.ID, cids []cid.Cid) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    cur := s.prepareSendingList.Front()
    for cur != nil {
        ack, ok := cur.Value.(TicketAck)
        if !ok {
            return SendingListTypeError
        }
        if ack.Receiver() == pid {
            for _, cid := range cids {
                if ack.Cid() == cid {
                    tmp := cur
                    cur = cur.Next()
                    s.prepareSendingList.Remove(tmp)
                    goto NEXT
                }
            }
        }
        cur = cur.Next()
        NEXT:
    }
    return nil
}

func (s *linkedTicketStore) PopSendingTasks(cids []cid.Cid) ([]TicketAck, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    var tasks []TicketAck
    cur := s.prepareSendingList.Front()
    for cur != nil {
        ack, ok := cur.Value.(TicketAck)
        if !ok {
            return tasks, SendingListTypeError
        }
        for _, cid := range cids {
            if ack.Cid() == cid {
                tmp := cur
                cur = cur.Next()
                tasks = append(tasks, ack)
                s.prepareSendingList.Remove(tmp)
                goto NEXT
            }
        }
        cur = cur.Next()
        NEXT:
    }

    return tasks, nil
}

func (s *linkedTicketStore) StoreReceivedTickets(tickets []Ticket) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    for _, ticket := range tickets {
        s.receivedTickets[ticket.Cid()] = ticket
    }
    return nil
}

func (s *linkedTicketStore) GetReceivedTicket(cids []cid.Cid) (map[cid.Cid] Ticket, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ticketMap := make(map[cid.Cid] Ticket)
    for _, cid := range cids {
        ticket, ok := s.receivedTickets[cid]
        if ok {
            ticketMap[cid] = ticket
        }
    }
    return ticketMap, nil
}

func (s *linkedTicketStore) PredictTime() int64{
	//TODO: Not implemented for now - Riften
	//return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
	return int64(s.storeNumber) * 100
}

func (s *linkedTicketStore) Loggable() map[string] interface{}{
	storeLoggable :=  make(map[string][]string) // {"pid": [cids]}
	//trackerLoggable := make(map[string]interface{})

	for cid, m := range(s.dataTracker){
		//tmplist := make([]string, 0, len(m))
		for pid, _ := range m{
			tmplist, ok := storeLoggable[pid.String()]
			if !ok {
				tmplist = make([]string,0)
				tmplist = append(tmplist, cid.String())
				storeLoggable[pid.String()] = tmplist
			} else {
				storeLoggable[pid.String()] = append(storeLoggable[pid.String()], cid.String())
			}
		}
	}
	return map[string]interface{}{
		"datastore": storeLoggable,
		//"dataTracker": trackerLoggable,
		"storeNumber": s.storeNumber,
		"storeSize": s.storeSize,
	}
}

func (s *linkedTicketStore) LoggableFull() map[string] interface{}{
	storeLoggable :=  make(map[string]interface{})
	trackerLoggable := make(map[string]interface{})

	for cid, l := range(s.dataStore){
		//storeLoggable[cid.String()] = make([]map[string]string, 0, l.Len())
		tmplist := make([]string, 0, l.Len())
		for t:=l.Front(); t!=nil; t=t.Next(){
			tmplist = append(tmplist, t.Value.(Ticket).SendTo().String())
		}
		storeLoggable[cid.String()] = tmplist
	}

	for cid, m := range(s.dataTracker){
		tmplist := make([]string, 0, len(m))
		for k, _ := range m{
			tmplist = append(tmplist, k.String())
		}
		trackerLoggable[cid.String()] = tmplist
	}
	return map[string]interface{}{
		"datastore": storeLoggable,
		"dataTracker": trackerLoggable,
		"storeNumber": s.storeNumber,
		"storeSize": s.storeSize,
	}
}
    
func (s *linkedTicketStore) SendTickets(tickets []Ticket, pid peer.ID) {
    var pids = []peer.ID{pid}
    s.peerManager.SendTicketMessage(tickets, pids)
}

func (s *linkedTicketStore) SendTicketAcks(acks []TicketAck, pid peer.ID) {
    var pids = []peer.ID{pid}
    s.peerManager.SendTicketAckMessage(acks, pids)
}
