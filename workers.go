package bitswap

import (
	"context"

	engine "github.com/SJTU-OpenNetwork/go-bitswap/decision"
	bsmsg "github.com/SJTU-OpenNetwork/go-bitswap/message"
	cid "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	process "github.com/jbenet/goprocess"
	procctx "github.com/jbenet/goprocess/context"
)

// TaskWorkerCount is the total number of simultaneous threads sending
// outgoing messages
var TaskWorkerCount = 4

// TicketWorkerCount is the total number of si,ultaneous threads sending outgoing messages
// TODO: Judge the number of ticketworkers - Riften
var TicketWorkerCount = 4

func (bs *Bitswap) startWorkers(ctx context.Context, px process.Process) {

	// Start up workers to handle requests from other nodes for the data on this node
	for i := 0; i < TaskWorkerCount; i++ {
		i := i
		px.Go(func(px process.Process) {
			bs.taskWorker(ctx, i)
		})
	}

	if bs.provideEnabled {
		// Start up a worker to manage sending out provides messages
		px.Go(func(px process.Process) {
			bs.provideCollector(ctx)
		})

		// Spawn up multiple workers to handle incoming blocks
		// consider increasing number if providing blocks bottlenecks
		// file transfers
		px.Go(bs.provideWorker)
	}

	// TODO: start ticketWorker - Jerry
//	for i := 0; i < TicketWorkerCount; i++ {
//		i := i
//		px.Go(func(px process.Process) {
//			bs.ticketWorker(ctx, i)
//		})
//	}
}

