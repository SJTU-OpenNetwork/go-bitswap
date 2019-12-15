// Package decision implements the decision engine for the bitswap service.
package decision

import (
	"context"
	"fmt"
	"github.com/SJTU-OpenNetwork/go-bitswap/tickets"
	"sync"
	"time"

	bsmsg "github.com/SJTU-OpenNetwork/go-bitswap/message"
	wl "github.com/SJTU-OpenNetwork/go-bitswap/wantlist"
	"github.com/google/uuid"
	cid "github.com/ipfs/go-cid"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log"
	"github.com/ipfs/go-peertaskqueue"
	"github.com/ipfs/go-peertaskqueue/peertask"
	process "github.com/jbenet/goprocess"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// TODO consider taking responsibility for other types of requests. For
// example, there could be a |cancelQueue| for all of the cancellation
// messages that need to go out. There could also be a |wantlistQueue| for
// the local peer's wantlists. Alternatively, these could all be bundled
// into a single, intelligent global queue that efficiently
// batches/combines and takes all of these into consideration.
//
// Right now, messages go onto the network for four reasons:
// 1. an initial `sendwantlist` message to a provider of the first key in a
//    request
// 2. a periodic full sweep of `sendwantlist` messages to all providers
// 3. upon receipt of blocks, a `cancel` message to all peers
// 4. draining the priority queue of `blockrequests` from peers
//
// Presently, only `blockrequests` are handled by the decision engine.
// However, there is an opportunity to give it more responsibility! If the
// decision engine is given responsibility for all of the others, it can
// intelligently decide how to combine requests efficiently.
//
// Some examples of what would be possible:
//
// * when sending out the wantlists, include `cancel` requests
// * when handling `blockrequests`, include `sendwantlist` and `cancel` as
//   appropriate
// * when handling `cancel`, if we recently received a wanted block from a
//   peer, include a partial wantlist that contains a few other high priority
//   blocks
//
// In a sense, if we treat the decision engine as a black box, it could do
// whatever it sees fit to produce desired outcomes (get wanted keys
// quickly, maintain good relationships with peers, etc).

var log = logging.Logger("engine")

const (
	// outboxChanBuffer must be 0 to prevent stale messages from being sent
	outboxChanBuffer = 0

	ticketOutboxChanBuffer = 10

	// Number of concurrent workers that pull tasks off the request queue
	taskWorkerCount = 8
	// maxMessageSize is the maximum size of the batched payload
	maxMessageSize = 512 * 1024
	// tagFormat is the tag given to peers associated an engine
	tagFormat = "bs-engine-%s-%s"

	// queuedTagWeight is the default weight for peers that have work queued
	// on their behalf.
	queuedTagWeight = 10

	// the alpha for the EWMA used to track short term usefulness
	shortTermAlpha = 0.5

	// the alpha for the EWMA used to track long term usefulness
	longTermAlpha = 0.05

	// long term ratio defines what "long term" means in terms of the
	// shortTerm duration. Peers that interact once every longTermRatio are
	// considered useful over the long term.
	longTermRatio = 10

	// long/short term scores for tagging peers
	longTermScore  = 10 // this is a high tag but it grows _very_ slowly.
	shortTermScore = 10 // this is a high tag but it'll go away quickly if we aren't using the peer.

	// Number of concurrent workers that process requests to the blockstore
	blockstoreWorkerCount = 128

	defaultRequestCapacity = 128
	defaultRequestSizeCapacity = 128 * 512 * 1024
)

var (
	// how frequently the engine should sample usefulness. Peers that
	// interact every shortTerm time period are considered "active".
	//
	// this is only a variable to make testing easier.
	shortTerm = 10 * time.Second
)

// Envelope contains a message for a Peer.
type Envelope struct {
	// Peer is the intended recipient.
	Peer peer.ID

	// Message is the payload.
	Message bsmsg.BitSwapMessage

	// A callback to notify the decision queue that the task is complete
	Sent func()
}

// PeerTagger covers the methods on the connection manager used by the decision
// engine to tag peers
type PeerTagger interface {
	TagPeer(peer.ID, string, int)
	UntagPeer(p peer.ID, tag string)
}

// Engine manages sending requested blocks to peers.
// TODO: sending blocks and tickets - Jerry
type Engine struct {
	// peerRequestQueue is a priority queue of requests received from peers.
	// Requests are popped from the queue, packaged up, and placed in the
	// outbox.
	peerRequestQueue *peertaskqueue.PeerTaskQueue
	// Capacity of RequestQueue
	// Tickets instead of blocks would be sent if request queue is full.
	requestRecorder *taskQueueRecorder

	ticketStore tickets.TicketStore

	// FIXME it's a bit odd for the client and the worker to both share memory
	// (both modify the peerRequestQueue) and also to communicate over the
	// workSignal channel. consider sending requests over the channel and
	// allowing the worker to have exclusive access to the peerRequestQueue. In
	// that case, no lock would be required.
	workSignal chan struct{}

	ticketSignal chan struct{}

	// outbox contains outgoing messages to peers. This is owned by the
	// taskWorker goroutine
	outbox chan (<-chan *Envelope)

	// TODO: ticketOutbox contains outgoing messages (tickets or ticket acks) to peers. - Jerry
	ticketOutbox chan *Envelope

	bsm *blockstoreManager

	peerTagger PeerTagger

	tagQueued, tagUseful string

	lock sync.Mutex // protects the fields immediatly below
	// ledgerMap lists Ledgers by their Partner key.
	ledgerMap map[peer.ID]*ledger

	ticker *time.Ticker

	taskWorkerLock  sync.Mutex
	taskWorkerCount int
}

// NewEngine creates a new block sending engine for the given block store
func NewEngine(ctx context.Context, bs bstore.Blockstore, peerTagger PeerTagger, ts tickets.TicketStore) *Engine {
	e := &Engine{
		ledgerMap:       make(map[peer.ID]*ledger),
		bsm:             newBlockstoreManager(ctx, bs, blockstoreWorkerCount),
		peerTagger:      peerTagger,
		outbox:          make(chan (<-chan *Envelope), outboxChanBuffer),	//TODO: Why not make outbox buffered directly ??? - Riften
		ticketOutbox:	 make(chan *Envelope, ticketOutboxChanBuffer),
		workSignal:      make(chan struct{}, 1),
		ticker:          time.NewTicker(time.Millisecond * 100),
		taskWorkerCount: taskWorkerCount,
		// TODO: add ticketOutbox - Jerry They can share the same outbox - Riften
		// ticketOutbox: 	 make(chan (<-chan *Envelope), outboxChanBuffer),
		requestRecorder: &taskQueueRecorder{
			blockNumber: 0,
			blockSize: 0,
			maxNumber:defaultRequestCapacity,
			maxSize:defaultRequestSizeCapacity,
		},
        ticketStore:     ts,
	}
	e.tagQueued = fmt.Sprintf(tagFormat, "queued", uuid.New().String())
	e.tagUseful = fmt.Sprintf(tagFormat, "useful", uuid.New().String())
	e.peerRequestQueue = peertaskqueue.New(
		peertaskqueue.OnPeerAddedHook(e.onPeerAdded),
		peertaskqueue.OnPeerRemovedHook(e.onPeerRemoved))
	go e.scoreWorker(ctx)
	return e
}

// Start up workers to handle requests from other nodes for the data on this node
// TODO: Start ticketWorkers - Riften
func (e *Engine) StartWorkers(ctx context.Context, px process.Process) {
	// Start up blockstore manager
	e.bsm.start(px)

	for i := 0; i < e.taskWorkerCount; i++ {
		px.Go(func(px process.Process) {
			e.taskWorker(ctx)
		})
	}
}

// scoreWorker keeps track of how "useful" our peers are, updating scores in the
// connection manager.
//
// It does this by tracking two scores: short-term usefulness and long-term
// usefulness. Short-term usefulness is sampled frequently and highly weights
// new observations. Long-term usefulness is sampled less frequently and highly
// weights on long-term trends.
//
// In practice, we do this by keeping two EWMAs. If we see an interaction
// within the sampling period, we record the score, otherwise, we record a 0.
// The short-term one has a high alpha and is sampled every shortTerm period.
// The long-term one has a low alpha and is sampled every
// longTermRatio*shortTerm period.
//
// To calculate the final score, we sum the short-term and long-term scores then
// adjust it ±25% based on our debt ratio. Peers that have historically been
// more useful to us than we are to them get the highest score.
func (e *Engine) scoreWorker(ctx context.Context) {
	ticker := time.NewTicker(shortTerm)
	defer ticker.Stop()

	type update struct {
		peer  peer.ID
		score int
	}
	var (
		lastShortUpdate, lastLongUpdate time.Time
		updates                         []update
	)

	for i := 0; ; i = (i + 1) % longTermRatio {
		var now time.Time
		select {
		case now = <-ticker.C:
		case <-ctx.Done():
			return
		}

		// The long term update ticks every `longTermRatio` short
		// intervals.
		updateLong := i == 0

		e.lock.Lock()
		for _, ledger := range e.ledgerMap {
			ledger.lk.Lock()

			// Update the short-term score.
			if ledger.lastExchange.After(lastShortUpdate) {
				ledger.shortScore = ewma(ledger.shortScore, shortTermScore, shortTermAlpha)
			} else {
				ledger.shortScore = ewma(ledger.shortScore, 0, shortTermAlpha)
			}

			// Update the long-term score.
			if updateLong {
				if ledger.lastExchange.After(lastLongUpdate) {
					ledger.longScore = ewma(ledger.longScore, longTermScore, longTermAlpha)
				} else {
					ledger.longScore = ewma(ledger.longScore, 0, longTermAlpha)
				}
			}

			// Calculate the new score.
			//
			// The accounting score adjustment prefers peers _we_
			// need over peers that need us. This doesn't help with
			// leeching.
			score := int((ledger.shortScore + ledger.longScore) * ((ledger.Accounting.Score())*.5 + .75))

			// Avoid updating the connection manager unless there's a change. This can be expensive.
			if ledger.score != score {
				// put these in a list so we can perform the updates outside _global_ the lock.
				updates = append(updates, update{ledger.Partner, score})
				ledger.score = score
			}
			ledger.lk.Unlock()
		}
		e.lock.Unlock()

		// record the times.
		lastShortUpdate = now
		if updateLong {
			lastLongUpdate = now
		}

		// apply the updates
		for _, update := range updates {
			if update.score == 0 {
				e.peerTagger.UntagPeer(update.peer, e.tagUseful)
			} else {
				e.peerTagger.TagPeer(update.peer, e.tagUseful, update.score)
			}
		}
		// Keep the memory. It's not much and it saves us from having to allocate.
		updates = updates[:0]
	}
}

func (e *Engine) onPeerAdded(p peer.ID) {
	e.peerTagger.TagPeer(p, e.tagQueued, queuedTagWeight)
}

func (e *Engine) onPeerRemoved(p peer.ID) {
	e.peerTagger.UntagPeer(p, e.tagQueued)
}

// WantlistForPeer returns the currently understood want list for a given peer
func (e *Engine) WantlistForPeer(p peer.ID) (out []wl.Entry) {
	partner := e.findOrCreate(p)
	partner.lk.Lock()
	defer partner.lk.Unlock()
	return partner.wantList.SortedEntries()
}

// LedgerForPeer returns aggregated data about blocks swapped and communication
// with a given peer.
func (e *Engine) LedgerForPeer(p peer.ID) *Receipt {
	ledger := e.findOrCreate(p)

	ledger.lk.Lock()
	defer ledger.lk.Unlock()

	return &Receipt{
		Peer:      ledger.Partner.String(),
		Value:     ledger.Accounting.Value(),
		Sent:      ledger.Accounting.BytesSent,
		Recv:      ledger.Accounting.BytesRecv,
		Exchanged: ledger.ExchangeCount(),
	}
}

// Each taskWorker pulls items off the request queue up and adds them to an
// envelope. The envelope is passed off to the bitswap workers, which send
// the message to the network.
func (e *Engine) taskWorker(ctx context.Context) {
	defer e.taskWorkerExit()
	for {
		oneTimeUse := make(chan *Envelope, 1) // buffer to prevent blocking
		select {
		case <-ctx.Done():
			return
		case e.outbox <- oneTimeUse:
		}
		// receiver is ready for an outoing envelope. let's prepare one. first,
		// we must acquire a task from the PQ...
		envelope, err := e.nextEnvelope(ctx)
		if err != nil {
			close(oneTimeUse)
			return // ctx cancelled
		}
		oneTimeUse <- envelope // buffered. won't block
		close(oneTimeUse)
	}
}

// taskWorkerExit handles cleanup of task workers
func (e *Engine) taskWorkerExit() {
	e.taskWorkerLock.Lock()
	defer e.taskWorkerLock.Unlock()

	e.taskWorkerCount--
	if e.taskWorkerCount == 0 {
		close(e.outbox)
	}
}

// TODO: add ticketWorker - sending tickets cannot use the same thread as normal messages - Jerry
// - Get ticket out of sendTicketStore
// - Pack it into envelop
// - put it to outbox
// For now it runs simultaneously with taskWorker and share the same outbox
func (e *Engine) ticketWorker(ctx context.Context) {
	//defer e.task
	// defer e.taskWorkerExit()
}


// TODO: add ticketWorkerExit - Jerry
func (e *Engine) ticketWorkerExit() {

}

//TODO: channel may be blocked if there is too more tickets to send. Try to send it using some kind of queue
//func (e )
func (e *Engine) SendTickets(target peer.ID, tickets []tickets.Ticket) {
	//1. Build Envelope
	//Envelope:
	//	- Peer: the target send ticket to
	//	- Message: the message contains the tickets
	//	- Sent: the call back func called when sending complete
	//
	//2. Send envelope to outbox

	msg := bsmsg.New(false)
	msg.AddTickets(tickets)
	if msg.Empty(){
		return
	}
	envelope := &Envelope{
		Peer:    target,
		Message: msg,
		Sent:    func(){
			e.ticketStore.AddTickets(tickets)
		},
	}
	e.ticketOutbox <- envelope
}
//func (e *Engine) nextTicketEnvelope(ctx context.Context) (*Envelope, error) {
//	for {
//		nextTask := e.ticketStore.PopTickets()
//		for nextTask == nil {
//			select{
//			case <- ctx.Done():
//				return nil, ctx.Err()
//			case <- e.ticketSignal:
//				nextTask = e.ticketStore.PopTickets()
//			}
//		}
//
//		msg := bsmsg.New(false)
//		for _, t := range nextTask.Tickets {
//			msg.AddTicket(t)
//		}
//		if msg.Empty() {
//			continue
//		}
//		return &Envelope{
//			Peer: nextTask.Target,
//			Message: msg,
//			Sent: func(){
//				select {
//				// Tell the taskworker
//				case e.workSignal <- struct{}{}:
//				default:
//				}
//			},
//		}, nil
//	}
//}

// nextEnvelope runs in the taskWorker goroutine. Returns an error if the
// context is cancelled before the next Envelope can be created.
func (e *Engine) nextEnvelope(ctx context.Context) (*Envelope, error) {
	for {
		nextTask := e.peerRequestQueue.PopBlock()
		e.requestRecorder.BlocksPop(nextTask)
		//WorkSignal will be sent when a new task is pushed into request queue.
		for nextTask == nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-e.workSignal:
				nextTask = e.peerRequestQueue.PopBlock()
				e.requestRecorder.BlocksPop(nextTask)
			case <-e.ticker.C:
				e.peerRequestQueue.ThawRound()
				nextTask = e.peerRequestQueue.PopBlock()
				e.requestRecorder.BlocksPop(nextTask)
			}
		}

		// with a task in hand, we're ready to prepare the envelope...
		blockCids := cid.NewSet()
		for _, t := range nextTask.Tasks {
			blockCids.Add(t.Identifier.(cid.Cid))
		}
		blks := e.bsm.getBlocks(ctx, blockCids.Keys())

		msg := bsmsg.New(true)
		for _, b := range blks {
			msg.AddBlock(b)
		}

		if msg.Empty() {
			// If we don't have the block, don't hold that against the peer
			// make sure to update that the task has been 'completed'
			nextTask.Done(nextTask.Tasks)
			continue
		}

		return &Envelope{
			Peer:    nextTask.Target,
			Message: msg,
			Sent: func() {
				nextTask.Done(nextTask.Tasks)
				select {
				case e.workSignal <- struct{}{}:
					// work completing may mean that our queue will provide new
					// work to be done.
				default:
				}
			},
		}, nil
	}
}

