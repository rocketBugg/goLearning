package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Weapon struct {
	Name       string
	Damage     int
	Durability int
}

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", LowDurability
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
	_, err := fmt.Fprintf(wr, "Weapon|%s|%d|%d", w.Name, w.Damage, w.Durability)
	return err
}

func (w *Weapon) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return ReadingErr
	}
	parts := strings.Split(string(data), "|")

	w.Name = parts[1]
	w.Damage, err = strconv.Atoi(parts[2])
	if err != nil {
		return DmgReadingErr
	}
	w.Durability, err = strconv.Atoi(parts[3])
	if err != nil {
		return DurabilityReadingErr
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
	_, err := fmt.Fprintf(wr, "Armor|%s|%d|%f", a.Name, a.Defense, a.Weight)
	return err
}

func (a *Armor) Deserialize(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return ReadingErr
	}
	parts := strings.Split(string(data), "|")

	a.Name = parts[1]
	a.Defense, err = strconv.Atoi(parts[2])
	if err != nil {
		return DefenceTransformErr
	}
	a.Weight, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return WeightTransformErr
	}
	return nil
}

type Potion struct {
	Name    string
	Effect  string
	Charges int
}

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", LowCharges
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

type Item interface {
	Use() (string, error)
	GetName() string
	GetWeight() float64
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

type Inventory struct {
	Items []Item
}
