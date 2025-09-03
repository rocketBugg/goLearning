package example

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Inventory struct {
	Items []Item
}

func SafeUse(item Item) (res string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %v", rec)
		}
	}()

	if w, ok := item.(*Weapon); ok && w.Name == "Ящик Пандоры" {
		panic("Проклятье Пандоры!")
	}

	return item.Use()
}

func DescribeItem(i Item) (string, error) {
	if i == nil {
		return "", MissingItem
	}

	return fmt.Sprintf("%s (вес: %.1f)", i.GetName(), i.GetWeight()), nil
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	var result []T

	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, len(items))

	for i, item := range items {
		result[i] = transform(item)
	}

	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

func (inv *Inventory) AddItem(item Item) error {
	if item == nil {
		return NilItem
	}
	inv.Items = append(inv.Items, item)
	return nil
}

func (inv *Inventory) GetWeapons() []*Weapon {
	weapons := Filter(inv.Items, func(item Item) bool {
		_, ok := item.(*Weapon)
		return ok
	})

	return Map(weapons, func(item Item) *Weapon {
		return item.(*Weapon)
	})
}

func (inv *Inventory) GetBrokenItems() []Item {
	return Filter(inv.Items, func(item Item) bool {
		switch v := item.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		default:
			return false
		}
	})
}

func (inv *Inventory) GetItemNames() []string {
	return Map(inv.Items, func(item Item) string {
		return item.GetName()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(item Item) bool {
		return item.GetName() == name
	})
}

func (inv *Inventory) Save(w io.Writer) error {
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			storable.Serialize(w)

			_, SaveInventoryErr = fmt.Fprintln(w)
		}
	}
	return SaveInventoryErr
}

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Weapon") {
			var w Weapon

			reader := strings.NewReader(line)

			w.Deserialize(reader)

			inv.AddItem(&w)
		} else if strings.HasPrefix(line, "Armor") {
			var a Armor

			reader := strings.NewReader(line)

			a.Deserialize(reader)

			inv.AddItem(&a)
		}
	}
	return LoadInventoryErr
}
