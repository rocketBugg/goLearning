package main

import (
	"fmt"
	example "main/init"
	"math"
	"os"
)

func main() {
	inv := example.Inventory{}

	sword := &example.Weapon{Name: "Меч", Damage: 10, Durability: 5}
	healthPotion := &example.Potion{Name: "Лечебное", Effect: "+50 HP", Charges: 3}
	pandoraBox := &example.Weapon{Name: "Ящик Пандоры", Damage: math.MaxInt, Durability: math.MaxInt}

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

	// SafeUse с обработкой паники для Ящика Пандоры
	fmt.Println("\nSafeUse Ящик Пандоры (ожидается panic):")
	msg, err := example.SafeUse(pandoraBox)
	if err != nil {
		fmt.Println("Поймана паника:", err)
	} else {
		fmt.Println(msg)
	}

	fmt.Println(example.DescribeItem(sword))
	fmt.Println(example.DescribeItem(nil))

	fmt.Println("\nСохраняем в файл")

	file, err := os.OpenFile("homework_solved.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(example.OpenFileErr)
		return
	}
	inv.Save(file)
	if err := inv.Save(file); err != nil {
		fmt.Println(example.SaveFileErr)
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
	inv = example.Inventory{}

	file, err = os.Open("homework_solved.txt")
	if err != nil {
		fmt.Println(example.OpenFileErr)
	}

	inv.Load(file)
	if err := inv.Load(file); err != nil {
		fmt.Println(example.LoadFileErr)
		return
	}

	names := inv.GetItemNames()

	fmt.Println("\nИмена предметов:", names)

	for _, item := range inv.Items {
		fmt.Println(example.DescribeItem(item))
	}
}
