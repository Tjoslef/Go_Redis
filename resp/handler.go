package resp

var Handler = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
	return Value{Typ: "string", Str: "PONG"}
}
