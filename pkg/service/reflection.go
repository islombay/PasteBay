package service

import "reflect"

func changeStructs(src interface{}, dest interface{}) {
	vSrc := reflect.ValueOf(src).Elem()
	vDest := reflect.ValueOf(dest).Elem()

	for i := 0; i < vSrc.NumField(); i++ {
		fieldSrc := vSrc.Field(i)
		fieldDest := vDest.FieldByName(vSrc.Type().Field(i).Name)

		if fieldDest.IsValid() && fieldDest.CanSet() {
			fieldDest.Set(fieldSrc)
		}
	}
}
