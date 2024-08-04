package resp

import "sync"

var Handler = map[string]func([]Value) Value{
	"PING": ping,
	"GET":  get,
	"SET":  set,
	"HGET": hget,
	"HSET": hset,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}
	return Value{Typ: "string", Str: args[0].Bulk}

}

var SETs = map[string]string{}
var SETmus = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) > 2 {
		return Value{Typ: "error", Str: "to many argument for SET"}
	}
	key := args[0].Bulk
	value := args[1].Bulk
	SETmus.Lock()
	SETs[key] = value
	SETmus.Unlock()
	return Value{Typ: "string", Str: "OK"}
}
func get(args []Value) Value {

	if len(args) > 1 {
		return Value{Typ: "error", Str: "to many argument for GET"}
	}
	key := args[0].Bulk
	SETmus.RLock()
	value, ok := SETs[key]
	if !ok {
		return Value{Typ: "null"}
	}
	return Value{Typ: "bulk", Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Typ: "string", Str: "wrong number of argument"}
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{Typ: "string", Str: "OK"}

}
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Typ: "string", Str: "wrong number of argument"}
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	if !ok {
		return Value{Typ: "null"}
	}
	return Value{Typ: "bulk", Str: value}
}
