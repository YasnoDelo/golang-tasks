package reflectspell

import (
	"log"
	"reflect"
)

type Spell interface {
	// название заклинания
	Name() string
	// характеристика, на которую воздействует
	Char() string
	// количественное значение
	Value() int
}

// CastReceiver — если объект удовлетворяет этом интерфейсу, то заклинание применяется через него
type CastReceiver interface {
	ReceiveSpell(s Spell)
}

func CastToAll(spell Spell, objects []interface{}) {
	for _, obj := range objects {
		CastTo(spell, obj)
	}
}

func CastTo(spell Spell, object interface{}) {
	// Попытка передать заклинание объекту, реализующему интерфейс CastReceiver
	if receiver, ok := object.(CastReceiver); ok {
		receiver.ReceiveSpell(spell)
		return
	}

	// Получение значения объекта с помощью рефлексии
	v := reflect.ValueOf(object)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Попытка найти поле по имени и применить заклинание
	field := v.FieldByName(spell.Char())
	if !field.IsValid() || !field.CanSet() {
		logInvalidField(spell, object)
		return
	}

	applySpellToField(spell, field)
}

func logInvalidField(spell Spell, object interface{}) {
	log.Printf("Field %s not found or cannot be set in object %v", spell.Char(), object)
}

func applySpellToField(spell Spell, field reflect.Value) {
	if field.Kind() == reflect.Int {
		field.SetInt(field.Int() + int64(spell.Value()))
	} else {
		logUnsupportedFieldType(spell, field)
	}
}

func logUnsupportedFieldType(spell Spell, field reflect.Value) {
	log.Printf("Cannot apply spell to field %s: unsupported type %s", spell.Char(), field.Kind())
}

type spell struct {
	name string
	char string
	val  int
}

func NewSpell(name string, char string, val int) Spell {
	return &spell{name: name, char: char, val: val}
}

func (s spell) Name() string {
	return s.name
}

func (s spell) Char() string {
	return s.char
}

func (s spell) Value() int {
	return s.val
}

type Player struct {
	// nolint: unused
	name   string
	health int
}

func (p *Player) ReceiveSpell(s Spell) {
	if s.Char() == "Health" {
		p.health += s.Value()
	}
}

type Zombie struct {
	Health int
}

type Daemon struct {
	Health int
}

type Orc struct {
	Health int
}

type Wall struct {
	Durability int
}

// func main() {
// 	player := &Player{
// 		name:   "Player_1",
// 		health: 100,
// 	}

// 	enemies := []interface{}{
// 		&Zombie{Health: 1000},
// 		&Zombie{Health: 1000},
// 		&Orc{Health: 500},
// 		&Orc{Health: 500},
// 		&Orc{Health: 500},
// 		&Daemon{Health: 1000},
// 		&Daemon{Health: 1000},
// 		&Wall{Durability: 100},
// 	}

// 	CastToAll(NewSpell("fire", "Health", -50), append(enemies, player))
// 	CastToAll(NewSpell("heal", "Health", 190), append(enemies, player))

// 	fmt.Println(player)
// 	for _, enemy := range enemies {
// 		fmt.Println(enemy)
// 	}
// }