// Outbox returns a channel of one-time use Envelope channels.
func (e *Engine) Outbox() <-chan (<-chan *Envelope) {
	return e.outbox
}

func (e *Engine) TicketOutbox() chan *Envelope{
	return e.ticketOutbox
}
//
//func (e *Engine) TicketOutbox() <-chan (<-chan *Envelope) {
//	return e.ticketOutbox
//}

// Peers returns a slice of Peers with whom the local node has active sessions.
func (e *Engine) Peers() []peer.ID {
	e.lock.Lock()
	defer e.lock.Unlock()

	response := make([]peer.ID, 0, len(e.ledgerMap))

	for _, ledger := range e.ledgerMap {
		response = append(response, ledger.Partner)
	}
	return response
}

func createTicketsFromEntry(target peer.ID, tasks []peertask.Task, blockSizes map[cid.Cid]int) []tickets.Ticket{
	res := make([]tickets.Ticket, 0)
	for _, task := range tasks{
		ticketSize, ok := blockSizes[task.Identifier.(cid.Cid)]
		if ok {
			ticket := tickets.CreateBasicTicket(target, task.Identifier.(cid.Cid), int64(ticketSize))
			res = append(res, ticket)
		} else {
			log.Error("Can not get block sizes when create ticket")
		}
	}
	return res
}

