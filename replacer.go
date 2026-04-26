// 1.0.0
package replacer

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 1.0.0
type Replacer struct {
	MChar rune
	r     map[string]interface{}
}

// Динамическая дата
func New() *Replacer {
	r := &Replacer{
		MChar: '%',
		r:     make(map[string]interface{}),
	}
	r.init()
	return r
}

// Фиксированная заданная дата/время на момент создания
func NewFix(tm time.Time) *Replacer {
	r := &Replacer{MChar: '%', r: make(map[string]interface{})}
	r.initFix(tm)
	return r
}

// Инициализация базовых тегов
func (r *Replacer) init() {
	if r.r == nil {
		r.r = make(map[string]interface{})
	}
	r.r["date"] = func() string { return time.Now().Format("2006-01-02") }
	r.r["date.Y"] = func() string { return time.Now().Format("2006") }
	r.r["date.M"] = func() string { return time.Now().Format("01") }
	r.r["date.D"] = func() string { return time.Now().Format("02") }
	r.r["time"] = func() string { return time.Now().Format("15:04:05") }
	r.r["time.H"] = func() string { return time.Now().Format("15") }
	r.r["time.M"] = func() string { return time.Now().Format("04") }
	r.r["time.S"] = func() string { return time.Now().Format("05") }
	r.r["time.Z"] = func() string { return time.Now().Format(".000")[1:] }     //Hack for time.Format()
	r.r["time.ns"] = func() string { return time.Now().Format(".000000")[1:] } //Hack for time.Format()
}

func (r *Replacer) initFix(tm time.Time) {
	r.r["date"] = tm.Format("2006-01-02")
	r.r["date.Y"] = tm.Format("2006")
	r.r["date.M"] = tm.Format("01")
	r.r["date.D"] = tm.Format("02")
	r.r["time"] = tm.Format("15:04:05")
	r.r["time.H"] = tm.Format("15")
	r.r["time.M"] = tm.Format("04")
	r.r["time.S"] = tm.Format("05")
	r.r["time.Z"] = tm.Format(".000")[1:]     //Hack for time.Format()
	r.r["time.ns"] = tm.Format(".000000")[1:] //Hack for time.Format(
}

// Произвести замену, удалить неопознанные теги, вернуть ошибки
func (r *Replacer) ReplaceCE(s string) (string, []error) {
	ret, e := r.replace(s, true)
	if len(e) > 0 {
		return ret, e
	}
	return ret, nil
}

// Произвести замену, не менять неопознанные теги, вернуть ошибки
func (r *Replacer) ReplaceE(s string) (string, []error) {
	ret, e := r.replace(s, false)
	if len(e) > 0 {
		return ret, e
	}
	return ret, nil
}

// Произвести замену, удалить неопознанные теги
func (r *Replacer) ReplaceC(s string) string {
	ret, _ := r.replace(s, true)
	return ret
}

// Произвести замену, не менять неопознаные теги
func (r *Replacer) Replace(s string) string {
	ret, _ := r.replace(s, false)
	return ret
}

func (r *Replacer) replace(s string, clr bool) (ret string, e []error) {
	e = make([]error, 0)
	b := bytes.NewBuffer([]byte(""))
	max := len(s)
	pi := 0
	for i := 0; i < max; {
		for i < max && s[i] != byte(r.MChar) {
			i++
		}
		b.WriteString(s[pi:i])
		pi = i
		i++
		if i > max {
			return b.String(), e
		}
		for i < max && s[i] != byte(r.MChar) {
			i++
		}
		key := strings.Split(s[pi+1:i], " ")

		//У ключа может быть параметр
		i++
		pi = i
		if v, ok := r.r[key[0]]; ok {
			switch vv := v.(type) {
			case string:
				b.WriteString(vv)
			case func() string: //Функция без параметра
				b.WriteString(vv())
			case func(string) string: //Функция с параметром
				if len(key) > 1 {
					b.WriteString(vv(key[1]))
				}
			case func(...string) string: //Функция с несколькими параметрами
				if len(key) > 1 {
					var args []reflect.Value
					for _, x := range key[1:] {
						args = append(args, reflect.ValueOf(x))
					}
					fn := reflect.ValueOf(vv)
					res := fn.Call(args)
					s := res[0].Interface().(string)
					b.WriteString(s)
				}
			default:
				//Unknown type
			}
		} else { //Если ключа не нашли -
			e = append(e, errors.New("Key not defined: "+string(r.MChar)+strings.Join(key, " ")+string(r.MChar)))
			if clr { //Исключаем ключ
			} else { //Не меняем ключ
				b.WriteString(string(r.MChar) + strings.Join(key, " ") + string(r.MChar))
			}
		}
	}
	return b.String(), e
}

// Добавить токен замены
//
//	1.0.0
func (r *Replacer) Add(key string, fn interface{}) error {
	switch v := fn.(type) {
	case string, func() string, func(string) string, func(...string) string:
		r.r[key] = v
	default:
		return errors.New("Unsupported type: " + fmt.Sprintf("%T", v))
	}
	return nil
}

// Удалить токен замены
func (r *Replacer) Del(key string) error {
	delete(r.r, key)
	return nil
}
