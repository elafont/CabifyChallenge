/* Description of articles available at the store  */

package items

import "errors"

type Item struct {
	name  string
	price float32
}

var ErrNotFound = errors.New("Article Not found")

var (
	VOUCHER = &Item{"Cabify Voucher", 5.00}
	MUG     = &Item{"Cabify Coffee Mug", 7.50}
	TSHIRT  = &Item{"Cabify T-Shirt", 20.00}
)

// As there are no DB, we use a map to store the items available on the store
// There is a small interface to the items storage just in case a real DB should appear later
// Not using a sync.Map as this map is read only.
var items = map[string]*Item{
	"VOUCHER": VOUCHER,
	"MUG":     MUG,
	"TSHIRT":  TSHIRT,
}

// SearchItem Search the storage area for the named item
func SearchItem(name string) (art *Item, err error) {
	art, ok := items[name]
	if !ok {
		return nil, ErrNotFound
	}
	return
}

// GetPrice Return the price of an item,
func (i *Item) GetPrice() float32 {
	return i.price
}

// GetName Return the Name of an item,
func (i *Item) GetName() string {
	return i.name
}
