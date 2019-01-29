// Deals with the checkout process

package tickets

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elafont/CabifyChallenge/items"
)

// Discounts, rule to apply for discounts, every rule will be applied to every line
// of the ticket, if no items are declared, the rule will be applied to the whole ticket
type rule func(int, int, float32) float32
type rules map[*items.Item]rule

type line struct {
	item     *items.Item //
	discount float32     // discount applied to this item
}

// Ticket, description of a ticket
type Ticket struct {
	sync.Mutex
	lines           []*line             // Lines in the ticket
	ticketDiscount  float32             // Discount applied to the whole ticket
	ts              time.Time           // Creation date for the ticket
	ttotal          float32             // Total for the ticket,
	up2date         bool                // false means ticket needs re-Calc
	discRules       rules               // Discount Rules to apply
	itemsDiscounted map[*items.Item]int // number of items being discounted
}

var todayDiscounts = rules{
	items.TSHIRT:  discount3plus,
	items.VOUCHER: discount2x1,
	// nil: discountTicket, // Used to add discounts to the whole ticket
}

// rule to discount half of the products
func discount2x1(line, quantity int, price float32) float32 {
	if (line % 2) == 0 {
		return price // half of the products will be gratis (free as in beer)
	}
	return 0
}

// rule to discount when buying 3 or more
func discount3plus(line, quantity int, price float32) float32 {
	if quantity > 2 {
		return 1 // 1$ discount for each item
	}
	return 0
}

// Added on my own, not asked for the challenge
func discountTicket(line, quantity int, price float32) float32 {
	if price > 500 {
		return 10 // 10$ discount for tickets over 500$
	}
	if quantity > 100 {
		return 5 // 5$ discount when buying over 100 items
	}
	return 0
}

// NewTicket, returns an empty ticket
func NewTicket(discounts rules) *Ticket {
	t := &Ticket{ts: time.Now()}
	t.SetDiscount(discounts)
	return t
}

// SetDiscount, Changes the discount rules applied to the ticket
func (t *Ticket) SetDiscount(discounts rules) {
	t.discRules = discounts
	t.up2date = false
}

// Add items to a ticket
func (t *Ticket) Add(item *items.Item) {
	t.Lock()
	t.lines = append(t.lines, &line{item: item})
	t.up2date = false
	t.Unlock()
}

// Calc, evaluates the discounts and total price of the ticket
func (t *Ticket) Calc() float32 {
	var ttotal float32

	t.Lock()

	t.itemsDiscounted = make(map[*items.Item]int)
	for _, line := range t.lines {
		item := line.item
		price := item.GetPrice()
		line.discount = t.calcDiscountLocked(item, price)
		ttotal += price - line.discount // Sums the total of the ticket discounts included
	}

	t.ticketDiscount = t.calcDiscountLocked(nil, ttotal) // Calc whole ticket discount if applicable

	t.Unlock()
	ttotal -= t.ticketDiscount
	t.ttotal = ttotal
	t.up2date = true
	return ttotal
}

func (t *Ticket) calcDiscountLocked(item *items.Item, price float32) float32 {
	if rule, ok := t.discRules[item]; ok { // if there are a discount rule for this item
		t.itemsDiscounted[item]++
		num := t.itemsDiscounted[item]
		quantity := t.countLocked(item)
		return rule(num, quantity, price) // return the discount to apply
	}
	return 0
}

// counts the number of items available on the ticket, nil item means count all items
func (t *Ticket) countLocked(item *items.Item) (total int) {
	if item == nil {
		return len(t.lines)
	}

	for _, line := range t.lines {
		if line.item == item {
			total++
		}
	}
	return
}

// String, returns the whole ticket in printed form
func (t *Ticket) String() string {
	var ticket strings.Builder
	if t.up2date == false {
		ttotal := t.Calc()
		t.ttotal = ttotal
	}

	header := "CABIFY STORE   date:%19s\n\n"
	subhead := "Article             Price  Disc. Total "
	separator := "------------------- ------ ----- ------"
	itemLine := "%19s % 6.2f % 5.2f % 6.2f\n"
	tDiscount := "		    								 Disc.:  % 5.2f\n"
	footer := "                         Total:  % 5.2f\n"

	ticket.WriteString(fmt.Sprintln())
	ticket.WriteString(fmt.Sprintf(header, t.ts.Format("02, Jan/2006 15:04")))
	ticket.WriteString(fmt.Sprintln(subhead))
	ticket.WriteString(fmt.Sprintln(separator))

	t.Lock()
	for _, line := range t.lines {
		price := line.item.GetPrice()
		ticket.WriteString(fmt.Sprintf(itemLine, line.item.GetName(), price, line.discount, price-line.discount))
	}
	t.Unlock()

	ticket.WriteString(fmt.Sprintln(separator))
	if t.ticketDiscount != 0 {
		ticket.WriteString(fmt.Sprintf(tDiscount, t.ticketDiscount))
	}

	ticket.WriteString(fmt.Sprintf(footer, t.ttotal))
	ticket.WriteString(fmt.Sprintln())
	return ticket.String()
}
