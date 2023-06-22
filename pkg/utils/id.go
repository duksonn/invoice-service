package utils

import (
	"github.com/google/uuid"
)

var GenerateUuid func() string = uuid.NewString

func PatchGenerateUuid(generateUuid func() string) {
	GenerateUuid = generateUuid
}

func RestoreGenerateUuid() {
	GenerateUuid = uuid.NewString
}
