package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// Кастомные ошибки
var (
	// TODO: Добавьте необходимые кастомные ошибки
	lowDurability        = errors.New("у предмета закончилась прочность")
	lowCharges           = errors.New("закончились заряды зелья")
	missingItem          = errors.New("предмет отсутствует")
	nilItem              = errors.New("предмет не существует")
	openFile             = errors.New("ошибка открытия файла")
	saveFileErr          = errors.New("ошибка сохранения файла")
	loadFileErr          = errors.New("ошибка загрузки файла")
	dmgReadingErr        = errors.New("ошибка чтения свойства Damage")
	durabilityReadingErr = errors.New("ошибка чтения свойства Durability")
	readingErr           = errors.New("ошибка чтения")
	defenceTransformErr  = errors.New("ошибка преобразования Defense")
	weightTranformErr    = errors.New("ошибка преобразования Weight")
	saveInventoryErr     = errors.New("не удалось сохранить инвентарь")
	loadInventoryErr     = errors.New("ошибка загрузки файла")
)

type Item interface {
	Use() (string, error)
	GetName() string
	GetWeight() float64
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

func (w *Weapon) Use() (string, error) {
	// TODO: Реализуйте возврат ошибки
	if w.Durability <= 0 {
		return "", lowDurability
	}

	w.Durability--

	return fmt.Sprintf("Атаковали %s (%d урона)", w.Name, w.Damage), nil
}

func (w *Weapon) GetName() string {
	return w.Name
}

func (w *Weapon) GetWeight() float64 {
	return 2.5
}

func (w *Weapon) Serialize(wr io.Writer) error {
	// TODO: Реализуйте возврат ошибки
	_, err := fmt.Fprintf(wr, "Weapon|%s|%d|%d", w.Name, w.Damage, w.Durability)
	return err
}

func (w *Weapon) Deserialize(r io.Reader) error {
	// TODO: Реализуйте возврат ошибок
	data, err := io.ReadAll(r)
	if err != nil {
		return readingErr
	}
	parts := strings.Split(string(data), "|")

	w.Name = parts[1]
	w.Damage, err = strconv.Atoi(parts[2])
	if err != nil {
		return dmgReadingErr
	}
	w.Durability, err = strconv.Atoi(parts[3])
	if err != nil {
		return durabilityReadingErr
	}
	return nil
}

type Armor struct {
	Name    string
	Defense int
	Weight  float64
}

func (a *Armor) Use() (string, error) {
	return fmt.Sprintf("Надели %s (+%d защиты)", a.Name, a.Defense), nil
}

func (a *Armor) GetName() string {
	return a.Name
}

func (a *Armor) GetWeight() float64 {
	return a.Weight
}

func (a *Armor) Serialize(wr io.Writer) error {
	// TODO: Реализуйте возврат ошибки
	_, err := fmt.Fprintf(wr, "Armor|%s|%d|%f", a.Name, a.Defense, a.Weight)
	return err
}

func (a *Armor) Deserialize(r io.Reader) error {
	// TODO: Реализуйте возврат ошибок
	data, err := io.ReadAll(r)
	if err != nil {
		return readingErr
	}
	parts := strings.Split(string(data), "|")

	a.Name = parts[1]
	a.Defense, err = strconv.Atoi(parts[2])
	if err != nil {
		return defenceTransformErr
	}
	a.Weight, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return weightTranformErr
	}
	return nil
}

type Potion struct {
	Name    string
	Effect  string
	Charges int
}

func (p *Potion) Use() (string, error) {
	// TODO: Реализуйте возврат ошибки
	if p.Charges <= 0 {
		return "", lowCharges
	}

	p.Charges--

	return fmt.Sprintf("Использовали %s (%s)", p.Name, p.Effect), nil
}

func (p *Potion) GetName() string {
	return p.Name
}

func (p *Potion) GetWeight() float64 {
	return 0.5
}

func DescribeItem(i Item) (string, error) {
	// TODO: Реализуйте возврат ошибки
	if i == nil {
		return "", missingItem
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

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) error {
	// TODO: Проверка на nil
	if item == nil {
		return nilItem
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
	// TODO: Реализуйте возврат ошибки
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			storable.Serialize(w)

			_, saveInventoryErr = fmt.Fprintln(w)
		}
	}
	return saveInventoryErr
}

func (inv *Inventory) Load(r io.Reader) error {
	// TODO: Реализуйте возврат ошибки
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
	return loadInventoryErr
}

func SafeUse(item Item) (res string, err error) {
	// TODO: Используйте defer с recover для перехвата паники
	// TODO: Для оружия с именем "Ящик Пандоры" вызовите панику
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

func main() {

	// TODO: Реализуйте логику/вызовы:
	// TODO: 1. Обработку ошибок везде
	// TODO: 2. Use предмета до потери прочности и обработку ошибки при потере прочности
	// TODO: 3. DescribeItem с предметом и с nil
	// TODO: 4. Обработку ошибок сохранения/загрузки в файл
	// TODO: 5. Обработку паники для "Ящика Пандоры"
	inv := Inventory{}

	sword := &Weapon{Name: "Меч", Damage: 10, Durability: 5}
	healthPotion := &Potion{Name: "Лечебное", Effect: "+50 HP", Charges: 3}
	pandoraBox := &Weapon{Name: "Ящик Пандоры", Damage: math.MaxInt, Durability: math.MaxInt}

	inv.AddItem(sword)
	inv.AddItem(healthPotion)
	inv.AddItem(pandoraBox)
	inv.AddItem(nil)

	for {
		msg, err := sword.Use()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(msg)
	}

	// SafeUse с обработкой паники для Ящик Пандоры
	fmt.Println("\nSafeUse Ящик Пандоры (ожидается panic):")
	msg, err := SafeUse(pandoraBox)
	if err != nil {
		fmt.Println("Поймана паника:", err)
	} else {
		fmt.Println(msg)
	}

	fmt.Println(DescribeItem(sword))
	fmt.Println(DescribeItem(nil))

	fmt.Println("\nСохраняем в файл")

	file, err := os.OpenFile("homework_solved.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(openFile)
		return
	}
	inv.Save(file)
	if err := inv.Save(file); err != nil {
		fmt.Println(saveFileErr)
		return
	}

	file.Close()

	fmt.Println("Ломаем файл")
	file, err = os.OpenFile("homework_solved.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		fmt.Fprintln(file, "Weapon||")
		file.Close()
	}
	fmt.Fprintf(file, "Weapon||")

	fmt.Println("Загружаем из файла")
	inv = Inventory{}

	file, err = os.Open("homework_solved.txt")
	if err != nil {
		fmt.Println(openFile)
	}

	inv.Load(file)
	if err := inv.Load(file); err != nil {
		fmt.Println(loadFileErr)
		return
	}

	names := inv.GetItemNames()

	fmt.Println("\nИмена предметов:", names)

	for _, item := range inv.Items {
		fmt.Println(DescribeItem(item))
	}
}
