package datastore

import (
	"fmt"

	"github.com/thisdougb/go-fsvault/fsvault"
)

var (
	orderKeyFmt = "/order/%d"
)

func ListOrders() []string {

	key := "/order/"
	return fsvault.List(getVaultRoot(), key)
}

func GetOrder(id int64) ([]byte, error) {
	key := fmt.Sprintf(orderKeyFmt, id)
	return fsvault.Get(getVaultRoot(), key)
}

func StoreOrder(id int64, data []byte) error {
	key := fmt.Sprintf(orderKeyFmt, id)
	return fsvault.Put(getVaultRoot(), key, data)
}