// MessageReceived performs book-keeping. Returns error if passed invalid
// arguments.
func (e *Engine) MessageReceived(ctx context.Context, p peer.ID, m bsmsg.BitSwapMessage) {
	if m.Empty() {
		log.Debugf("received empty message from %s", p)
	}

	newWorkExists := false
	defer func() {
		if newWorkExists {
			e.signalNewWork()
		}
	}()

	// Get block sizes
	entries := m.Wantlist()
	wantKs := cid.NewSet()
	for _, entry := range entries {
		if !entry.Cancel {
			wantKs.Add(entry.Cid)
		}
	}

	blockSizes := e.bsm.getBlockSizes(ctx, wantKs.Keys())

	l := e.findOrCreate(p)		//Find the ledger
	l.lk.Lock()
	defer l.lk.Unlock()
	if m.Full() {
		l.wantList = wl.New()
	}

	var msgSize int
	var activeEntries []peertask.Task
	var activeTickets []tickets.Ticket
	var queryCids []cid.Cid

	for _, entry := range m.Wantlist() {
		if entry.Cancel {
			log.Debugf("%s cancel %s", p, entry.Cid)
			l.CancelWant(entry.Cid)
			e.peerRequestQueue.Remove(entry.Cid, p)
		} else {
			log.Debugf("wants %s - %d", entry.Cid, entry.Priority)
			l.Wants(entry.Cid, entry.Priority)
			blockSize, ok := blockSizes[entry.Cid]
			if ok {
				// we have the block
				// although we have the block, but the network resources are limited
				// I may just send a ticket, which means I will send the block when the 
				// resources are sufficient - Jerry
				newWorkExists = true
				if msgSize+blockSize > maxMessageSize {
					if(e.requestRecorder.IsFull()){	//We do not have enough network resources. Send tickets instead of send the block
						//e.ticketStore.GetTickets()
						//Create ticket
						tmptickets := createTicketsFromEntry(p, activeEntries, blockSizes)
						activeTickets = append(activeTickets, tmptickets...)
						// Add it to sender
					} else {
						e.peerRequestQueue.PushBlock(p, activeEntries...)
					}
					activeEntries = []peertask.Task{}
					msgSize = 0
				}
				activeEntries = append(activeEntries, peertask.Task{Identifier: entry.Cid, Priority: entry.Priority})
				msgSize += blockSize
			} else {
				// although we do not have the block, but I have the ticket. - Jerry
				// So I can send a ticket with higher level - Jerry
				// Judge whether we have the tickets
				queryCids = append(queryCids, entry.Cid)
			}
		}
	}
	if len(activeEntries) > 0 {
		if e.requestRecorder.IsFull() {
			tmptickets := createTicketsFromEntry(p, activeEntries, blockSizes)
			activeTickets = append(activeTickets, tmptickets...)
		}else{
			e.peerRequestQueue.PushBlock(p, activeEntries...)
		}
	}

	tmpticketsmap, err := e.ticketStore.GetReceivedTicket(queryCids)
	if err != nil {
		log.Error(err)
	} else {
		for _, t := range tmpticketsmap{
			activeTickets = append(activeTickets, t)
		}
	}

	// TODO: Judge whether it is necessary to add ticket to ticket store early so that we can judge send time more correctly
	// TODO: add activeTickets to ticket sender and ticket store

	// Receive blocks
	blockCids := make([]cid.Cid, 0)
	for _, block := range m.Blocks() {
		log.Debugf("got block %s %d bytes", block, len(block.RawData()))
		l.ReceivedBytes(len(block.RawData()))
		blockCids = append(blockCids, block.Cid())
	}

    // Receive tickets
	if m.Tickets() != nil {
		e.handleTickets(ctx, m.Tickets())
	}

	// Receive Ticket Acks - Jerry 2019/12/14
	acks := m.TicketAcks()
    if acks != nil {
        e.handleTicketAcks(ctx, p, acks)
    }
}

