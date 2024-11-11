// Package rob implemented an reorder buffer for memory requests.
package rob

import (
	"container/list"

	"github.com/sarchlab/akita/v3/mem/mem"
	"github.com/sarchlab/akita/v3/sim"
	"github.com/sarchlab/akita/v3/tracing"
)

// 不需要记录向core的RSP，因为这意味着transaction的完成
type transaction struct {
	reqFromTop    mem.AccessReq
	reqToBottom   mem.AccessReq
	rspFromBottom mem.AccessRsp
}

// ReorderBuffer can maintain the returning order of memory transactions.
type ReorderBuffer struct {
	// 嵌入是将一个结构体作为另一个结构体的字段，不需要显式地命名字段
	// ROB直接嵌入TickingComponent，继承其所有字段和方法
	*sim.TickingComponent

	// ROB拥有的端口
	topPort     sim.Port
	bottomPort  sim.Port
	controlPort sim.Port

	// 不属于ROB的Port，属于另一个组件，作为请求发送的目的地
	BottomUnit sim.Port

	bufferSize     int
	numReqPerCycle int

	// list.Element表示链表中的一个节点。每个 list.Element 包含一个值和指向前后元素的指针
	// toBottomReqIDToTransactionTable 可以用来实现高效的查找和更新功能
	// 映射 req.id --> 发送给BottomUnit的req对应的事务
	toBottomReqIDToTransactionTable map[string]*list.Element

	// *list.List 是一个指向链表的指针，记录当前ROB缓存的正在进行的所有事务，用来实现高效增加和删除事务
	// 链表中的元素通过Value取出，其中类型为*transaction 类型的指针
	transactions *list.List
	isFlushing   bool
}

// Tick updates the status of the ReorderBuffer.
// Tick()是一个interface函数
// 返回True说明当前周期组件可以取得进步，可以调度新的event
func (b *ReorderBuffer) Tick(now sim.VTimeInSec) (madeProgress bool) {
	madeProgress = b.processControlMsg(now) || madeProgress

	if !b.isFlushing {
		madeProgress = b.runPipeline(now) || madeProgress
	}

	return madeProgress
}

func (b *ReorderBuffer) processControlMsg(
	now sim.VTimeInSec,
) (madeProgress bool) {
	item := b.controlPort.Peek()
	if item == nil {
		return false
	}

	msg := item.(*mem.ControlMsg)
	if msg.DiscardTransations {
		return b.discardTransactions(now, msg)
	} else if msg.Restart {
		return b.restart(now, msg)
	}

	panic("never")
}

func (b *ReorderBuffer) discardTransactions(
	now sim.VTimeInSec,
	msg *mem.ControlMsg,
) (madeProgress bool) {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(b.controlPort).
		WithDst(msg.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := b.controlPort.Send(rsp)
	if err != nil {
		return false
	}

	b.isFlushing = true
	b.toBottomReqIDToTransactionTable = make(map[string]*list.Element)
	b.transactions.Init()
	b.controlPort.Retrieve(now)

	// fmt.Printf("%.10f, %s, rob flushed\n", now, b.Name())

	return true
}

func (b *ReorderBuffer) restart(
	now sim.VTimeInSec,
	msg *mem.ControlMsg,
) (madeProgress bool) {
	rsp := mem.ControlMsgBuilder{}.
		WithSrc(b.controlPort).
		WithDst(msg.Src).
		WithSendTime(now).
		ToNotifyDone().
		Build()

	err := b.controlPort.Send(rsp)
	if err != nil {
		return false
	}

	b.isFlushing = false
	b.toBottomReqIDToTransactionTable = make(map[string]*list.Element)
	b.transactions.Init()

	for b.topPort.Retrieve(now) != nil {
	}

	for b.bottomPort.Retrieve(now) != nil {
	}

	b.controlPort.Retrieve(now)

	// fmt.Printf("%.10f, %s, rob restarted\n", now, b.Name())

	return true
}

// bottomUp：先看transaction buffer的头部，取出front并且已经得到rsp的事务，从buffer中删除事务并且通过topPort回复response
// parseBottom：根据bottomPort的响应message修改对应事务的状态
// topDown：从topPort中的message转发到bottomPort，并且在transaction buffer中新建事务记录
func (b *ReorderBuffer) runPipeline(now sim.VTimeInSec) (madeProgress bool) {
	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.bottomUp(now) || madeProgress
	}

	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.parseBottom(now) || madeProgress
	}

	for i := 0; i < b.numReqPerCycle; i++ {
		madeProgress = b.topDown(now) || madeProgress
	}

	return madeProgress
}

// 接收来自core的请求，把请求插入buffer，然后把请求转发到底部端口
func (b *ReorderBuffer) topDown(now sim.VTimeInSec) bool {

	if b.isFull() {
		return false
	}

	// 从topPort取出一个request
	item := b.topPort.Peek()
	if item == nil {
		return false
	}

	req := item.(mem.AccessReq)

	// 形成该request对应的transaction
	trans := b.createTransaction(req)

	// reqToBottom的src是b.bottomPort，dst是其他组件的BottomUnit
	trans.reqToBottom.Meta().Src = b.bottomPort
	trans.reqToBottom.Meta().SendTime = now

	// 发送该request到BottomUnit
	err := b.bottomPort.Send(trans.reqToBottom)
	if err != nil {
		return false
	}

	// 发送成功则把该事务加入buffer
	b.addTransaction(trans)
	// 从topPort的buffer中移除该message
	b.topPort.Retrieve(now)

	tracing.TraceReqReceive(req, b)
	tracing.TraceReqInitiate(trans.reqToBottom, b,
		tracing.MsgIDAtReceiver(req, b))

	return true
}

