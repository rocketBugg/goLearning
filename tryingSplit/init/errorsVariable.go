package example

import (
	"errors"
)

var (
	LowDurability        = errors.New("у предмета закончилась прочность")
	LowCharges           = errors.New("закончились заряды зелья")
	MissingItem          = errors.New("предмет отсутствует")
	NilItem              = errors.New("предмет не существует")
	OpenFileErr          = errors.New("ошибка открытия файла")
	SaveFileErr          = errors.New("ошибка сохранения файла")
	LoadFileErr          = errors.New("ошибка загрузки файла")
	DmgReadingErr        = errors.New("ошибка чтения свойства Damage")
	DurabilityReadingErr = errors.New("ошибка чтения свойства Durability")
	ReadingErr           = errors.New("ошибка чтения")
	DefenceTransformErr  = errors.New("ошибка преобразования Defense")
	WeightTransformErr   = errors.New("ошибка преобразования Weight")
	SaveInventoryErr     = errors.New("не удалось сохранить инвентарь")
	LoadInventoryErr     = errors.New("ошибка загрузки файла")
)
