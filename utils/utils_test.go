package utils

import "testing"

func TestGenerateUUID(t *testing.T) {
	uniqueIDs := make(map[string]bool)

	t.Run("test unique id", func(t *testing.T) {

		for i := 0; i < 1000000; i++ {
			gotUuid := GenerateUUID()
			if uniqueIDs[gotUuid] {
				t.Fatalf("This not uniq id: %v", gotUuid)
			}

			uniqueIDs[gotUuid] = true
		}
	})
}
