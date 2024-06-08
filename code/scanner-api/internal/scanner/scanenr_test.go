package scanner

import "testing"

func TestGetRAMSize(t *testing.T) {
	ramSize, err := getRAMSize()
	if err != nil {
		t.Errorf("Error getting RAM Size: %v", err)
	}

	if ramSize == 0 {
		t.Errorf("RAM Size is 0")
	}

	t.Logf("RAM Size: %d", ramSize)

}