//func (e *Engine) handleReceiveBlocks(cids []cid.Cid) {
//	acks, _ := e.ticketStore.PopSendingTasks(cids)
//	var activeEntries []peertask.Task
//    msgSize := 0
//	for index, ack := range acks {
//		blockSize, _ := ack.Size()
//		if msgSize + blockSize > maxMessageSize {
//			e.peerRequestQueue.PushBlock(ack.Receiver(), activeEntries...)
//			activeEntries = []peertask.Task{}
//			msgSize = 0
//		}
//		activeEntries = append(activeEntries, peertask.Task{Identifier:cids[index], Priority:10})
//		msgSize += blockSize
//	}
//	if len(activeEntries) > 0 {
//		e.peerRequestQueue.PushBlock(p, activeEntries...)
//	}
//}

func (e *Engine) handleTickets(ctx context.Context, tks []tickets.Ticket) {
	var cids, noblocks []cid.Cid
	ticketMap := make(map[cid.Cid] tickets.Ticket)
    rejectsMap := make(map[peer.ID] []tickets.Ticket)
    acceptsMap := make(map[peer.ID] []tickets.Ticket)

    for _, ticket := range tks {
        cids = append(cids, ticket.Cid())
        ticketMap[ticket.Cid()] = ticket
    }
    blockSizes := e.bsm.getBlockSizes(ctx, cids)
    for _, cid := range cids {
        _, ok := blockSizes[cid]
        if ok {
            ticket := ticketMap[cid]
            rejectsMap[ticket.SendTo()] = append(rejectsMap[ticket.SendTo()], ticket)
        } else {
            noblocks = append(noblocks, cid)
        }
    }

    localMap, _ := e.ticketStore.GetReceivedTicket(noblocks)
    for _, cid := range noblocks {
        local, ok := localMap[cid]
        nt := ticketMap[cid]
        if ok {
            if local.Level() <= nt.Level() {
                rejectsMap[nt.SendTo()] = append(rejectsMap[nt.SendTo()], nt)
            } else {
                rejectsMap[local.SendTo()] = append(rejectsMap[local.SendTo()], local)
                acceptsMap[nt.SendTo()] = append(acceptsMap[nt.SendTo()], nt)
            }
        } else {
            acceptsMap[nt.SendTo()] = append(acceptsMap[nt.SendTo()], nt)
        }
    }
    for reject := range rejectsMap {
        //reject rejectsMap[reject]
    }
    for accept := range acceptsMap {
        //accept acceptsMap[reject]
    }

}

