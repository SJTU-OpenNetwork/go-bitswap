// (Overview)

package tickets

import (
	"fmt"
	"time"

	pb "github.com/SJTU-OpenNetwork/go-bitswap/message/pb"
	//"io"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

//type ticketState int
//type ackType int
//type storeType int

/*
type PeerDecodeFail struct{
	pid string
}

func (e *PeerDecodeFail) Error() string{
	return fmt.Sprintf("Cannot decode peer id from string %s", e.pid)
}
*/

const(
	STATE_NEW int32 = iota
	STATE_ACK
    STATE_REJECT

    //unused
	STATE_CANCEL
	STATE_TIMEOUT
)

const(
	ACK_ACCEPT int32 = iota
	ACK_CANCEL
)

const(
	STORE_SEND int32 = iota
	STORE_RECV
)

//var TicketAck_Type_value_name =

type Loggable interface{
	Loggable() map[string]interface{}
}

type Ticket interface {
	Publisher() peer.ID		// Needed because the ticket receiver should send ack to ticket sender
	SendTo() peer.ID
	Cid() cid.Cid
	TimeStamp() int64 // Higher level, lower priority
	//Valid(peer.ID, peer.ID) bool
	GetSize() int64
	GetState() int32
	BasicInfo() string
	//Called when receive ack for this Ticket
	ACKed()
	//Called when receive cancel for thie Ticket
	Canceled()
	//TODO: Add digital signature to ticket. - Riften
	//Signature String
	ToProto() *pb.Ticket

	SetPublisher(peer.ID)
	SetTimeStamp(timeStamp int64)
	SetState(int32)

	Loggable
	//Loggable() map[string]interface{}
}

type TicketAck interface {
	// Publisher, Receiver, Cid should return the same value as Ticket
	// They identify the ack for specific ticket
	Publisher() peer.ID
	Receiver() peer.ID
	Cid() cid.Cid

	// Type of ack
	ACK() int32
	ToProto() *pb.TicketAck

	//Loggable() map[string]interface{}
	Loggable
}

// TODO: Is it possible that the timestamps of received ticket and sent ticket are different? Does that matter? - Riften

type BasicTicket struct{
	publisher peer.ID
	sendTo	peer.ID
	contentId cid.Cid
	timeStamp int64
	blockSize int64
	state int32
}

type BasicTicketAck struct{
	publisher peer.ID
	receiver peer.ID
	cid cid.Cid
	aType int32
}

func NewBasicTicket(ticket *pb.Ticket) (*BasicTicket, error){
	//publisherId, err := peer.IDB58Decode(ticket.GetPublisher())
	//if(err != nil){
	//	return nil, err
	//}
	sendtoId, err := peer.IDB58Decode(ticket.GetSendTo())
	if(err != nil){
		return nil, err
	}
	contentId, err := cid.Decode(ticket.GetCid())
	if(err != nil){
		return nil, err
	}

	tmpState := pb.Ticket_State_value[ticket.GetState().String()]

	//pb.Ticket.GetState().EnumDescriptor()
	return &BasicTicket{
		//publisher: publisherId,
		sendTo:    sendtoId,
		contentId: contentId,
		timeStamp: ticket.GetTimeStamp(),
		blockSize: ticket.GetByteSize(),
		state:     tmpState,
	}, nil
}

func CreateBasicTicket(sendTo peer.ID, contentId cid.Cid, blockSize int64) *BasicTicket{
	return &BasicTicket{
		publisher: peer.ID(""),
		sendTo:    sendTo,
		contentId: contentId,
		timeStamp: 0,
		blockSize: blockSize,
		state:     STATE_NEW,
	}
}

func CreateBasicTicketWithTime(sendTo peer.ID, contentId cid.Cid, blockSize int64, timeStamp int64) *BasicTicket{
	return &BasicTicket{
		publisher: peer.ID(""),
		sendTo: sendTo,
		contentId: contentId,
		blockSize: blockSize,
		timeStamp: timeStamp,
		state: STATE_NEW,
	}
}

func (t *BasicTicket) Publisher() peer.ID{
	return t.publisher
}

func (t *BasicTicket) SetPublisher(p peer.ID) {
	t.publisher = p
}

func (t *BasicTicket) SendTo() peer.ID{
	return t.sendTo
}

func (t *BasicTicket) Cid() cid.Cid{
	return t.contentId
}

func (t* BasicTicket) TimeStamp() int64{
	return t.timeStamp
}

//func (t* BasicTicket) Valid(sender peer.ID, receiver peer.ID) bool{
//	return t.Publisher() == sender && t.SendTo() == receiver
//}

func (t* BasicTicket) GetSize() int64{
	return t.blockSize
}

func (t* BasicTicket) GetState() int32{
	return t.state
}

func (t* BasicTicket) SetState(state int32) {
	t.state = state
}

func (t* BasicTicket) GetStateString() string{
	return pb.Ticket_State_name[t.GetState()]
}

func (t* BasicTicket) ACKed(){
	t.state = STATE_ACK
}

func (t* BasicTicket) Canceled(){
	t.state = STATE_CANCEL
}

func (t *BasicTicket) ToProto() *pb.Ticket{
	//publisher := t.publisher.String()
	receiver := t.sendTo.String()
	cid := t.contentId.String()

	re := &pb.Ticket{
		//Publisher: publisher,
		SendTo:    receiver,
		Cid:       cid,
		ByteSize:  t.blockSize,
		TimeStamp: t.timeStamp,
		State:     pb.Ticket_State(t.state),
	}
	return re
}

func (t *BasicTicket) BasicInfo() string{
	return fmt.Sprintf("To: %s\nCid: %s\nbyteSize: %d\ntimeStamp: %d\nState:%s",
		t.sendTo.Pretty(), t.contentId.String(), t.blockSize, t.timeStamp, pb.Ticket_State_name[t.state])
}

func (t *BasicTicket) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"publisher": t.Publisher().String(),
		"receiver": t.SendTo().String(),
		"cid": t.Cid().String(),
		"timestamp": t.TimeStamp(),
		"state" : pb.Ticket_State_name[t.GetState()],
	}
}


