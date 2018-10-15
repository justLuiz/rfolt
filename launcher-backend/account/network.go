package account

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"rfolt/configuration"
	"rfolt/launcher-backend/client"
	"time"
)

const (
	StopOnCheck      int = 0
	StopOnServerList int = 1
)

const (
	_Opcode_Idle              uint16 = 0x6501
	_Opcode_Out_Handshake     uint16 = 0x150C
	_Opcode_In_EncKeys        uint16 = 0x150D
	_Opcode_Out_Elp           uint16 = 0x1503
	_Opcode_In_AccInfo        uint16 = 0x1504
	_Opcode_Out_GetServerList uint16 = 0x1505
	_Opcode_In_ServerList     uint16 = 0x1506
	_Opcode_Out_ConnectTo     uint16 = 0x1507
	_Opcode_In_ServerInfo     uint16 = 0x1508
)

type WorkResult struct {
	NetworkState  int
	AccountStatus int
	ServerList    []string
	DefaultSet    DefaultSet
	BinPid        int
}

func loginServerWork(cfg configuration.Network, login, pwd string, stopOn int) (result WorkResult) {
	socket := LenOpDataSocket{}
	var err error
	dialer := net.Dialer{Timeout: 3 * time.Second}
	socket.Conn, err = dialer.Dial("tcp", fmt.Sprintf("%v:%v", cfg.Ip, cfg.LoginPort))
	if err != nil {
		log.Printf("Error with connecting to login server: %v", err)
		result.NetworkState = -1
		return
	}
	defer socket.Conn.Close()

	//Отправим запрос на ключи шифрования
	socket.Send(_Opcode_Out_Handshake, make([]byte, 1))

	for socket.Conn != nil {
		opcode, data, err := socket.Read()
		if err != nil {
			log.Printf("Error with reading login server: %v", err)
			result.NetworkState = -1
			return
		}

		switch opcode {
		case _Opcode_Idle: //пришел idle-пакет
			socket.Send(_Opcode_Idle, make([]byte, 1)) //отправим назад такой же
		case _Opcode_In_EncKeys: //пришли ключи для шифрования логина и пароля
			plusKey := data[0] + cfg.PlusSalt
			xorKey := data[1] + cfg.XorSalt

			var buf bytes.Buffer   //буфер
			buf.WriteString(login) //логин и нули
			buf.Write(make([]byte, 13-len(login)))
			buf.WriteString(pwd) //пароль и нули
			buf.Write(make([]byte, 13-len(pwd)))
			buf.WriteByte(0x00) //завершающий ноль
			send := buf.Bytes()

			//Зашифруем
			for i := 0; i < len(send)-1; i++ {
				send[i] = (send[i] + plusKey) ^ xorKey
			}

			socket.Send(_Opcode_Out_Elp, send)
		case _Opcode_In_AccInfo:
			result.AccountStatus = int(data[0]) //получим состояние аккаунта

			//если хотели только проверить состояние аккаунта или аккаунт не ок - выходим
			if stopOn == StopOnCheck || result.AccountStatus != 0 {
				return
			}

			result.DefaultSet.AccId = binary.LittleEndian.Uint32(data[1:5]) //id аккаунта
			result.DefaultSet.Premium = uint32(data[5])                     //премиум-флаг

			//запросим спиок серверов
			socket.Send(_Opcode_Out_GetServerList, make([]byte, 4))
		case _Opcode_In_ServerList: //пришел список серверов
			//todo обязательно парсить спиок серверов
			if stopOn == StopOnServerList {
				return
			}

			//подключимся к выбранному серверу
			result.DefaultSet.ServerIndex = cfg.PreferIndex
			serverIndexBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(serverIndexBytes, cfg.PreferIndex)
			socket.Send(_Opcode_Out_ConnectTo, serverIndexBytes)
		case _Opcode_In_ServerInfo: //пришла инфа о сервере
			//Запишем дефолт-сет
			result.DefaultSet.Ip = binary.LittleEndian.Uint32(data[1:5])
			result.DefaultSet.Port = binary.LittleEndian.Uint16(data[5:7])
			result.DefaultSet.Key1 = binary.LittleEndian.Uint32(data[7:11])
			result.DefaultSet.Key2 = binary.LittleEndian.Uint32(data[11:15])
			result.DefaultSet.Key3 = binary.LittleEndian.Uint32(data[15:19])
			result.DefaultSet.Key4 = binary.LittleEndian.Uint32(data[19:23])
			result.DefaultSet.Locale = 0x32D0

			//Запустим игру
			result.BinPid = client.StartRfOnlineBin(result.DefaultSet.Encode())
			return
		}
	}

	return
}
