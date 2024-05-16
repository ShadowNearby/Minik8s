package utils

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

// GenerateUUID generate a random uuid
func GenerateUUID() string {
	return fmt.Sprintf("%d", rand.Int())
}

// GenerateContainerIDByName receive container name + pod name as param
func GenerateContainerIDByName(containerName string, podUUID string) string {
	id := uuid.NewMD5(uuid.NameSpaceDNS, []byte(containerName+podUUID))
	return id.String()[:12]
}
