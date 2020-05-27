// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package bitswap_message_pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Ticket_State int32

const (
	Ticket_New     Ticket_State = 0
	Ticket_ACK     Ticket_State = 1
	Ticket_CANCEL  Ticket_State = 2
	Ticket_TIMEOUT Ticket_State = 3
)

var Ticket_State_name = map[int32]string{
	0: "New",
	1: "ACK",
	2: "CANCEL",
	3: "TIMEOUT",
}

var Ticket_State_value = map[string]int32{
	"New":     0,
	"ACK":     1,
	"CANCEL":  2,
	"TIMEOUT": 3,
}

func (x Ticket_State) String() string {
	return proto.EnumName(Ticket_State_name, int32(x))
}

func (Ticket_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0, 0}
}

type TicketAck_Type int32

const (
	TicketAck_ACCEPT TicketAck_Type = 0
	TicketAck_CANCEL TicketAck_Type = 1
)

var TicketAck_Type_name = map[int32]string{
	0: "ACCEPT",
	1: "CANCEL",
}

var TicketAck_Type_value = map[string]int32{
	"ACCEPT": 0,
	"CANCEL": 1,
}

func (x TicketAck_Type) String() string {
	return proto.EnumName(TicketAck_Type_name, int32(x))
}

func (TicketAck_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2, 0}
}

type Ticket struct {
	//string publisher    = 1;
	SendTo   string `protobuf:"bytes,1,opt,name=sendTo,proto3" json:"sendTo,omitempty"`
	Cid      string `protobuf:"bytes,2,opt,name=cid,proto3" json:"cid,omitempty"`
	ByteSize int64  `protobuf:"varint,3,opt,name=byteSize,proto3" json:"byteSize,omitempty"`
	//int64  timeStamp    = 4;
	State                Ticket_State `protobuf:"varint,4,opt,name=state,proto3,enum=bitswap.message.pb.Ticket_State" json:"state,omitempty"`
	Duration             int64        `protobuf:"varint,5,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Ticket) Reset()         { *m = Ticket{} }
func (m *Ticket) String() string { return proto.CompactTextString(m) }
func (*Ticket) ProtoMessage()    {}
func (*Ticket) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Ticket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ticket.Unmarshal(m, b)
}
func (m *Ticket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ticket.Marshal(b, m, deterministic)
}
func (m *Ticket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ticket.Merge(m, src)
}
func (m *Ticket) XXX_Size() int {
	return xxx_messageInfo_Ticket.Size(m)
}
func (m *Ticket) XXX_DiscardUnknown() {
	xxx_messageInfo_Ticket.DiscardUnknown(m)
}

var xxx_messageInfo_Ticket proto.InternalMessageInfo

func (m *Ticket) GetSendTo() string {
	if m != nil {
		return m.SendTo
	}
	return ""
}

func (m *Ticket) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *Ticket) GetByteSize() int64 {
	if m != nil {
		return m.ByteSize
	}
	return 0
}

func (m *Ticket) GetState() Ticket_State {
	if m != nil {
		return m.State
	}
	return Ticket_New
}

func (m *Ticket) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

type Ticketlist struct {
	Items                []*Ticket `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Ticketlist) Reset()         { *m = Ticketlist{} }
func (m *Ticketlist) String() string { return proto.CompactTextString(m) }
func (*Ticketlist) ProtoMessage()    {}
func (*Ticketlist) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *Ticketlist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ticketlist.Unmarshal(m, b)
}
func (m *Ticketlist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ticketlist.Marshal(b, m, deterministic)
}
func (m *Ticketlist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ticketlist.Merge(m, src)
}
func (m *Ticketlist) XXX_Size() int {
	return xxx_messageInfo_Ticketlist.Size(m)
}
func (m *Ticketlist) XXX_DiscardUnknown() {
	xxx_messageInfo_Ticketlist.DiscardUnknown(m)
}

var xxx_messageInfo_Ticketlist proto.InternalMessageInfo

