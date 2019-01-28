package tickets

import (
	"sync"
	"testing"

	"github.com/elafont/CabifyChallenge/items"
)

func TestTicket(t *testing.T) {
	var wg sync.WaitGroup
	var buyList = []struct {
		list     string
		total    float32
		articles []*items.Item
	}{
		{"Normal", 32.50, []*items.Item{items.VOUCHER, items.TSHIRT, items.MUG}},
		{"2x1", 25.00, []*items.Item{items.VOUCHER, items.TSHIRT, items.VOUCHER}},
		{"3plus", 81.00, []*items.Item{items.TSHIRT, items.TSHIRT, items.TSHIRT, items.VOUCHER, items.TSHIRT}},
		{"AllDisc", 74.50, []*items.Item{items.VOUCHER, items.TSHIRT, items.VOUCHER, items.VOUCHER, items.MUG, items.TSHIRT, items.TSHIRT}},
	}

	for _, bl := range buyList {
		wg.Add(1)
		go t.Run(bl.list, func(t *testing.T) {
			ticket := NewTicket(todayDiscounts)
			for _, item := range bl.articles {
				ticket.Add(item)
			}
			ticket.Calc()
			t.Log(*ticket.ttotal)
			if *ticket.ttotal != bl.total {
				t.Error("Bad Total")
			}
			// fmt.Println(ticket.String())
			wg.Done()
		})
		wg.Wait()
	}
}

func TestDiscount2x1(t *testing.T) {
	ticket := NewTicket(todayDiscounts)
	voucher := items.VOUCHER
	price := voucher.GetPrice()
	ticket.Add(voucher)
	ticket.Add(voucher)
	if ticket.Calc() != price {
		t.Error("Discount not applied")
	}

	ticket.Add(voucher)
	if ticket.Calc() != (2 * price) {
		t.Error("Total badly calculated")
	}

	ticket.Add(voucher)
	if ticket.Calc() != (2 * price) {
		t.Error("Total badly calculated")
	}

}

func TestDiscount3plus(t *testing.T) {
	ticket := NewTicket(todayDiscounts)
	shirt := items.TSHIRT
	price := shirt.GetPrice()
	ticket.Add(shirt)
	ticket.Add(shirt)

	if ticket.Calc() != (2 * price) {
		t.Error("Total badly calculated")
	}

	ticket.Add(shirt)
	if ticket.Calc() != 3*(price-1) {
		t.Error("Discount not applied")
	}

	ticket.Add(shirt)
	if ticket.Calc() != 4*(price-1) {
		t.Error("Discount not applied")
	}
}
