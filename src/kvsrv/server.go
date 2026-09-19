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
	mu    sync.Mutex
	store map[string]string
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key := args.Key

	reply.Value = kv.store[key]
	if DEBUG {
		log.Printf("Get called: key= %s, value= %s", args.Key, reply.Value)
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key, value := args.Key, args.Value
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

	key, value := args.Key, args.Value
	if DEBUG {
		log.Printf("Append called: key= %s, value= %s", key, value)
	}

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

	if DEBUG {
		log.Printf("Append result: key= %s, value= %s, stored= %s", args.Key, reply.Value, kv.store[key])
	}

}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.store = make(map[string]string)

	return kv
}