func (m *Ticketlist) GetItems() []*Ticket {
	if m != nil {
		return m.Items
	}
	return nil
}

type TicketAck struct {
	Publisher            string         `protobuf:"bytes,1,opt,name=publisher,proto3" json:"publisher,omitempty"`
	Receiver             string         `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Cid                  string         `protobuf:"bytes,3,opt,name=cid,proto3" json:"cid,omitempty"`
	Type                 TicketAck_Type `protobuf:"varint,4,opt,name=type,proto3,enum=bitswap.message.pb.TicketAck_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *TicketAck) Reset()         { *m = TicketAck{} }
func (m *TicketAck) String() string { return proto.CompactTextString(m) }
func (*TicketAck) ProtoMessage()    {}
func (*TicketAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *TicketAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TicketAck.Unmarshal(m, b)
}
func (m *TicketAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TicketAck.Marshal(b, m, deterministic)
}
func (m *TicketAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TicketAck.Merge(m, src)
}
func (m *TicketAck) XXX_Size() int {
	return xxx_messageInfo_TicketAck.Size(m)
}
func (m *TicketAck) XXX_DiscardUnknown() {
	xxx_messageInfo_TicketAck.DiscardUnknown(m)
}

var xxx_messageInfo_TicketAck proto.InternalMessageInfo

func (m *TicketAck) GetPublisher() string {
	if m != nil {
		return m.Publisher
	}
	return ""
}

func (m *TicketAck) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *TicketAck) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *TicketAck) GetType() TicketAck_Type {
	if m != nil {
		return m.Type
	}
	return TicketAck_ACCEPT
}

type TicketAcklist struct {
	Items                []*TicketAck `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TicketAcklist) Reset()         { *m = TicketAcklist{} }
func (m *TicketAcklist) String() string { return proto.CompactTextString(m) }
func (*TicketAcklist) ProtoMessage()    {}
func (*TicketAcklist) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *TicketAcklist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TicketAcklist.Unmarshal(m, b)
}
func (m *TicketAcklist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TicketAcklist.Marshal(b, m, deterministic)
}
func (m *TicketAcklist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TicketAcklist.Merge(m, src)
}
func (m *TicketAcklist) XXX_Size() int {
	return xxx_messageInfo_TicketAcklist.Size(m)
}
func (m *TicketAcklist) XXX_DiscardUnknown() {
	xxx_messageInfo_TicketAcklist.DiscardUnknown(m)
}

var xxx_messageInfo_TicketAcklist proto.InternalMessageInfo

func (m *TicketAcklist) GetItems() []*TicketAck {
	if m != nil {
		return m.Items
	}
	return nil
}

type Message struct {
	Wantlist *Message_Wantlist `protobuf:"bytes,1,opt,name=wantlist,proto3" json:"wantlist,omitempty"`
	Blocks   [][]byte          `protobuf:"bytes,2,rep,name=blocks,proto3" json:"blocks,omitempty"`
	Payload  []*Message_Block  `protobuf:"bytes,3,rep,name=payload,proto3" json:"payload,omitempty"`
	//Ticketlist ticketlist = 4;
	//TicketAcklist ticketAcklist = 5;
	Ticketlist           []*Ticket    `protobuf:"bytes,4,rep,name=ticketlist,proto3" json:"ticketlist,omitempty"`
	TicketAcklist        []*TicketAck `protobuf:"bytes,5,rep,name=ticketAcklist,proto3" json:"ticketAcklist,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetWantlist() *Message_Wantlist {
	if m != nil {
		return m.Wantlist
	}
	return nil
}

