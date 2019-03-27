package utils

import (
	"log"
	"testing"
)

func TestGenerateUID(t *testing.T) {
	uniqueIDs := make(map[UID]bool)

	t.Run("test unique id", func(t *testing.T) {

		for i := 0; i < 1000000; i++ {
			gotUuid, err := GenerateUID()
			if err != nil {
				t.Fatalf("GenerateUID() error: %v", err)
			}
			if uniqueIDs[gotUuid] {
				t.Fatalf("This not uniq id: %v", gotUuid)
			}

			uniqueIDs[gotUuid] = true

			log.Println(gotUuid)
		}
	})
}
