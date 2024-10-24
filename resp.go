package main

import (
	"net"
	"strconv"
	"errors"
	"io"
	"bufio"
)

// Constants representing the RESP types
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value struct to represent RESP types
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// Marshal method to convert Value to RESP format
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// marshalString converts a simple string value to RESP format
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalBulk converts a bulk string value to RESP format
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalArray converts an array value to RESP format
func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}

// marshallError converts an error value to RESP format
func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshallNull converts a null value to RESP format
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

// NewResp initializes a new RESP connection
type NewResp struct {
	conn net.Conn
}

// Read reads a RESP message from the connection
func (r *NewResp) Read() (Value, error) {
	// Read the first byte to determine the type
	var buf [1]byte
	_, err := r.conn.Read(buf[:])
	if err != nil {
		return Value{}, err
	}

	typ := buf[0]

	switch typ {
	case STRING:
		return r.readString()
	case ERROR:
		return r.readError()
	case INTEGER:
		return r.readInteger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return Value{}, errors.New("unknown RESP type")
	}
}

// readString reads a string from the RESP format
func (r *NewResp) readString() (Value, error) {
	return r.readValue(STRING)
}

// readError reads an error from the RESP format
func (r *NewResp) readError() (Value, error) {
	return r.readValue(ERROR)
}

// readInteger reads an integer from the RESP format
func (r *NewResp) readInteger() (Value, error) {
	return r.readValue(INTEGER)
}

// readBulk reads a bulk string from the RESP format
func (r *NewResp) readBulk() (Value, error) {
	var lengthBuf [32]byte

	// Read length of bulk string
	n, err := r.conn.Read(lengthBuf[:])
	if err != nil {
		return Value{}, err
	}

	lengthStr := string(lengthBuf[:n])
	length, err := strconv.Atoi(lengthStr[:n-2]) // Remove CRLF
	if err != nil {
		return Value{}, err
	}

	if length == -1 {
		return Value{typ: "null"}, nil // Handle null bulk string
	}

	// Read the bulk string
	bulkData := make([]byte, length+2) // +2 for CRLF
	_, err = io.ReadFull(r.conn, bulkData)
	if err != nil {
		return Value{}, err
	}

	// Return the bulk string
	return Value{typ: "bulk", bulk: string(bulkData[:length])}, nil
}

// readArray reads an array from the RESP format
func (r *NewResp) readArray() (Value, error) {
	var lengthBuf [32]byte

	// Read length of array
	n, err := r.conn.Read(lengthBuf[:])
	if err != nil {
		return Value{}, err
	}

	lengthStr := string(lengthBuf[:n])
	length, err := strconv.Atoi(lengthStr[:n-2]) // Remove CRLF
	if err != nil {
		return Value{}, err
	}

	array := make([]Value, length)
	for i := 0; i < length; i++ {
		value, err := r.Read() // Read each value in the array
		if err != nil {
			return Value{}, err
		}
		array[i] = value
	}

	return Value{typ: "array", array: array}, nil
}

// readValue reads a RESP value based on the provided type
func (r *NewResp) readValue(typ byte) (Value, error) {
	// Read until CRLF and parse the value
	buf := bufio.NewReader(r.conn)
	valueBytes, err := buf.ReadBytes('\n')
	if err != nil {
		return Value{}, err
	}

	valueStr := string(valueBytes[:len(valueBytes)-2]) // Remove CRLF

	switch typ {
	case STRING:
		return Value{typ: "string", str: valueStr}, nil
	case ERROR:
		return Value{typ: "error", str: valueStr}, nil
	case INTEGER:
		num, err := strconv.Atoi(valueStr)
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "integer", num: num}, nil
	default:
		return Value{}, errors.New("unknown type")
	}
}
