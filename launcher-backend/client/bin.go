package client

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Запускает бин и возврщает его pid
func StartRfOnlineBin(defaultSet []byte) int {
	//Запишем дефолтсет
	err := ioutil.WriteFile("system/defaultset.tmp", defaultSet, os.ModeExclusive)
	if err != nil {
		log.Printf("Error with writing defaultset: %v", err)
		return -1
	}

	//Запустим бин
	bin := exec.Command("rf_online.bin")
	err = bin.Start()
	if err != nil {
		log.Printf("Error with starting bin: %v", err)
		return -2
	}

	//Возвращаем pid
	return bin.Process.Pid
}
