package enums

import (
	"testing"
)

func TestIntToBTSeverity(t *testing.T) {
	v, err := BTSeverityMap.ByIndex(1)

	if err != nil {
		t.Fatalf("error getting severity: %v", err)
	}

	if v != BTSeverityCritical {
		t.Fatalf("invalid severity. expected %v got %v", BTSeverityCritical, v)
	}

	v, err = BTSeverityMap.ByIndex("2")
	if err != nil {
		t.Fatalf("error getting severity: %v", err)
	}

	if v != BTSeverityHigh {
		t.Fatalf("invalid severity. expected %v got %v", BTSeverityHigh, v)
	}

	iv, err := BTSeverityMap.ByIndex(10000)
	if err == nil {
		t.Fatalf("expected error getting invalid severity. got %v", iv)
	}

}