func (e *Engine) handleTicketAcks(ctx context.Context, p peer.ID, acks []tickets.TicketAck) {
	ackMap := make(map[cid.Cid] tickets.TicketAck)
	accepts := make([]cid.Cid,0)
	rejects := make([]cid.Cid,0)
	prepare := make([]tickets.TicketAck,0)

    for _, ack := range acks {
		switch ack.ACK()  {
		case tickets.ACK_ACCEPT:
			accepts = append(accepts, ack.Cid())
		case tickets.ACK_CANCEL:
            rejects = append(rejects, ack.Cid())
		}
		ackMap[ack.Cid()] = ack
	}
    // handle reject Acks
    e.ticketStore.RemoveSendingTask(p, rejects)
	e.ticketStore.RemoveTickets(p, rejects)

    // handle accept Acks
	var activeEntries []peertask.Task
    blockSizes := e.bsm.getBlockSizes(ctx, accepts)
    msgSize := 0
	for _, c := range accepts {
		blockSize, ok := blockSizes[c]
		if ok {
			// I have corresponding block, send it
			if msgSize+blockSize > maxMessageSize {
				e.peerRequestQueue.PushBlock(p, activeEntries...)
				activeEntries = []peertask.Task{}
				msgSize = 0
			}
			activeEntries = append(activeEntries, peertask.Task{Identifier: ackMap[c].Cid(), Priority: 10})
			msgSize += blockSize
		} else {
			// I don't have corresponding block
            prepare = append(prepare, ackMap[c])
		}
	}
	if len(activeEntries) > 0 {
		e.peerRequestQueue.PushBlock(p, activeEntries...)
	}
	e.ticketStore.PrepareSending(prepare)
}

