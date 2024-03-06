package agent

import (
	"bytes"
	"regexp"
	"sort"
	"sync"

	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
)

type CommonHandler struct {
	OIDs  sort.StringSlice
	Items map[string]*Item

	mutex sync.Mutex
}

type Item struct {
	T pdu.VariableType
	V any
}

// Add registry oid func for getting value
func (c *CommonHandler) Add(oid string, t pdu.VariableType, v interface{}) *CommonHandler {
	if c.Items == nil {
		c.Items = make(map[string]*Item)
	}

	c.OIDs = append(c.OIDs, oid)
	c.Items[oid] = &Item{
		T: t,
		V: v,
	}
	return c
}

func (c *CommonHandler) Sort() {
	c.OIDs.Sort()
}

// Get tries to find the provided oid and returns the corresponding value.
func (c *CommonHandler) Get(oid value.OID) (value.OID, pdu.VariableType, interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Items == nil {
		return nil, pdu.VariableTypeNoSuchObject, nil, nil
	}

	item, ok := c.Items[oid.String()]
	if !ok {
		re := regexp.MustCompile(`\.0$`)
		cleaned := re.ReplaceAllString(oid.String(), "")
		item, ok = c.Items[cleaned]
		if ok {
			return oid, item.T, item.V, nil
		}
		return nil, pdu.VariableTypeNoSuchObject, nil, nil
	}

	return oid, item.T, item.V, nil
}

// GetNext tries to find the value that follows the provided oid and returns it.
func (c *CommonHandler) GetNext(from value.OID, includeFrom bool, to value.OID) (value.OID, pdu.VariableType, interface{}, error) {
	if c.Items == nil {
		return nil, pdu.VariableTypeNoSuchObject, nil, nil
	}

	fromOID, toOID := from.String(), to.String()
	for _, oid := range c.OIDs {
		if oidWithin(oid, fromOID, includeFrom, toOID) {
			return c.Get(value.MustParseOID(oid))
		}
	}

	return nil, pdu.VariableTypeNoSuchObject, nil, nil
}

func oidWithin(oid string, from string, includeFrom bool, to string) bool {
	oidBytes, fromBytes, toBytes := []byte(oid), []byte(from), []byte(to)

	fromCompare := bytes.Compare(fromBytes, oidBytes)
	toCompare := bytes.Compare(toBytes, oidBytes)

	return (fromCompare == -1 || (fromCompare == 0 && includeFrom)) && (toCompare == 1)
}

func (c *CommonHandler) UpdateHander(new *CommonHandler) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Items = new.Items
	c.OIDs = new.OIDs
}
