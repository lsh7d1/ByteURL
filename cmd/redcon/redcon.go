package main

import (
	"log"
	"strings"
	"time"

	"byteurl/cache"

	"github.com/tidwall/redcon"
)

var addr string = ":6380"

func main() {
	c := cache.NewCache("test", time.Minute, cache.WithAroundCapLimit(11919810))
	log.Printf("Started Server at: %s", addr)
	err := redcon.ListenAndServe(
		addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
			case "ping":
				conn.WriteString("PONG")
			case "set":
				if len(cmd.Args) != 3 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				_, _ = c.Take(string(cmd.Args[1]), func() (string, error) {
					return string(cmd.Args[2]), nil
				})
				// c.Set(string(cmd.Args[1]), string(cmd.Args[2]))
				conn.WriteString("OK")
			case "get":
				if len(cmd.Args) != 2 {
					conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
					return
				}
				val, ok := c.Get(string(cmd.Args[1]))
				if !ok {
					conn.WriteNull()
				} else {
					conn.WriteBulkString(val)
				}
				// case "del":
				// TODO: Cache's Del has no return value, and the redis client needs the server's return value.
			}
		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)
	if err != nil {
		log.Fatalf("redcon.ListenAndServe failed, err: %#v\n", err)
	}
}