func (e *Engine) addBlocks(ks []cid.Cid) {
	work := false

	for _, l := range e.ledgerMap {
		l.lk.Lock()
		for _, k := range ks {
			if entry, ok := l.WantListContains(k); ok {
				e.peerRequestQueue.PushBlock(l.Partner, peertask.Task{
					Identifier: entry.Cid,
					Priority:   entry.Priority,
				})
				work = true
			}
		}
		l.lk.Unlock()
	}

    // send blocks if there is tickets
	acks, _ := e.ticketStore.PopSendingTasks(ks)
	for _, ack := range acks {
		e.peerRequestQueue.PushBlock(ack.Receiver(), peertask.Task{
            Identifier: ack.Cid(),
            Priority: 10000, // how to set the priority?
        })
        work = true
	}

	if work {
		e.signalNewWork()
	}
}

// AddBlocks is called when new blocks are received and added to a block store,
// meaning there may be peers who want those blocks, so we should send the blocks
// to them.
func (e *Engine) AddBlocks(ks []cid.Cid) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.addBlocks(ks)
}

// TODO add contents of m.WantList() to my local wantlist? NB: could introduce
// race conditions where I send a message, but MessageSent gets handled after
// MessageReceived. The information in the local wantlist could become
// inconsistent. Would need to ensure that Sends and acknowledgement of the
// send happen atomically