func (m *Message) GetBlocks() [][]byte {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func (m *Message) GetPayload() []*Message_Block {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetTicketlist() []*Ticket {
	if m != nil {
		return m.Ticketlist
	}
	return nil
}

func (m *Message) GetTicketAcklist() []*TicketAck {
	if m != nil {
		return m.TicketAcklist
	}
	return nil
}

type Message_Wantlist struct {
	Entries              []*Message_Wantlist_Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	Full                 bool                      `protobuf:"varint,2,opt,name=full,proto3" json:"full,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *Message_Wantlist) Reset()         { *m = Message_Wantlist{} }
func (m *Message_Wantlist) String() string { return proto.CompactTextString(m) }
func (*Message_Wantlist) ProtoMessage()    {}
func (*Message_Wantlist) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4, 0}
}

func (m *Message_Wantlist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Wantlist.Unmarshal(m, b)
}
func (m *Message_Wantlist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Wantlist.Marshal(b, m, deterministic)
}
func (m *Message_Wantlist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Wantlist.Merge(m, src)
}
func (m *Message_Wantlist) XXX_Size() int {
	return xxx_messageInfo_Message_Wantlist.Size(m)
}
func (m *Message_Wantlist) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Wantlist.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Wantlist proto.InternalMessageInfo

func (m *Message_Wantlist) GetEntries() []*Message_Wantlist_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Message_Wantlist) GetFull() bool {
	if m != nil {
		return m.Full
	}
	return false
}

type Message_Wantlist_Entry struct {
	Block                []byte   `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	Priority             int32    `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	Cancel               bool     `protobuf:"varint,3,opt,name=cancel,proto3" json:"cancel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_Wantlist_Entry) Reset()         { *m = Message_Wantlist_Entry{} }
func (m *Message_Wantlist_Entry) String() string { return proto.CompactTextString(m) }
func (*Message_Wantlist_Entry) ProtoMessage()    {}
func (*Message_Wantlist_Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4, 0, 0}
}

func (m *Message_Wantlist_Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Wantlist_Entry.Unmarshal(m, b)
}
func (m *Message_Wantlist_Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Wantlist_Entry.Marshal(b, m, deterministic)
}
func (m *Message_Wantlist_Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Wantlist_Entry.Merge(m, src)
}
func (m *Message_Wantlist_Entry) XXX_Size() int {
	return xxx_messageInfo_Message_Wantlist_Entry.Size(m)
}
func (m *Message_Wantlist_Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Wantlist_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Wantlist_Entry proto.InternalMessageInfo

func (m *Message_Wantlist_Entry) GetBlock() []byte {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *Message_Wantlist_Entry) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *Message_Wantlist_Entry) GetCancel() bool {
	if m != nil {
		return m.Cancel
	}
	return false
}

type Message_Block struct {
	Prefix               []byte   `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message_Block) Reset()         { *m = Message_Block{} }
func (m *Message_Block) String() string { return proto.CompactTextString(m) }
func (*Message_Block) ProtoMessage()    {}
func (*Message_Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4, 1}
}

func (m *Message_Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message_Block.Unmarshal(m, b)
}
func (m *Message_Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message_Block.Marshal(b, m, deterministic)
}
func (m *Message_Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message_Block.Merge(m, src)
}
func (m *Message_Block) XXX_Size() int {
	return xxx_messageInfo_Message_Block.Size(m)
}
func (m *Message_Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Message_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Message_Block proto.InternalMessageInfo

func (m *Message_Block) GetPrefix() []byte {
	if m != nil {
		return m.Prefix
	}
	return nil
}

func (m *Message_Block) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("bitswap.message.pb.Ticket_State", Ticket_State_name, Ticket_State_value)
	proto.RegisterEnum("bitswap.message.pb.TicketAck_Type", TicketAck_Type_name, TicketAck_Type_value)
	proto.RegisterType((*Ticket)(nil), "bitswap.message.pb.Ticket")
	proto.RegisterType((*Ticketlist)(nil), "bitswap.message.pb.Ticketlist")
	proto.RegisterType((*TicketAck)(nil), "bitswap.message.pb.TicketAck")
	proto.RegisterType((*TicketAcklist)(nil), "bitswap.message.pb.TicketAcklist")
	proto.RegisterType((*Message)(nil), "bitswap.message.pb.Message")
	proto.RegisterType((*Message_Wantlist)(nil), "bitswap.message.pb.Message.Wantlist")
	proto.RegisterType((*Message_Wantlist_Entry)(nil), "bitswap.message.pb.Message.Wantlist.Entry")
	proto.RegisterType((*Message_Block)(nil), "bitswap.message.pb.Message.Block")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 572 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdd, 0x8e, 0x93, 0x40,
	0x14, 0x5e, 0x0a, 0xf4, 0xe7, 0xb4, 0x35, 0x64, 0x62, 0x0c, 0x21, 0xfe, 0x20, 0xf1, 0xa2, 0x31,
	0x91, 0x35, 0x6d, 0xb2, 0x17, 0x5e, 0x98, 0xb0, 0x58, 0x13, 0x7f, 0x76, 0xd5, 0x59, 0x8c, 0xd7,
	0x40, 0xa7, 0xdd, 0x49, 0x29, 0x10, 0x98, 0x5a, 0xf1, 0x71, 0x7c, 0x01, 0x9f, 0xc0, 0x7b, 0x9f,
	0xc0, 0x4b, 0x9f, 0xc5, 0xcc, 0x30, 0xb0, 0xdd, 0xe8, 0xae, 0x7b, 0x77, 0xbe, 0xe1, 0x3b, 0xdf,
	0x9c, 0xef, 0x9b, 0x13, 0x60, 0xbc, 0x21, 0x65, 0x19, 0xae, 0x88, 0x9b, 0x17, 0x19, 0xcb, 0x10,
	0x8a, 0x28, 0x2b, 0x77, 0x61, 0xee, 0xb6, 0xc7, 0x91, 0xf5, 0x64, 0x45, 0xd9, 0xf9, 0x36, 0x72,
	0xe3, 0x6c, 0x73, 0xb8, 0xca, 0x56, 0xd9, 0xa1, 0xa0, 0x46, 0xdb, 0xa5, 0x40, 0x02, 0x88, 0xaa,
	0x96, 0x70, 0x7e, 0x29, 0xd0, 0x0d, 0x68, 0xbc, 0x26, 0x0c, 0xdd, 0x81, 0x6e, 0x49, 0xd2, 0x45,
	0x90, 0x99, 0x8a, 0xad, 0x4c, 0x06, 0x58, 0x22, 0x64, 0x80, 0x1a, 0xd3, 0x85, 0xd9, 0x11, 0x87,
	0xbc, 0x44, 0x16, 0xf4, 0xa3, 0x8a, 0x91, 0x33, 0xfa, 0x95, 0x98, 0xaa, 0xad, 0x4c, 0x54, 0xdc,
	0x62, 0x74, 0x04, 0x7a, 0xc9, 0x42, 0x46, 0x4c, 0xcd, 0x56, 0x26, 0xb7, 0xa6, 0xb6, 0xfb, 0xf7,
	0x8c, 0x6e, 0x7d, 0xa1, 0x7b, 0xc6, 0x79, 0xb8, 0xa6, 0x73, 0xcd, 0xc5, 0xb6, 0x08, 0x19, 0xcd,
	0x52, 0x53, 0xaf, 0x35, 0x1b, 0xec, 0x4c, 0x41, 0x17, 0x5c, 0xd4, 0x03, 0xf5, 0x94, 0xec, 0x8c,
	0x03, 0x5e, 0x78, 0xfe, 0x1b, 0x43, 0x41, 0x00, 0x5d, 0xdf, 0x3b, 0xf5, 0xe7, 0x6f, 0x8d, 0x0e,
	0x1a, 0x42, 0x2f, 0x78, 0x75, 0x32, 0x7f, 0xf7, 0x31, 0x30, 0x54, 0xe7, 0x39, 0x40, 0x7d, 0x4d,
	0x42, 0x4b, 0x86, 0x9e, 0x82, 0x4e, 0x19, 0xd9, 0x94, 0xa6, 0x62, 0xab, 0x93, 0xe1, 0xd4, 0xba,
	0x7a, 0x2a, 0x5c, 0x13, 0x9d, 0xef, 0x0a, 0x0c, 0xea, 0x13, 0x2f, 0x5e, 0xa3, 0xbb, 0x30, 0xc8,
	0xb7, 0x51, 0x42, 0xcb, 0x73, 0x52, 0xc8, 0x78, 0x2e, 0x0e, 0xf8, 0xec, 0x05, 0x89, 0x09, 0xfd,
	0x4c, 0x0a, 0x19, 0x53, 0x8b, 0x9b, 0xf4, 0xd4, 0x8b, 0xf4, 0x8e, 0x40, 0x63, 0x55, 0xde, 0x04,
	0xe4, 0x5c, 0x3d, 0x8a, 0x17, 0xaf, 0xdd, 0xa0, 0xca, 0x09, 0x16, 0x7c, 0xe7, 0x3e, 0x68, 0x1c,
	0x71, 0xcb, 0x9e, 0xef, 0xcf, 0xdf, 0x07, 0xc6, 0xc1, 0x9e, 0x7d, 0xc5, 0x79, 0x01, 0xe3, 0xb6,
	0x4f, 0x98, 0x9e, 0x5d, 0x36, 0x7d, 0xef, 0xda, 0x9b, 0x1a, 0xdf, 0xdf, 0x34, 0xe8, 0x9d, 0xd4,
	0xdf, 0xd1, 0x4b, 0xe8, 0xef, 0xc2, 0x54, 0x24, 0x28, 0x4c, 0x0f, 0xa7, 0x8f, 0xfe, 0xa5, 0x21,
	0xe9, 0xee, 0x27, 0xc9, 0x3d, 0xd6, 0x7e, 0xfe, 0x7e, 0x70, 0x80, 0xdb, 0x5e, 0xbe, 0x59, 0x51,
	0x92, 0xc5, 0xeb, 0xd2, 0xec, 0xd8, 0xea, 0x64, 0x84, 0x25, 0x42, 0x1e, 0xf4, 0xf2, 0xb0, 0x4a,
	0xb2, 0x90, 0xe7, 0xc3, 0x47, 0x7c, 0x78, 0x9d, 0xfc, 0x31, 0x6f, 0x92, 0xda, 0x4d, 0x1f, 0x7a,
	0x06, 0xc0, 0xda, 0x67, 0x36, 0xb5, 0xff, 0xbe, 0xee, 0x1e, 0x1b, 0xf9, 0x30, 0x66, 0xfb, 0x81,
	0x99, 0xfa, 0x4d, 0x72, 0xba, 0xdc, 0x63, 0xfd, 0x50, 0xa0, 0xdf, 0x18, 0x47, 0xaf, 0xa1, 0x47,
	0x52, 0x56, 0x50, 0xd2, 0x64, 0xfe, 0xf8, 0x26, 0x79, 0xb9, 0xf3, 0x94, 0x15, 0x55, 0xe3, 0x4c,
	0x0a, 0x20, 0x04, 0xda, 0x72, 0x9b, 0x24, 0x62, 0xa1, 0xfa, 0x58, 0xd4, 0xd6, 0x07, 0xd0, 0x05,
	0x17, 0xdd, 0x06, 0x5d, 0x64, 0x28, 0x9e, 0x65, 0x84, 0x6b, 0xc0, 0xf7, 0x30, 0x2f, 0x68, 0x56,
	0x50, 0x56, 0x89, 0x36, 0x1d, 0xb7, 0x98, 0xbf, 0x41, 0x1c, 0xa6, 0x31, 0x49, 0xc4, 0x2a, 0xf6,
	0xb1, 0x44, 0xd6, 0x0c, 0x74, 0x11, 0x2c, 0x27, 0xe4, 0x05, 0x59, 0xd2, 0x2f, 0x52, 0x53, 0x22,
	0x3e, 0xc7, 0x22, 0x64, 0xa1, 0x10, 0x1c, 0x61, 0x51, 0x47, 0x5d, 0xf1, 0xf3, 0x98, 0xfd, 0x09,
	0x00, 0x00, 0xff, 0xff, 0x4f, 0xcc, 0x46, 0xd1, 0x90, 0x04, 0x00, 0x00,
}
