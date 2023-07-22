package luaUtils

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"log"
	"reflect"
)

type Clock struct {
	TotalTimeInPhase     int     `lua:"totaltimeinphase"`
	Cycles               int     `lua:"cycles"`
	Phase                string  `lua:"phase"`
	RemainingTimeInPhase float64 `lua:"remainingtimeinphase"`
	MooomPhaseCycle      int     `lua:"mooomphasecycle"`
	Segs                 Segs    `lua:"segs"`
}

type Segs struct {
	Night int `lua:"night"`
	Day   int `lua:"day"`
	Dusk  int `lua:"dusk"`
}

type IsRandom struct {
	Summer bool `lua:"summer"`
	Autumn bool `lua:"autumn"`
	Spring bool `lua:"spring"`
	Winter bool `lua:"winter"`
}

type Lengths struct {
	Summer int `lua:"summer"`
	Autumn int `lua:"autumn"`
	Spring int `lua:"spring"`
	Winter int `lua:"winter"`
}

type Seasons struct {
	Premode               bool                   `lua:"premode"`
	Season                string                 `lua:"season"`
	ElapsedDaysInSeason   int                    `lua:"elapseddaysinseason"`
	IsRandom              IsRandom               `lua:"israndom"`
	Lengths               Lengths                `lua:"lengths"`
	RemainingDaysInSeason int                    `lua:"remainingdaysinseason"`
	Mode                  string                 `lua:"mode"`
	TotalDaysInSeason     int                    `lua:"totaldaysinseason"`
	Segs                  map[string]interface{} `lua:"segs"`
}

type Data struct {
	Clock   Clock   `lua:"clock"`
	Seasons Seasons `lua:"seasons"`
}

func mapTableToStruct(table *lua.LTable, v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("lua")
		if tag == "" {
			continue
		}

		value := table.RawGetString(tag)
		if value == lua.LNil {
			continue
		}

		fieldValue := v.Field(i)
		switch fieldValue.Kind() {
		case reflect.Struct:
			if err := mapTableToStruct(value.(*lua.LTable), fieldValue); err != nil {
				return err
			}
		case reflect.Map:
			mapType := fieldValue.Type()
			keyType := mapType.Key()
			valueType := mapType.Elem()
			mapValue := reflect.MakeMap(mapType)

			tableValue, ok := value.(*lua.LTable)
			if !ok {
				return fmt.Errorf("expected lua.LTable, got %T", value)
			}

			tableValue.ForEach(func(key, value lua.LValue) {
				mapKey := reflect.New(keyType).Elem()
				err := luaValueToValue(key, mapKey)
				if err != nil {
					return
				}

				mapValue := reflect.New(valueType).Elem()
				err = luaValueToValue(value, mapValue)
				if err != nil {
					return
				}

				mapValue.Set(mapValue)
			})

			fieldValue.Set(mapValue)
		default:
			if err := luaValueToValue(value, fieldValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func luaValueToValue(lv lua.LValue, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if lv.Type() != lua.LTNumber {
			return fmt.Errorf("expected lua.LNumber, got %T", lv)
		}
		v.SetInt(int64(lv.(lua.LNumber)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if lv.Type() != lua.LTNumber {
			return fmt.Errorf("expected lua.LNumber, got %T", lv)
		}
		v.SetUint(uint64(lv.(lua.LNumber)))
	case reflect.Float32, reflect.Float64:
		if lv.Type() != lua.LTNumber {
			return fmt.Errorf("expected lua.LNumber, got %T", lv)
		}
		v.SetFloat(float64(lv.(lua.LNumber)))
	case reflect.Bool:
		if lv.Type() != lua.LTBool {
			return fmt.Errorf("expected lua.LBool, got %T", lv)
		}
		v.SetBool(bool(lv.(lua.LBool)))
	case reflect.String:
		if lv.Type() != lua.LTString {
			return fmt.Errorf("expected lua.LString, got %T", lv)
		}
		v.SetString(string(lv.(lua.LString)))
	case reflect.Map:
		if lv.Type() != lua.LTTable {
			return fmt.Errorf("expected lua.LTable, got %T", lv)
		}

		mapType := v.Type()
		keyType := mapType.Key()
		valueType := mapType.Elem()
		mapValue := reflect.MakeMap(mapType)

		table := lv.(*lua.LTable)
		table.ForEach(func(key, value lua.LValue) {
			mapKey := reflect.New(keyType).Elem()
			err := luaValueToValue(key, mapKey)
			if err != nil {
				return
			}

			mapValue := reflect.New(valueType).Elem()
			err = luaValueToValue(value, mapValue)
			if err != nil {
				return
			}

			mapValue.Set(mapValue)
		})

		v.Set(mapValue)
	default:
		return fmt.Errorf("unsupported kind: %v", v.Kind())
	}

	return nil
}

func mapTableToMap(table *lua.LTable, m map[string]interface{}) {
	table.ForEach(func(key, value lua.LValue) {
		keyStr := key.String()

		switch value.Type() {
		case lua.LTTable:
			nestedTable := value.(*lua.LTable)
			nestedMap := map[string]interface{}{}
			mapTableToMap(nestedTable, nestedMap)
			m[keyStr] = nestedMap
		case lua.LTNumber:
			m[keyStr] = float64(value.(lua.LNumber))
		case lua.LTString:
			m[keyStr] = string(value.(lua.LString))
		case lua.LTBool:
			m[keyStr] = bool(value.(lua.LBool))
		default:
			m[keyStr] = nil
		}
	})
}

func LuaTable2Map(script string) (map[string]interface{}, error) {
	L := lua.NewState()
	defer L.Close()

	err := L.DoString(script)
	if err != nil {
		return map[string]interface{}{}, err
	}
	table := L.Get(-1)
	L.Pop(1)
	data := map[string]interface{}{}
	mapTableToMap(table.(*lua.LTable), data)
	return data, nil
}

func LuaTable2Struct(script string, v reflect.Value) error {
	L := lua.NewState()
	defer L.Close()

	err := L.DoString(script)
	if err != nil {
		log.Println(err)
		return err
	}

	table := L.Get(-1)
	L.Pop(1)

	err = mapTableToStruct(table.(*lua.LTable), v)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
