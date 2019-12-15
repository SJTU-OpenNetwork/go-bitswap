package tickets

import (
	"crypto/rand"
	pb "github.com/SJTU-OpenNetwork/go-bitswap/message/pb"
	cid "github.com/ipfs/go-cid"
	u "github.com/ipfs/go-ipfs-util"
	"github.com/libp2p/go-libp2p-core/peer"
	peertest "github.com/libp2p/go-libp2p-core/test"
	mh "github.com/multiformats/go-multihash"
	"testing"

	"encoding/json"
)

func makeRandomCid() cid.Cid {
	p := make([]byte, 256)
	rand.Read(p)
	h, _ := mh.Sum(p, mh.SHA3, 4)
	cid := cid.NewCidV1(7, h)
	return cid
}

func mkFakeCid(s string) cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(s)))
}

func mkFakePid() (peer.ID, error){
	return peertest.RandPeerID()
}

func loggable2json(loginfo Loggable) string {
	jsonString, err := json.MarshalIndent(loginfo.Loggable(), "", "\t")
	if err != nil{
		log.Error(err)
		return ""
	}
	return string(jsonString)
}

var peer1 peer.ID
var peer2 peer.ID
var cids []cid.Cid
const numCids = 10

func init(){
	peer1, _ = mkFakePid()
	peer2, _ = mkFakePid()
	cids = make([]cid.Cid, 0, numCids)
	for i:=0; i<numCids; i++{
		cids = append(cids, makeRandomCid())
	}
}

func TestBasicTicket(t *testing.T) {
	ticket1 := CreateBasicTicket(peer2, cids[0], 1024)

	t.Logf("Create ticket:\n%s", ticket1.BasicInfo())

	ticket1proto := ticket1.ToProto()
	t.Log("Convert to proto.")

	ticket1byte, err := ticket1proto.Marshal()
	if(err != nil){
		t.Error(err)
	}
	var ticket1protoresume pb.Ticket
	ticket1protoresume.Unmarshal(ticket1byte)
	t.Log("Marshal and unmarshal ticket.")

	ticket1resume, err := NewBasicTicket(&ticket1protoresume)
	if(err != nil){
		t.Fail()
	}

	t.Logf("Resumed ticket:\n%s", ticket1resume.BasicInfo())

	ticket1resume.ACKed()
	t.Logf("Accept the ticket.\nTicket State:%s", ticket1resume.GetStateString())
	ticket1resume.Canceled()
	t.Logf("Cancel the ticket.\nTicket State:%s", ticket1resume.GetStateString())
}

func TestBasicTicketAck(t *testing.T) {
	ticketack1 := CreateBasicTicketAck(peer1, peer2, cids[0], 0)
	t.Log("Create ticket ack:\n", loggable2json(ticketack1))

	ticketack1proto := ticketack1.ToProto()
	t.Log("Generate proto:\n", ticketack1proto)

	ticketack1byte,_ := ticketack1proto.Marshal()
	var ticketack1protoresume pb.TicketAck
	ticketack1protoresume.Unmarshal(ticketack1byte)
	ticketack1resume, _ := NewBasicTicketAck(&ticketack1protoresume)

	t.Log("Marshal and Unmarshal ack:\n", ticketack1resume.Loggable())
}

func TestLinkedTicketStore(t *testing.T) {

}