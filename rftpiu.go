package rftpiu

import (
	"bufio"
	"errors"
	"net"
	"regexp"
	"strings"
)

/**
	ReadMessage returns the rftp message send by the other host. Returns error in
	case of invalid message or conn.Read failing
**/
func ReadMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, errM := reader.ReadString(';')
	if errM != nil {
		return "", errM
	}
	return message, nil
}

/**
	SendMessage sends a rftp message to a tcp connection, it returns error in case of
	conn.Write failing
**/
func SendMessage(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	_, errW := writer.WriteString(message)
	if errW != nil {
		return errW
	}
	errF := writer.Flush()
	if errF != nil {
		return errF
	}
	return nil
}

/**
	GetKey returns the value of a key in a rftp message, it returns error when key not found
	or value not found for provided key, if error != nil string = ""
**/
func GetKey(message string, key string) (string, error) {
	regExp, errReg := regexp.Compile(key + ":" + " ([a-z]|[A-Z]|[0-9])+")
	if errReg != nil {
		return "", nil
	}
	result := regExp.Find([]byte(message))
	if result == nil {
		return "", errors.New("key not found")
	}

	aux := string(result)
	_, after, found := strings.Cut(aux, ":")
	if !found {
		return "", errors.New("value not found for provided key")
	}
	after = strings.TrimSpace(after)
	return after, nil
}
