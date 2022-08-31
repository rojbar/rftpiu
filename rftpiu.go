package rftpiu

import (
	"bufio"
	"errors"
	"io"
	"net"
	"regexp"
	"strings"

	"go.uber.org/zap"
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

/**
	Reads to the buffer and writes from the buffer
**/
func ReadThenWrite(reader io.Reader, writer bufio.Writer, buffer []byte) error {
	_, errR := io.ReadFull(reader, buffer)
	if errR != nil {
		if errR == io.EOF {
			print(errR)
		}
		return errR
	}
	_, errW := writer.Write(buffer)
	if errW != nil {
		return errW
	}
	errF := writer.Flush()
	if errF != nil {
		return errF
	}

	return nil
}

type Queue[T any] struct {
	Data []T
	Head int
	Tail int
	Size int
}

func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{Data: make([]T, size+1), Head: 0, Tail: 0, Size: size}
}

func (queue *Queue[T]) Enqueue(data T) error {
	if queue.Tail+1 == queue.Head || (queue.Tail == queue.Size && queue.Head == 0) {
		return errors.New("trying to enqueue full queue")
	}

	queue.Data[queue.Tail] = data
	if queue.Tail == queue.Size {
		queue.Tail = 0
	} else {
		queue.Tail += 1
	}

	return nil
}

func (queue *Queue[T]) Retrieve() (T, error) {
	if queue.Head == queue.Tail {
		var aux T
		return aux, errors.New("trying to retrieve empty queue")
	}
	data := queue.Data[queue.Head]

	return data, nil
}

func (queue *Queue[T]) Dequeue() (T, error) {
	if queue.Head == queue.Tail {
		var aux T
		return aux, errors.New("trying to dequeue empty queue")
	}

	data := queue.Data[queue.Head]

	if queue.Head == queue.Size {
		queue.Head = 0
	} else {
		queue.Head += 1
	}

	return data, nil
}

var Logger *zap.Logger

func InitializeLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

}