func makeTimestamp(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}


func (a *BasicTicketAck) Publisher() peer.ID{
	return a.publisher
}

func (a *BasicTicketAck) Receiver() peer.ID{
	return a.receiver
}

func (a *BasicTicketAck)Cid() cid.Cid{
	return a.cid
}

func (a *BasicTicketAck) ACK() int32{
	return a.aType
}



func NewBasicTicketAck(ack *pb.TicketAck) (*BasicTicketAck, error){
	publisherId, err := peer.IDB58Decode(ack.GetPublisher())
	if(err != nil){
		return nil, err
	}
	receiverId, err := peer.IDB58Decode(ack.GetReceiver())
	if(err != nil){
		return nil, err
	}
	contentId, err := cid.Decode(ack.GetCid())
	if(err != nil){
		return nil, err
	}
	return &BasicTicketAck{
		publisher: publisherId,
		receiver:  receiverId,
		cid:       contentId,
		aType:     pb.Ticket_State_value[ack.GetType().String()],
	}, nil
}

func CreateBasicTicketAck(publisher peer.ID, receiver peer.ID, contentId cid.Cid, acktype int32) *BasicTicketAck{
	return &BasicTicketAck{
		publisher: publisher,
		receiver:  receiver,
		cid:       contentId,
		aType:     acktype,
	}
}

func (a *BasicTicketAck) ToProto() *pb.TicketAck{
	publisher := a.publisher.String()
	receiver := a.receiver.String()
	cid := a.cid.String()
	return &pb.TicketAck{
		Publisher: publisher,
		Receiver:  receiver,
		Cid:       cid,
		Type:      pb.TicketAck_Type(a.aType),
	}
}

func (a *BasicTicketAck) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"publisher": a.Publisher().String(),
		"receiver": a.Receiver().String(),
		"cid": a.Cid().String(),
		"type": pb.TicketAck_Type_name[a.ACK()],
	}
}

func AckFromTicket(p peer.ID, t Ticket, ackType int32)TicketAck{
	return CreateBasicTicketAck(p, t.SendTo(), t.Cid(), ackType)
}

func GetRejectAcks(p peer.ID, ts []Ticket) []TicketAck{
	acks := make([]TicketAck, 0)
	for _, t :=range ts{
		acks = append(acks, AckFromTicket(p, t, ACK_CANCEL))
	}
	return acks
}

func GetAcceptAcks(p peer.ID, ts []Ticket) []TicketAck{
	acks := make([]TicketAck, 0)
	for _, t :=range ts{
		acks = append(acks, AckFromTicket(p, t, ACK_ACCEPT))
	}
	return acks
}

func (a *BasicTicket) SetTimeStamp(timeStamp int64){
	a.timeStamp = timeStamp
}
