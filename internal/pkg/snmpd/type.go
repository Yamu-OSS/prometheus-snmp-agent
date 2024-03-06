package snmpd

import (
	"strconv"

	"github.com/posteo/go-agentx/pdu"

	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
)

const (
	SubTypeTable = "TABLE"
	SubTypeList  = "LIST"
)

var valueType = map[string]func(string) (pdu.VariableType, any){
	"CEILINT": ToCEILINT,
	"STRING":  ToSTRING,
}

func ToCEILINT(s string) (pdu.VariableType, any) {
	floatVal, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Error(err, "Error parsing float")
	}

	return pdu.VariableTypeInteger, int32(floatVal)
}

func ToSTRING(s string) (pdu.VariableType, any) {
	return pdu.VariableTypeOctetString, s
}
