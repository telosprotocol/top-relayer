package ethbeacon_rpc

import "testing"

func TestGetBeforeSlotInSamePeriod(t *testing.T) {
	slot, err := getBeforeSlotInSamePeriod(2302239)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(slot)
}
