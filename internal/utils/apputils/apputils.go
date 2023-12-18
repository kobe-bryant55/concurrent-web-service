package apputils

import (
	"bytes"
	"encoding/json"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
	"strconv"
)

func InterfaceToStruct(in, out interface{}) error {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONEncode, err)
	}

	err = json.NewDecoder(buf).Decode(out)
	if err != nil {
		return errorutils.New(errorutils.ErrJSONDecode, err)
	}

	return nil
}

func StringToUINT64(s string) (*uint64, error) {
	pu, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}

	return &pu, nil
}
