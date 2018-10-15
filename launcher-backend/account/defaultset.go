package account

import (
	"bytes"
	"encoding/binary"
)

type DefaultSet struct {
	Login                  string
	AccId                  uint32
	Premium                uint32
	Ip                     uint32
	Port                   uint16
	Key1, Key2, Key3, Key4 uint32
	ServerIndex            uint16
	Locale                 uint16
}

func (dset DefaultSet) Encode() []byte {
	var buf bytes.Buffer //буфер содержащий дефолт-сет

	tmp32 := make([]byte, 4) //буфер для конвертирования 4байтных чисел
	tmp16 := make([]byte, 2) //буфер для конвертирования 2байтных чисел

	binary.LittleEndian.PutUint32(tmp32, dset.Ip^0xCB9C4B3A)
	buf.Write(tmp32)
	binary.LittleEndian.PutUint16(tmp16, dset.Port^0x4FB6)
	buf.Write(tmp16)

	buf.WriteString(dset.Login)
	buf.Write(make([]byte, 13-len(dset.Login)))

	binary.LittleEndian.PutUint32(tmp32, dset.AccId^0x6E65E0AF)
	buf.Write(tmp32)

	binary.LittleEndian.PutUint32(tmp32, dset.Key1^0xCFCF22E6)
	buf.Write(tmp32)
	binary.LittleEndian.PutUint32(tmp32, dset.Key2^0x5BBCDE6F)
	buf.Write(tmp32)
	binary.LittleEndian.PutUint32(tmp32, dset.Key3^0xACDF5EDA)
	buf.Write(tmp32)
	binary.LittleEndian.PutUint32(tmp32, dset.Key4^0xBCCD1B37)
	buf.Write(tmp32)

	binary.LittleEndian.PutUint16(tmp16, dset.ServerIndex^0x4B3A)
	buf.Write(tmp16)

	binary.LittleEndian.PutUint32(tmp32, dset.AccId^0xC89C183A)
	buf.Write(tmp32)
	binary.LittleEndian.PutUint32(tmp32, dset.Premium^0xC89C183A)
	buf.Write(tmp32)

	buf.Write(make([]byte, 4))
	binary.LittleEndian.PutUint16(tmp16, dset.Locale)
	buf.Write(tmp16)

	return buf.Bytes()
}
