package account

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type LenOpDataSocket struct {
	Conn net.Conn
}

func (socket LenOpDataSocket) Send(op uint16, data []byte) error {
	//Буффер содержащий пакет
	var buf bytes.Buffer

	//Вычислим длинну пакета
	packetLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(packetLen, uint16(2+2+len(data)))

	//Получим опкод
	opcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, op)

	buf.Write(packetLen) //длинна
	buf.Write(opcode)    //опкод
	buf.Write(data)      //данные

	send := buf.Bytes()

	i, err := socket.Conn.Write(buf.Bytes())
	if err != nil {
		return err
	}

	log.Printf("Send %v bytes: %x\n", i, send)
	return nil
}

func (socket LenOpDataSocket) Read() (uint16, []byte, error) {
	//Прочитаем данные
	buf := make([]byte, 64)
	_, err := socket.Conn.Read(buf)
	if err != nil {
		return 0, nil, err
	}

	if len(buf) < 4 {
		return 0, nil, fmt.Errorf("real len < 4")
	}

	//Вычислим длинну
	packetLen := binary.LittleEndian.Uint16(buf[:2])

	if len(buf) < int(packetLen) {
		return 0, nil, fmt.Errorf("real len < packet len")
	}

	//Готово
	log.Printf("Recv %v bytes: %x\n", packetLen, buf[:packetLen])
	return binary.BigEndian.Uint16(buf[2:4]), buf[4:packetLen], nil
}
