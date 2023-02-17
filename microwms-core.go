package microwms_core

import (
	"fmt"
	"github.com/mlplabs/microwms-core/whs"
)

func Version() {
	fmt.Println("Version 1.0.0")
}

func GetStorage() *whs.Storage {
	storage := new(whs.Storage)
	return storage
}
