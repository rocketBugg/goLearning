package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// TODO: Реализуйте структуры:
// TODO: - Weapon: Name (string), Damage (int), Durability (int)
// TODO: - Armor: Name (string), Defense (int), Weight (float64)
// TODO: - Potion: Name (string), Effect (string), Charges (int)
// TODO:
// TODO: Можете добавить свои структуры :)
// TODO:
// TODO: Для каждой структуры реализуйте:
// TODO: - Метод Use() string (описание использования, например "Используется <имя>" и т.д.)
// TODO: - Методы интерфейса Item

type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

func (w *Weapon) Use() string {
	if w.Durability <= 0 {
		fmt.Printf("%s сломан и не может быть использован!\n", w.Name)
		return w.Name
	}
	w.Durability--
	fmt.Printf("Используется %s, прочность: %d \n", w.Name, w.Durability)
	return w.Name
}

func (w *Weapon) GetName() string {
	return w.Name
}

func (w *Weapon) GetWeight() float64 {
	return float64(w.Damage * 3)
}

func (w *Weapon) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Weapon|%s|%d|%d\n", w.Name, w.Damage, w.Durability)
}

func (w *Weapon) Deserialize(r io.Reader) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		fileContent := scanner.Text()
		parts := strings.Split(fileContent, "|")
		w.Name = parts[1]
		w.Damage, _ = strconv.Atoi(parts[2])
		w.Durability, _ = strconv.Atoi(parts[3])
	}
}

type Armor struct {
	Name    string
	Defense int
	Weight  float64
}

func (a *Armor) Use() string {
	fmt.Sprintf("Броня - %s экипирована\n", a.Name)
	return a.Name
}

func (a *Armor) GetName() string {
	return a.Name
}

func (a *Armor) GetWeight() float64 {
	return a.Weight
}

func (a *Armor) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Armor|%s|%d|%.2f\n", a.Name, a.Defense, a.Weight)
}

func (a *Armor) Deserialize(r io.Reader) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		fileContent := scanner.Text()
		parts := strings.Split(fileContent, "|")
		a.Name = parts[1]
		a.Defense, _ = strconv.Atoi(parts[2])
		a.Weight, _ = strconv.ParseFloat(parts[3], 64)
	}
}

const PotionWeight float64 = 15

type Potion struct {
	Name    string
	Effect  string
	Charges int
}

func (p *Potion) Use() string {
	if p.Charges <= 0 {
		fmt.Printf("Заряды у %s закночились и зелье не может быть использовано!\n", p.Name)
		return p.Name
	}
	p.Charges--
	fmt.Printf("Используется %s\n", p.Name)
	return p.Name
}

func (p *Potion) GetName() string {
	return p.Name
}

func (p *Potion) GetWeight() float64 {
	return PotionWeight
}

type Item interface {
	GetName() string
	GetWeight() float64
	Use() string
}

// TODO: Реализуйте функцию
func DescribeItem(i Item) string {
	// Функция должна возвращать:
	// - "Предмет отсутствует" если i == nil
	// - "<название> (вес: <вес>)" в остальных случаях
	if i == nil {
		fmt.Println("Предмет отсутствует")
		return "Предмет отсутствует"
	}
	return fmt.Sprintf("%s (вес: %.2f)", i.GetName(), i.GetWeight())
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	// TODO: Верните новый слайс только с элементами, для которых predicate вернул true
	res := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			res = append(res, item)
		}
	}
	return res
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	// TODO: Примените transform к каждому элементу и верните слайс с результатами
	result := make([]R, 0)
	for _, i := range items {
		item := transform(i)
		result = append(result, item)
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	// TODO: Найдите первый элемент, удовлетворяющий condition и верните элемент и true или zero value и false
	for _, i := range items {
		if condition(i) {
			return i, true
		}
	}
	var zero T
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) GetWeapons() []*Weapon {
	// TODO: Используйте Filter для отбора Weapon, затем Map для преобразования Item -> *Weapon
	weapons := Filter(inv.Items, func(item Item) bool {
		_, ok := item.(*Weapon)
		return ok
	})
	return Map(weapons, func(item Item) *Weapon {
		return item.(*Weapon)
	})
}

