package wms_core

import (
	"fmt"
	"github.com/mlplabs/microwms-core/models"
)

func Version() {
	fmt.Println("Version 1.0.0")
}

func GetStorage() *models.Storage {
	storage := new(models.Storage)
	return storage
}