// MessageSent is called when a message has successfully been sent out, to record
// changes.
func (e *Engine) MessageSent(p peer.ID, m bsmsg.BitSwapMessage) {
	l := e.findOrCreate(p)
	l.lk.Lock()
	defer l.lk.Unlock()

	for _, block := range m.Blocks() {
		l.SentBytes(len(block.RawData()))
		l.wantList.Remove(block.Cid())
	}
}

// PeerConnected is called when a new peer connects, meaning we should start
// sending blocks.
func (e *Engine) PeerConnected(p peer.ID) {
	e.lock.Lock()
	defer e.lock.Unlock()
	l, ok := e.ledgerMap[p]
	if !ok {
		l = newLedger(p)
		e.ledgerMap[p] = l
	}
	l.lk.Lock()
	defer l.lk.Unlock()
	l.ref++
}

// PeerDisconnected is called when a peer disconnects.
func (e *Engine) PeerDisconnected(p peer.ID) {
	e.lock.Lock()
	defer e.lock.Unlock()
	l, ok := e.ledgerMap[p]
	if !ok {
		return
	}
	l.lk.Lock()
	defer l.lk.Unlock()
	l.ref--
	if l.ref <= 0 {
		delete(e.ledgerMap, p)
	}
}

func (e *Engine) numBytesSentTo(p peer.ID) uint64 {
	// NB not threadsafe
	return e.findOrCreate(p).Accounting.BytesSent
}

func (e *Engine) numBytesReceivedFrom(p peer.ID) uint64 {
	// NB not threadsafe
	return e.findOrCreate(p).Accounting.BytesRecv
}

// ledger lazily instantiates a ledger
func (e *Engine) findOrCreate(p peer.ID) *ledger {
	e.lock.Lock()
	defer e.lock.Unlock()
	l, ok := e.ledgerMap[p]
	if !ok {
		l = newLedger(p)
		e.ledgerMap[p] = l
	}
	return l
}

func (e *Engine) signalNewWork() {
	// Signal task generation to restart (if stopped!)
	// workSignal will wake up task worker
	select {
	case e.workSignal <- struct{}{}:
	default:
	}
}
