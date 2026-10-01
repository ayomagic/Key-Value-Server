package kvsrv

import (
	"crypto/rand"
	"log"
	"math/big"
	"sync"

	"6.5840/labrpc"
)

type Clerk struct {
	server   *labrpc.ClientEnd
	mu       sync.Mutex
	op       sync.Mutex
	opId     int64
	clientId int64
}

func (ck *Clerk) getOperationId() int64 {
	ck.op.Lock()
	defer ck.op.Unlock()

	ck.opId += 1
	return ck.opId
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	ck.opId = 0
	ck.clientId = nrand()
	// You'll have to add code here.
	// TODO: ???
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()

	// You will have to modify this function.
	args := GetArgs{
		Key:      key,
		OpId:     ck.getOperationId(),
		ClientId: ck.clientId,
	}
	reply := GetReply{
		Value: "",
	}
	ok := ck.server.Call("KVServer.Get", &args, &reply)
	if !ok {
		// TODO: do something (handle network failure)
	}
	if DEBUG {
		log.Printf("Clerk Get: key= %s, result= %s", key, reply.Value)
	}

	return reply.Value
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	ck.mu.Lock()
	defer ck.mu.Unlock()

	// You will have to modify this function.
	args := PutAppendArgs{
		Key:      key,
		Value:    value,
		OpId:     ck.getOperationId(),
		ClientId: ck.clientId,
	}
	reply := PutAppendReply{
		Value: "",
	}
	ok := ck.server.Call("KVServer."+op, &args, &reply)
	if !ok {
		// TODO: do something (handle network failure)
	}
	if DEBUG {
		log.Printf("Clerk call %s: key= %s, value= %s, result= %s", op, key, value, reply.Value)
	}

	return reply.Value
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
