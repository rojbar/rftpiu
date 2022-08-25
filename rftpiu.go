package rftpiu

import (
	"bufio"
	"errors"
	"net"
	"regexp"
	"strings"
)

//OK
func ReadMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, errM := reader.ReadString(';')
	if errM != nil {
		return "", errM
	}
	return message, nil
}

//OK
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

//OK
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
