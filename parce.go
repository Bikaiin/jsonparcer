package testparcer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrorWhileReadingFile     = errors.New("error while reading file")
	ErrorWhileUnmarshaling    = errors.New("error while unmarshaling")
	ErrorWhileChekingRequired = errors.New("error while cheking requiered fields")
	ErrorWhileSettingDefault  = errors.New("error while setting fields")
)

// Parce принимает путь файла и с труктуру в которю распарсит json, возвращает ошибку
func Parce(filepath string, target interface{}) error {
	err := readJSON(filepath, target)
	if err != nil {
		err := fmt.Errorf("%w: %v", ErrorWhileReadingFile, err)
		return err
	}

	err = checkeRequiredFields(target)
	if err != nil {
		err := fmt.Errorf("%w: %v", ErrorWhileChekingRequired, err)
		return err
	}

	err = setDefaultFields(target)
	if err != nil {
		err := fmt.Errorf("%w: %v", ErrorWhileUnmarshaling, err)
		return err
	}

	return nil
}

func checkeRequiredFields(target interface{}) error {
	fields := reflect.ValueOf(target).Elem()

	for i := 0; i < fields.NumField(); i++ {
		tagStr := fields.Type().Field(i).Tag.Get("json")

		if strings.Contains(tagStr, "required") && fields.Field(i).IsZero() {
			err := fmt.Sprintf(`required field "%v" (tag "%v") is missing`, fields.Type().Field(i).Name, strings.Split(tagStr, ",")[0])
			return errors.New(err)
		}

		switch fields.Field(i).Kind() {
		case reflect.Struct:
			err := checkeRequiredFields(fields.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		case reflect.Map:
			for _, key := range fields.Field(i).MapKeys() {
				if (fields.Field(i).MapIndex(key).Kind() == reflect.Struct ||
					fields.Field(i).MapIndex(key).Kind() == reflect.Ptr) && !fields.Field(i).MapIndex(key).IsZero() {
					err := checkeRequiredFields(fields.Field(i).MapIndex(key).Interface())
					if err != nil {
						return err
					}
				}
			}
		case reflect.Slice:
			for j := 0; j < fields.Field(i).Len(); j++ {
				if (fields.Field(i).Index(j).Kind() == reflect.Struct ||
					fields.Field(i).Index(j).Kind() == reflect.Ptr) && !fields.Field(i).Index(j).IsZero() {
					err := checkeRequiredFields(fields.Field(i).Index(j).Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func setDefaultFields(target interface{}) error {
	fields := reflect.ValueOf(target).Elem()

	for i := 0; i < fields.NumField(); i++ {
		tagStr := fields.Type().Field(i).Tag.Get("default")

		if tagStr != "" && fields.Field(i).IsZero() {
			switch fields.Field(i).Kind() {
			case reflect.Int:
				val, err := strconv.ParseInt(tagStr, 10, 32)
				if err != nil {
					return err
				}
				fields.Field(i).SetInt(val)
			case reflect.Int8:
				val, err := strconv.ParseInt(tagStr, 10, 8)
				if err != nil {
					return err
				}
				fields.Field(i).SetInt(val)
			case reflect.Int16:
				val, err := strconv.ParseInt(tagStr, 10, 16)
				if err != nil {
					return err
				}
				fields.Field(i).SetInt(val)
			case reflect.Int32:
				val, err := strconv.ParseInt(tagStr, 10, 32)
				if err != nil {
					return err
				}
				fields.Field(i).SetInt(val)
			case reflect.Int64:
				val, err := strconv.ParseInt(tagStr, 10, 64)
				if err != nil {
					return err
				}
				fields.Field(i).SetInt(val)
			case reflect.Uint:
				u, err := strconv.ParseUint(tagStr, 10, 32)
				if err != nil {
					return err
				}
				fields.Field(i).SetUint(u)
			case reflect.Uint8:
				u, err := strconv.ParseUint(tagStr, 10, 8)
				if err != nil {
					return err
				}
				fields.Field(i).SetUint(u)
			case reflect.Uint16:
				u, err := strconv.ParseUint(tagStr, 10, 16)
				if err != nil {
					return err
				}
				fields.Field(i).SetUint(u)
			case reflect.Uint32:
				u, err := strconv.ParseUint(tagStr, 10, 32)
				if err != nil {
					return err
				}
				fields.Field(i).SetUint(u)
			case reflect.Uint64:
				u, err := strconv.ParseUint(tagStr, 10, 64)
				if err != nil {
					return err
				}
				fields.Field(i).SetUint(u)
			case reflect.Float32:
				f, err := strconv.ParseFloat(tagStr, 32)
				if err != nil {
					return err
				}
				fields.Field(i).SetFloat(f)
			case reflect.Float64:
				f, err := strconv.ParseFloat(tagStr, 64)
				if err != nil {
					return err
				}
				fields.Field(i).SetFloat(f)
			case reflect.String:
				fields.Field(i).SetString(tagStr)
			default:
				err := fmt.Sprintf(`type of field "%v" (type %v) is not support setting defaul value`, fields.Type().Field(i).Name, fields.Type().Field(i).Type)
				return errors.New(err)
			}
		}

		switch fields.Field(i).Kind() {
		case reflect.Struct:
			err := setDefaultFields(fields.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		case reflect.Slice:
			for j := 0; j < fields.Field(i).Len(); j++ {
				if (fields.Field(i).Index(j).Kind() == reflect.Struct ||
					fields.Field(i).Index(j).Kind() == reflect.Ptr) && !fields.Field(i).Index(j).IsZero() {
					err := setDefaultFields(fields.Field(i).Index(j).Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		case reflect.Map:
			for _, key := range fields.Field(i).MapKeys() {
				if (fields.Field(i).MapIndex(key).Kind() == reflect.Struct ||
					fields.Field(i).MapIndex(key).Kind() == reflect.Ptr) && !fields.Field(i).MapIndex(key).IsZero() {
					err := setDefaultFields(fields.Field(i).MapIndex(key).Interface())
					if err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

func readJSON(filepath string, target interface{}) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewReader(f)
	d := json.NewDecoder(b)
	d.DisallowUnknownFields()
	err = d.Decode(target)
	if err != nil {
		return err
	}

	return nil
}
