package kvsrv

import (
	"log"
	"os"
	"sync"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if DEBUG {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu          sync.Mutex
	store       map[string]string
	clientStore map[int64]Response
}

type Response struct {
	Value string
	opId  int64
}

func (kv *KVServer) lookupCacheResponse(clientId int64, operationId int64) (Response, bool) {
	response, ok := kv.clientStore[clientId]
	if ok && response.opId == operationId {
		return response, true
	}
	return Response{}, false
}

func (kv *KVServer) cacheResponse(clientId int64, data Response) {
	kv.clientStore[clientId] = data
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key, clientId := args.Key, args.ClientId
	delete(kv.clientStore, clientId)
	reply.Value = kv.store[key]
	if DEBUG {
		log.Printf("Get called: key= %s, value= %s", args.Key, reply.Value)
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key, value, clientId := args.Key, args.Value, args.ClientId
	delete(kv.clientStore, clientId)
	if DEBUG {
		log.Printf("Put called: key= %s, value= %s", key, value)
	}

	kv.store[key] = value
	reply.Value = kv.store[key]

	if DEBUG {
		log.Printf("Put result: key= %s, value= %s, stored= %s", key, value, kv.store[key])
	}

}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key, value, clienId, opId := args.Key, args.Value, args.ClientId, args.OpId
	if DEBUG {
		log.Printf("Append called: key= %s, value= %s", key, value)
	}

	// client response is cached
	if resp, ok := kv.lookupCacheResponse(clienId, opId); ok {
		reply.Value = resp.Value
	} else {
		// client response is not cached
		if _, ok := kv.store[key]; !ok {
			if DEBUG {
				log.Printf("key= %s not found", key)
			}
			kv.store[key] = value
			reply.Value = ""
		} else {
			if DEBUG {
				log.Printf("key= %s found", key)
			}
			reply.Value = kv.store[key]
			kv.store[key] += value
		}
		kv.cacheResponse(clienId, Response{Value: reply.Value, opId: opId})
	}

	if DEBUG {
		log.Printf("Append result: key= %s, value= %s, stored= %s", args.Key, reply.Value, kv.store[key])
	}

}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.store = make(map[string]string)
	kv.clientStore = make(map[int64]Response)
	return kv
}