func (bs *Bitswap) taskWorker(ctx context.Context, id int) {
	idmap := logging.LoggableMap{"ID": id}
	defer log.Debug("bitswap task worker shutting down...")
	for {
		log.Event(ctx, "Bitswap.TaskWorker.Loop", idmap)
		select {
		//case ticketEnvelope := <- bs.engine.TicketOutbox():
		//	bs.sendTicketsOrAck(ctx, ticketEnvelope)
		//case ticketAckEnvelope := <- bs.engine.TicketAckOutbox():
		//	bs.sendTicketsOrAck(ctx, ticketAckEnvelope)
		case nextEnvelope := <-bs.engine.Outbox():
			select {
			case envelope, ok := <-nextEnvelope:
				if !ok {
					continue
				}
                bs.counterLk.Lock()
                *bs.wcounter ++
                bs.counterLk.Unlock()
				// update the BS ledger to reflect sent message
				// TODO: Should only track *useful* messages in ledger
				outgoing := bsmsg.New(false)
				for _, block := range envelope.Message.Blocks() {
					log.Event(ctx, "Bitswap.TaskWorker.Work", logging.LoggableF(func() map[string]interface{} {
						return logging.LoggableMap{
							"ID":     id,
							"Target": envelope.Peer.Pretty(),
							"Block":  block.Cid().String(),
						}
					}))
					outgoing.AddBlock(block)
				}
				bs.engine.MessageSent(envelope.Peer, outgoing)

				bs.sendBlocks(ctx, envelope)
				bs.counterLk.Lock()
				for _, block := range envelope.Message.Blocks() {
					bs.counters.blocksSent++
					bs.counters.dataSent += uint64(len(block.RawData()))
				}
                *bs.wcounter --
				bs.counterLk.Unlock()

			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}


// TODO: ticketWorker is used to send ticket or ticket ack - Jerry
func (bs *Bitswap) ticketWorker(ctx context.Context, id int) {
	idmap := logging.LoggableMap{"ID": id}
	defer log.Debug("bitswap ticket worker shutting down...")
	for {
		log.Event(ctx, "Bitswap.TicketWorker.Loop", idmap)
		select{
		case ticketEnvelope := <- bs.engine.TicketOutbox():
			bs.sendTicketsOrAck(ctx, ticketEnvelope)
		case ticketAckEnvelope := <- bs.engine.TicketAckOutbox():
			bs.sendTicketsOrAck(ctx, ticketAckEnvelope)
		case <-ctx.Done():
			return
		}
	}
}

// TODO: add sendTicket - Jerry
//func (bs *Bitswap) sendTicket(ctx context.Context, env *engine.Envelope) {
//	defer env.Sent()
//
//	msg := bsmsg.New(false)
//	for _, ticket := range env.Message.Tickets(){
//		msg.AddTicket(ticket)
//	}
//	err := bs.network.SendMessage(ctx, env.Peer, msg)
//	if err != nil {
//		log.Infof("sendticket error :%s", err)
//	}
//}
//
//// TODO: add sendTicketAck - Jerry
//func (bs *Bitswap) sendTicketAck(ctx context.Context, env *engine.Envelope) {
//	defer env.Sent()
//
//	msg := bsmsg.New(false)
//	for _, ack := range env.Message.TicketAcks(){
//		msg.AddTicketAck(ack)
//	}
//	err := bs.network.SendMessage(ctx, env.Peer, msg)
//	if err != nil {
//		log.Infof("sendticketack error :%s", err)
//	}
//
//}

func (bs *Bitswap) sendBlocks(ctx context.Context, env *engine.Envelope) {
	// Blocks need to be sent synchronously to maintain proper backpressure
	// throughout the network stack
	defer env.Sent()

	msgSize := 0
	msg := bsmsg.New(false)
	for _, block := range env.Message.Blocks() {
		msgSize += len(block.RawData())
		msg.AddBlock(block)
		//log.Infof("Sending block %s to %s", block, env.Peer)
		log.Debugf("[BLKSEND] Cid %s, SendTo %s", block.Cid().String(), env.Peer)
	}

	bs.sentHistogram.Observe(float64(msgSize))
	err := bs.network.SendMessage(ctx, env.Peer, msg)
	if err != nil {
		//log.Infof("sendblock error: %s", err)
		log.Error(err)
	}
	for _, block := range env.Message.Blocks() {
		msgSize += len(block.RawData())
		msg.AddBlock(block)
		//log.Infof("Sending block %s to %s", block, env.Peer)
		log.Debugf("[BLKSENDED] Cid %s, SendTo %s", block.Cid().String(), env.Peer)
	}
}

func (bs *Bitswap) sendTicketsOrAck(ctx context.Context, env *engine.Envelope) {
	defer env.Sent()
	log.Debugf("Begin to send message to %s", env.Peer.String())
	err := bs.network.SendMessage(ctx, env.Peer, env.Message)
	if err != nil{
		log.Infof("sendTickets errorL %s", err)
	}
}

//func (bs *Bitswap) sendTicketAcks()


func (bs *Bitswap) provideWorker(px process.Process) {
	// FIXME: OnClosingContext returns a _custom_ context type.
	// Unfortunately, deriving a new cancelable context from this custom
	// type fires off a goroutine. To work around this, we create a single
	// cancelable context up-front and derive all sub-contexts from that.
	//
	// See: https://github.com/ipfs/go-ipfs/issues/5810
	ctx := procctx.OnClosingContext(px)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	limit := make(chan struct{}, provideWorkerMax)

	limitedGoProvide := func(k cid.Cid, wid int) {
		defer func() {
			// replace token when done
			<-limit
		}()
		ev := logging.LoggableMap{"ID": wid}

		defer log.EventBegin(ctx, "Bitswap.ProvideWorker.Work", ev, k).Done()

		ctx, cancel := context.WithTimeout(ctx, provideTimeout) // timeout ctx
		defer cancel()

		if err := bs.network.Provide(ctx, k); err != nil {
			log.Warning(err)
		}
	}

	// worker spawner, reads from bs.provideKeys until it closes, spawning a
	// _ratelimited_ number of workers to handle each key.
	for wid := 2; ; wid++ {
		ev := logging.LoggableMap{"ID": 1}
		log.Event(ctx, "Bitswap.ProvideWorker.Loop", ev)

		select {
		case <-px.Closing():
			return
		case k, ok := <-bs.provideKeys:
			if !ok {
				log.Debug("provideKeys channel closed")
				return
			}
			select {
			case <-px.Closing():
				return
			case limit <- struct{}{}:
				go limitedGoProvide(k, wid)
			}
		}
	}
}

func (bs *Bitswap) provideCollector(ctx context.Context) {
	defer close(bs.provideKeys)
	var toProvide []cid.Cid
	var nextKey cid.Cid
	var keysOut chan cid.Cid

	for {
		select {
		case blkey, ok := <-bs.newBlocks:
			if !ok {
				log.Debug("newBlocks channel closed")
				return
			}

			if keysOut == nil {
				nextKey = blkey
				keysOut = bs.provideKeys
			} else {
				toProvide = append(toProvide, blkey)
			}
		case keysOut <- nextKey:
			if len(toProvide) > 0 {
				nextKey = toProvide[0]
				toProvide = toProvide[1:]
			} else {
				keysOut = nil
			}
		case <-ctx.Done():
			return
		}
	}
}
