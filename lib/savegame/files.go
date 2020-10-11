package savegame

import (
	"encoding/binary"
	"io/ioutil"
	"os"

	"gopkg.in/restruct.v1"
)

func LoadFile(path string, obj interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = restruct.Unpack(buf, binary.LittleEndian, &obj); err != nil {
		return err
	}
	return nil
}

func SaveFile(path string, obj interface{}) error {
	buf, err := restruct.Pack(binary.LittleEndian, &obj)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, buf, os.ModePerm); err != nil {
		return err
	}
	return nil
}
