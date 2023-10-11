package SocketIO

import (
	"bytes"

	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

const (
	EVENT        = 2
	BINARY_EVENT = 5
)

type EmitterOpts struct {
	// Host means hostname like localhost
	Host string
	// Port means port number, like 6379
	Port int
	// Key means redis subscribe key
	Key string
	// Protocol, like tcp
	Protocol string
	// Address, like localhost:6379
	Addr     string
	Password string
	DB       int
}

type Emitter struct {
	Redis *redis.Client
	Key   string
	rooms []string
	flags map[string]interface{}
}

func NewEmitter(conn *redis.Client) (*Emitter, error) {
	emitter := &Emitter{
		Redis: conn,
		Key:   "socket.io#/#",
	}
	return emitter, nil
}

func (emitter *Emitter) Join() *Emitter {
	emitter.flags["join"] = true
	return emitter
}

func (emitter *Emitter) Volatile() *Emitter {
	emitter.flags["volatile"] = true
	return emitter
}

func (emitter *Emitter) Broadcast() *Emitter {
	emitter.flags["broadcast"] = true
	return emitter
}

/**
 * Limit emission to a certain `room`.
 *
 * @param {String} room
 */
func (emitter *Emitter) In(room string) *Emitter {
	for _, r := range emitter.rooms {
		if r == room {
			return emitter
		}
	}
	emitter.rooms = append(emitter.rooms, room)
	return emitter
}

func (emitter *Emitter) To(room string) *Emitter {
	return emitter.In(room)
}

func (emitter *Emitter) ToRooms(rooms []string) *Emitter {
	var temp []string
	// 合并rooms和 emitter.rooms, 并去重
	for _, r := range rooms {
		temp = append(temp, r)
	}
	for _, r := range emitter.rooms {
		temp = append(temp, r)
	}
	var result []string
	for _, r := range temp {
		var flag bool
		for _, rr := range result {
			if r == rr {
				flag = true
			}
		}
		if !flag {
			result = append(result, r)
		}
	}

	emitter.rooms = result
	return emitter
}

/**
 * Limit emission to certain `namespace`.
 *
 * @param {String} namespace
 */
func (emitter *Emitter) Of(namespace string) *Emitter {
	emitter.flags["nsp"] = namespace
	return emitter
}

// send the packet by string, json, etc
// Usage:
// Emit("event name", "data")
func (emitter *Emitter) Emit(event string, data ...interface{}) (*Emitter, error) {
	d := []interface{}{event}
	d = append(d, data...)
	packet := map[string]interface{}{
		"type": EVENT,
		"data": d,
	}
	return emitter.emit(packet)
}

// send the packet by binary
// Usage:
// EmitBinary("event name", []byte{0x01, 0x02, 0x03})
func (emitter *Emitter) EmitBinary(event string, data ...interface{}) (*Emitter, error) {
	d := []interface{}{event}
	d = append(d, data...)
	packet := map[string]interface{}{
		"type": BINARY_EVENT,
		"data": d,
	}
	return emitter.emit(packet)
}

func (emitter *Emitter) emit(packet map[string]interface{}) (*Emitter, error) {
	if emitter.flags["nsp"] != nil {
		packet["nsp"] = emitter.flags["nsp"]
		delete(emitter.flags, "nsp")
	} else {
		packet["nsp"] = "/"
	}
	var pack []interface{} = make([]interface{}, 0)
	pack = append(pack, "uid")
	pack = append(pack, packet)
	// interface BroadcastOptions {
	//     rooms: Set<Room>;
	//     except?: Set<Room>;
	//     flags?: BroadcastFlags;
	// }
	pack = append(pack, map[string]interface{}{
		"rooms": emitter.rooms,
		"flags": emitter.flags,
	})
	if len(emitter.rooms) == 1 {
		emitter.Key = emitter.Key + emitter.rooms[0] + "#"
	}
	buf := &bytes.Buffer{}
	enc := msgpack.NewEncoder(buf)
	error := enc.Encode(pack)
	if error != nil {
		return nil, error
	}
	// binary hack for socket.io-parser
	if packet["type"] == BINARY_EVENT {
		expectedBytes := []byte{byte(0xd8), byte(0x00), byte(0x09)}
		replaceBytes := []byte{byte(0xa9)}
		buf = bytes.NewBuffer(bytes.Replace(buf.Bytes(), replaceBytes, expectedBytes, 1))
	}

	emitter.Redis.Publish(emitter.Key, buf.String())
	emitter.rooms = make([]string, 0, 0)
	emitter.flags = make(map[string]interface{})
	return emitter, nil
}

func (emitter *Emitter) Close() {
	if emitter.Redis != nil {
		defer emitter.Redis.Close()
	}
}