func (inv *Inventory) GetBrokenItems() []Item {
	// TODO: Используйте Filter для отбора:
	// TODO: - Weapon: Durability <= 0
	// TODO: - Potion: Charges <= 0
	// TODO:
	// TODO: Подсказка: поможет приведение типов - item.(type)
	return Filter(inv.Items, func(item Item) bool {
		switch v := item.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		}
		return false
	})
}

func (inv *Inventory) GetItemNames() []string {
	// TODO: Используйте Map для преобразования []Item -> []string
	return Map(inv.Items, func(item Item) string {
		return item.GetName()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	// TODO: Используйте Find для поиска по имени
	return Find(inv.Items, func(item Item) bool {
		return item.GetName() == name
	})
}

// TODO: Бонус: реализуйте интефейс Storable для Weapon и Armor:
// TODO: - Weapon: формат "Weapon|Name|Damage|Durability"
// TODO: - Armor: формат "Armor|Name|Defense|Weight"

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (inv *Inventory) Save(w io.Writer) {
	// TODO: Бонус: сделайте сохранение/загрузку инвентаря в/из файла
	for _, item := range inv.Items {
		switch t := item.(type) {
		case *Weapon:
			t.Serialize(w)
		case *Armor:
			t.Serialize(w)
		}
	}
}

func (inv *Inventory) Load(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		switch parts[0] {
		case "Weapon":
			damage, _ := strconv.Atoi(parts[2])
			durability, _ := strconv.Atoi(parts[3])
			inv.AddItem(&Weapon{
				Name:       parts[1],
				Damage:     damage,
				Durability: durability,
			})
		case "Armor":
			defense, _ := strconv.Atoi(parts[2])
			weight, _ := strconv.ParseFloat(parts[3], 64)
			inv.AddItem(&Armor{
				Name:    parts[1],
				Defense: defense,
				Weight:  weight,
			})
		}
	}
}

func main() {

	// TODO: Создайте инвентарь и добавьте:
	inventory := &Inventory{}
	// TODO: - Оружие: "Меч" (урон 10, прочность 5)
	inventory.AddItem(&Weapon{"Меч", 10, 5})
	// TODO: - Броню: "Щит" (защита 5, вес 4.5)
	inventory.AddItem(&Armor{"Щит", 5, 4.5})
	// TODO: - Зелье: "Лечебное" (+50 HP, 3 заряда)
	inventory.AddItem(&Potion{"Лечебное", "+50 HP", 3})
	// TODO: - Оружие: "Сломанный лук" (урон 5, прочность 0)
	inventory.AddItem(&Weapon{"Сломанный лук", 5, 0})
	// TODO: Реализуйте логику/вызовы:
	// TODO: 1. Use предмета с выводом в консоль
	inventory.Items[0].Use()
	// TODO: 2. DescribeItem с предметом и с nil, так же с выводом в консоль
	DescribeItem(inventory.Items[2])
	fmt.Println(DescribeItem(inventory.Items[2]))
	DescribeItem(nil)
	//TODO: 3. Вывести в консоль результат вызова GetWeapons (должны вернуться только меч и лук)
	for _, item := range inventory.GetWeapons() {
		fmt.Println(DescribeItem(item))
	}
	// TODO: 4. Вывести в консоль результат вызова GetBrokenItems (должен вернуть сломанный лук)
	for _, item := range inventory.GetBrokenItems() {
		fmt.Println(DescribeItem(item))
	}
	// TODO: 5. Вывести в консоль результат вызова GetItemNames (все названия)
	inventory.GetItemNames()
	fmt.Println(inventory.GetItemNames())
	// TODO: 6. Вывести в консоль результат вызова FindItemByName (поиск "Щит")
	item, found := inventory.FindItemByName("Щит")
	if found {
		fmt.Println("Найдено:", DescribeItem(item))
	} else {
		fmt.Println("Щит не найден")
	}
	// TODO: Бонус: сделайте сохранение инвентаря в файл и загрузку инвентаря из файла
	// Load
	newInventory := &Inventory{}
	f, _ := os.Open(".inventory.text")
	newInventory.Load(f)
	defer f.Close()

	// Save

	f2, _ := os.Create(".inventory.text")
	newInventory.Save(f2)
	defer f2.Close()
}