// 接收来自底部端口的rsp，匹配并更新buffer中对应的事务状态
func (b *ReorderBuffer) parseBottom(now sim.VTimeInSec) bool {
	item := b.bottomPort.Peek()
	if item == nil {
		return false
	}

	rsp := item.(mem.AccessRsp)
	// 找到需要响应的request id
	rspTo := rsp.GetRspTo()

	// 找到原request对应的事务
	transElement, found := b.toBottomReqIDToTransactionTable[rspTo]

	// 更新事务对应的状态
	if found {
		trans := transElement.Value.(*transaction)
		trans.rspFromBottom = rsp

		tracing.TraceReqFinalize(trans.reqToBottom, b)
	}

	b.bottomPort.Retrieve(now)

	return true
}

// 始终查看buffer的头部的事务，观察其是否已经完成
func (b *ReorderBuffer) bottomUp(now sim.VTimeInSec) bool {
	elem := b.transactions.Front()
	if elem == nil {
		return false
	}

	// 未完成返回false
	trans := elem.Value.(*transaction)
	if trans.rspFromBottom == nil {
		return false
	}

	// 完成后通过TopPort向请求者回复rsp
	rsp := b.duplicateRsp(trans.rspFromBottom, trans.reqFromTop.Meta().ID)
	// rsp的src为topPort，dst为req的src
	rsp.Meta().Dst = trans.reqFromTop.Meta().Src
	rsp.Meta().Src = b.topPort
	rsp.Meta().SendTime = now

	// 通过topPort回复request
	err := b.topPort.Send(rsp)
	if err != nil {
		return false
	}

	b.deleteTransaction(elem)

	tracing.TraceReqComplete(trans.reqFromTop, b)

	return true
}

// ROB维护的事务已满
func (b *ReorderBuffer) isFull() bool {
	return b.transactions.Len() >= b.bufferSize
}

// 形成该request对应的transaction，即把请求转发到bottom
func (b *ReorderBuffer) createTransaction(req mem.AccessReq) *transaction {
	return &transaction{
		reqFromTop:  req,
		reqToBottom: b.duplicateReq(req),
	}
}

// 把transaction记录在ROB的buffer中
func (b *ReorderBuffer) addTransaction(trans *transaction) {
	elem := b.transactions.PushBack(trans)
	b.toBottomReqIDToTransactionTable[trans.reqToBottom.Meta().ID] = elem
}

// 从map和链表中移除对应的事务
func (b *ReorderBuffer) deleteTransaction(elem *list.Element) {
	trans := elem.Value.(*transaction)
	b.transactions.Remove(elem)
	delete(b.toBottomReqIDToTransactionTable, trans.reqToBottom.Meta().ID)
}

// 用于转发请求，每一个req的id都不一样
func (b *ReorderBuffer) duplicateReq(req mem.AccessReq) mem.AccessReq {
	switch req := req.(type) {
	case *mem.ReadReq:
		return b.duplicateReadReq(req)
	case *mem.WriteReq:
		return b.duplicateWriteReq(req)
	default:
		panic("unsupported type")
	}
}

// 转发read request，把request目的地改为BottomUnit
func (b *ReorderBuffer) duplicateReadReq(req *mem.ReadReq) *mem.ReadReq {
	return mem.ReadReqBuilder{}.
		WithAddress(req.Address).
		WithByteSize(req.AccessByteSize).
		WithPID(req.PID).
		WithDst(b.BottomUnit).
		Build()
}

// 转发write request，把request的目的改为BottomUnit
func (b *ReorderBuffer) duplicateWriteReq(req *mem.WriteReq) *mem.WriteReq {
	return mem.WriteReqBuilder{}.
		WithAddress(req.Address).
		WithPID(req.PID).
		WithData(req.Data).
		WithDirtyMask(req.DirtyMask).
		WithDst(b.BottomUnit).
		Build()
}

// 转发bottomPort的response到topPort
// rsp mem.AccessRsp：是bottomPort端口得到的响应
// rspTo string：是这个rsp对应的事务最开始的request的id
func (b *ReorderBuffer) duplicateRsp(
	rsp mem.AccessRsp,
	rspTo string,
) mem.AccessRsp {
	switch rsp := rsp.(type) {
	case *mem.DataReadyRsp:
		return b.duplicateDataReadyRsp(rsp, rspTo)
	case *mem.WriteDoneRsp:
		return b.duplicateWriteDoneRsp(rsp, rspTo)
	default:
		panic("type not supported")
	}
}

// 回复数据
// rsp mem.AccessRsp：是bottomPort端口得到的响应
// rspTo string：是这个rsp对应的事务最开始的request的id
func (b *ReorderBuffer) duplicateDataReadyRsp(
	rsp *mem.DataReadyRsp,
	rspTo string,
) *mem.DataReadyRsp {
	return mem.DataReadyRspBuilder{}.
		WithData(rsp.Data).
		WithRspTo(rspTo).
		Build()
}

// 回复写入完成的ack
func (b *ReorderBuffer) duplicateWriteDoneRsp(
	rsp *mem.WriteDoneRsp,
	rspTo string,
) *mem.WriteDoneRsp {
	return mem.WriteDoneRspBuilder{}.
		WithRspTo(rspTo).
		Build()
}
