package message

import (
	"encoding/binary"
	"fmt"
	"github.com/SJTU-OpenNetwork/go-bitswap/tickets"
	"io"

	pb "github.com/SJTU-OpenNetwork/go-bitswap/message/pb"
	wantlist "github.com/SJTU-OpenNetwork/go-bitswap/wantlist"
	blocks "github.com/ipfs/go-block-format"

	cid "github.com/ipfs/go-cid"
	pool "github.com/libp2p/go-buffer-pool"
	msgio "github.com/libp2p/go-msgio"

	"github.com/libp2p/go-libp2p-core/network"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("Bitswap.message")

// BitSwapMessage is the basic interface for interacting building, encoding,
// and decoding messages sent on the BitSwap protocol.
type BitSwapMessage interface {
	// Wantlist returns a slice of unique keys that represent data wanted by
	// the sender.
	Wantlist() []Entry

	// Blocks returns a slice of unique blocks.
	Blocks() []blocks.Block

	// Ticks returns a slice of tickets.
	Tickets() []tickets.Ticket //TODO: add definitions - Jerry

	TicketAcks() []tickets.TicketAck

	// AddEntry adds an entry to the Wantlist.
	AddEntry(key cid.Cid, priority int)

	Cancel(key cid.Cid)

	Empty() bool

	// A full wantlist is an authoritative copy, a 'non-full' wantlist is a patch-set
	Full() bool

	AddBlock(blocks.Block)

	AddTicket(tickets.Ticket)
	AddTickets([]tickets.Ticket)

	AddTicketAck(tickets.TicketAck)
	AddTicketAcks([]tickets.TicketAck)
	Exportable

	Loggable() map[string]interface{}
	//LogCat()
}

// Exportable is an interface for structures than can be
// encoded in a bitswap protobuf.
type Exportable interface {
	ToProtoV0() *pb.Message
	ToProtoV1() *pb.Message
	ToNetV0(w io.Writer) error
	ToNetV1(w io.Writer) error
}

type impl struct {
	full     bool
	wantlist map[cid.Cid]*Entry
	blocks   map[cid.Cid]blocks.Block
	tickets  map[cid.Cid]tickets.Ticket
	ticketAcks map[cid.Cid]tickets.TicketAck
}

// New returns a new, empty bitswap message
func New(full bool) BitSwapMessage {
	return newMsg(full)
}

func newMsg(full bool) *impl {
	return &impl{
		blocks:   make(map[cid.Cid]blocks.Block),
		wantlist: make(map[cid.Cid]*Entry),
		tickets: make(map[cid.Cid]tickets.Ticket),
		ticketAcks: make(map[cid.Cid]tickets.TicketAck),
		full:     full,
	}
}

// Entry is an wantlist entry in a Bitswap message (along with whether it's an
// add or cancel).
type Entry struct {
	wantlist.Entry
	Cancel bool
}

func newMessageFromProto(pbm pb.Message) (BitSwapMessage, error) {
	m := newMsg(pbm.Wantlist.Full)
	for _, e := range pbm.Wantlist.Entries {
		c, err := cid.Cast([]byte(e.Block))
		if err != nil {
			return nil, fmt.Errorf("incorrectly formatted cid in wantlist: %s", err)
		}
		m.addEntry(c, int(e.Priority), e.Cancel)
	}

	// deprecated
	for _, d := range pbm.Blocks {
		// CIDv0, sha256, protobuf only
		b := blocks.NewBlock(d)
		m.AddBlock(b)
	}
	//

	for _, b := range pbm.GetPayload() {
		pref, err := cid.PrefixFromBytes(b.GetPrefix())
		if err != nil {
			return nil, err
		}

		c, err := pref.Sum(b.GetData())
		if err != nil {
			return nil, err
		}

		blk, err := blocks.NewBlockWithCid(b.GetData(), c)
		if err != nil {
			return nil, err
		}

		m.AddBlock(blk)
	}

	// Generate Tickets and Ticket Acks
	for _, t := range pbm.GetTicketlist() {
		tk, err := tickets.NewBasicTicket(t)
		if(err != nil){
			return nil, err
		}
		m.AddTicket(tk)
	}

	for _, a := range pbm.GetTicketAcklist() {
		ack, err := tickets.NewBasicTicketAck(a)
		if(err != nil){
			return nil, err
		}
		m.AddTicketAck(ack)
	}

	return m, nil
}

func (m *impl) Full() bool {
	return m.full
}

func (m *impl) Empty() bool {
	return len(m.blocks) == 0 && len(m.wantlist) == 0 && len(m.tickets) == 0  && len(m.ticketAcks) == 0
}

func (m *impl) Wantlist() []Entry {
	out := make([]Entry, 0, len(m.wantlist))
	for _, e := range m.wantlist {
		out = append(out, *e)
	}
	return out
}

func (m *impl) Blocks() []blocks.Block {
	bs := make([]blocks.Block, 0, len(m.blocks))
	for _, block := range m.blocks {
		bs = append(bs, block)
	}
	return bs
}

//TODO add function Tickets() - Jerry
func (m *impl) Tickets() []tickets.Ticket{
	ts := make([]tickets.Ticket, 0, len(m.tickets))
	for _, ticket := range m.tickets {
		ts = append(ts, ticket)
	}
	return ts
}

func (m *impl) TicketAcks() []tickets.TicketAck{
	ts := make([]tickets.TicketAck, 0, len(m.ticketAcks))
	for _, ticketack := range m.ticketAcks {
		ts = append(ts, ticketack)
	}
	return ts
}

func (m *impl) Cancel(k cid.Cid) {
	delete(m.wantlist, k)
	m.addEntry(k, 0, true)
}

func (m *impl) AddEntry(k cid.Cid, priority int) {
	m.addEntry(k, priority, false)
}

func (m *impl) addEntry(c cid.Cid, priority int, cancel bool) {
	e, exists := m.wantlist[c]
	if exists {
		e.Priority = priority
		e.Cancel = cancel
	} else {
		m.wantlist[c] = &Entry{
			Entry: wantlist.Entry{
				Cid:      c,
				Priority: priority,
			},
			Cancel: cancel,
		}
	}
}

func (m *impl) AddBlock(b blocks.Block) {
	m.blocks[b.Cid()] = b
}

//TODO add function AddTicket() - Jerry
func (m *impl) AddTicket(t tickets.Ticket){
	m.tickets[t.Cid()] = t;
}
func (m *impl) AddTickets(ts []tickets.Ticket){
	for _, t := range ts{
		m.AddTicket(t)
	}
}

func (m *impl) AddTicketAck(ack tickets.TicketAck){
	m.ticketAcks[ack.Cid()] = ack
}

func (m *impl) AddTicketAcks(acks []tickets.TicketAck){
	for _, ack := range acks{
		m.AddTicketAck(ack)
	}
}

// FromNet generates a new BitswapMessage from incoming data on an io.Reader.
func FromNet(r io.Reader) (BitSwapMessage, error) {
	reader := msgio.NewVarintReaderSize(r, network.MessageSizeMax)
	return FromMsgReader(reader)
}

// FromPBReader generates a new Bitswap message from a gogo-protobuf reader
func FromMsgReader(r msgio.Reader) (BitSwapMessage, error) {
	msg, err := r.ReadMsg()
	if err != nil {
		return nil, err
	}

	var pb pb.Message
	err = pb.Unmarshal(msg)
	r.ReleaseMsg(msg)
	if err != nil {
		return nil, err
	}

	return newMessageFromProto(pb)
}

func (m *impl) ToProtoV0() *pb.Message {
	pbm := new(pb.Message)
	pbm.Wantlist.Entries = make([]pb.Message_Wantlist_Entry, 0, len(m.wantlist))
	for _, e := range m.wantlist {
		pbm.Wantlist.Entries = append(pbm.Wantlist.Entries, pb.Message_Wantlist_Entry{
			Block:    e.Cid.Bytes(),
			Priority: int32(e.Priority),
			Cancel:   e.Cancel,
		})
	}
	pbm.Wantlist.Full = m.full

	blocks := m.Blocks()
	pbm.Blocks = make([][]byte, 0, len(blocks))
	for _, b := range blocks {
		pbm.Blocks = append(pbm.Blocks, b.RawData())
	}

	// Convert ticket and acks
	pbm.Ticketlist = make([]*pb.Ticket, 0, len(m.tickets))
	for _, t := range m.tickets {
		pbm.Ticketlist = append(pbm.Ticketlist, t.ToProto())
	}

	pbm.TicketAcklist = make([]*pb.TicketAck, 0, len(m.ticketAcks))
	for _, a := range m.ticketAcks{
		pbm.TicketAcklist = append(pbm.TicketAcklist, a.ToProto())
	}

	return pbm
}

func (m *impl) ToProtoV1() *pb.Message {
	pbm := new(pb.Message)
	pbm.Wantlist.Entries = make([]pb.Message_Wantlist_Entry, 0, len(m.wantlist))
	for _, e := range m.wantlist {
		pbm.Wantlist.Entries = append(pbm.Wantlist.Entries, pb.Message_Wantlist_Entry{
			Block:    e.Cid.Bytes(),
			Priority: int32(e.Priority),
			Cancel:   e.Cancel,
		})
	}
	pbm.Wantlist.Full = m.full

	blocks := m.Blocks()
	pbm.Payload = make([]pb.Message_Block, 0, len(blocks))
	for _, b := range blocks {
		pbm.Payload = append(pbm.Payload, pb.Message_Block{
			Data:   b.RawData(),
			Prefix: b.Cid().Prefix().Bytes(),
		})
	}

	// Convert ticket and acks
	pbm.Ticketlist = make([]*pb.Ticket, 0, len(m.tickets))
	for _, t := range m.tickets {
		pbm.Ticketlist = append(pbm.Ticketlist, t.ToProto())
	}

	pbm.TicketAcklist = make([]*pb.TicketAck, 0, len(m.ticketAcks))
	for _, a := range m.ticketAcks{
		pbm.TicketAcklist = append(pbm.TicketAcklist, a.ToProto())
	}
	return pbm
}

//TODO add or modify function ToProtoV2()/ToProtoV1() - Jerry

func (m *impl) ToNetV0(w io.Writer) error {
	return write(w, m.ToProtoV0())
}

func (m *impl) ToNetV1(w io.Writer) error {
	return write(w, m.ToProtoV1())
}

//TODO add or modify function ToNetV2()/ToNetV1() - Jerry

func write(w io.Writer, m *pb.Message) error {
	size := m.Size()

	buf := pool.Get(size + binary.MaxVarintLen64)
	defer pool.Put(buf)

	n := binary.PutUvarint(buf, uint64(size))

	written, err := m.MarshalTo(buf[n:])
	if err != nil {
		return err
	}
	n += written

	_, err = w.Write(buf[:n])
	return err
}

func (m *impl) Loggable() map[string]interface{} {
	blocks := make([]string, 0, len(m.blocks))
	for _, v := range m.blocks {
		blocks = append(blocks, v.Cid().String())
	}

	tickets := make([]string, 0, len(m.tickets))
	for _, t := range m.tickets {
		tickets = append(tickets, t.Cid().String())
	}

	ticketacks := make([]map[string]interface{}, 0, len(m.ticketAcks))
	for _, a := range m.ticketAcks {
		ticketacks = append(ticketacks, map[string]interface{}{
			"cid": a.Cid().String(),
			"type": pb.TicketAck_Type_name[a.ACK()],
		})
	}

	return map[string]interface{}{
		"blocks": blocks,
		"wants":  m.Wantlist(),
		"tickets": tickets,
		"ticket acks": ticketacks,
	}
}
