package service

import "testing"

func TestIsAlphanumeric(t *testing.T) {
	got := isAlphaNumeric(22)
	want := false
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestAssignPoints(t *testing.T) {
	items := []Item{
		{
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		},
	}
	got, err := assignPoints(ReceiptToProcess{Retailer: "M&M Corner Market", PurchaseDate: "2022-03-20", PurchaseTime: "14:33", Items: items, Total: "9.00"})
	if err != nil {
		t.Errorf("errored out")
	}

	want := 109
	if got != int32(want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
