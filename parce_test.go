package testparcer

import (
	"log"
	"testing"
)

type testDefaultStruct struct {
	Name string `json:"name,required"`
	Age  int    `json:"age" default:"18"`
}

type testStructWithNested struct {
	Name   string            `json:"name,required"`
	Age    int               `json:"age" default:"18"`
	Parent testDefaultStruct `json:"parent,required"`
}

type testStructWithNestedMapOfPointers struct {
	Name   string                        `json:"name,required"`
	Age    int                           `json:"age" default:"18"`
	Parent map[string]*testDefaultStruct `json:"parent,required"`
}

type testStructWithNestedArray struct {
	Name   string              `json:"name,required"`
	Age    int                 `json:"age" default:"18"`
	Parent []testDefaultStruct `json:"parent,required"`
}

type testStructWithNestedMap struct {
	Name   string            `json:"name,required"`
	Age    int               `json:"age" default:"18"`
	Parent map[string]string `json:"parent,required"`
}

func Test_checkeRequiredFields(t *testing.T) {
	type args struct {
		target interface{}
		m      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test_1",
			args{
				&testDefaultStruct{
					Name: "Jin",
				},
				map[string]interface{}{
					"name": "Jin",
			},
			},
			false,
		},
		{
			"test_2",
			args{
				&testDefaultStruct{},
				map[string]interface{}{},
			},
			true,
		},
		{
			"test_3",
			args{
				&testStructWithNested{
					Name: "jin",
					Parent: testDefaultStruct{
						Name: "pa",
					},
				},
				map[string]interface{}{
					"name": "Jin",
					"parent": map[string]interface{}{
						"name": "pa",
			},
				},
			},
			false,
		},
		{
			"test_3.1",
			args{
				&testStructWithNested{
					Name: "jin",
					Parent: testDefaultStruct{
						Age: 200,
					},
				},
				map[string]interface{}{
					"name": "Jin",
					"parent": map[string]interface{}{
						"age": 200,
					},
				},
			},
			false,
		},
		{
			"test_4",
			args{
				&testStructWithNested{
					Name:   "jin",
					Parent: testDefaultStruct{},
				},
				map[string]interface{}{
					"name": "Jin",
			},
			},
			true,
		},
		{
			"test_5",
			args{
				&testStructWithNested{
					Name: "jin",
				},
				map[string]interface{}{},
			},
			true,
		},
		{
			"test_6",
			args{
				&testStructWithNestedMapOfPointers{
					Name: "jin",
					Parent: map[string]*testDefaultStruct{
						"mam": &testDefaultStruct{Name: "helen"},
						"pa":  &testDefaultStruct{Name: "jo"},
					},
				},
				map[string]interface{}{
					"name": "jin",
					"parent": map[string]interface{}{
						"mam": map[string]interface{}{
							"name": "helen",
						},
						"pa": map[string]interface{}{
							"name": "jo",
						},
					},
			},
			},
			false,
		},
		{
			"test_7",
			args{
				&testStructWithNestedMapOfPointers{
					Name: "jin",
					Parent: map[string]*testDefaultStruct{
						"mam": &testDefaultStruct{Name: "helen"},
						"pa":  &testDefaultStruct{},
					},
				},
				map[string]interface{}{
					"name": "Jin",
					"parent": map[string]interface{}{
						"mam": map[string]interface{}{
							"name": "helen",
						},
						"pa": map[string]interface{}{},
					},
				},
			},
			true,
		},
		{
			"test_8",
			args{
				&testStructWithNestedArray{
					Name: "jin",
					Parent: []testDefaultStruct{
						testDefaultStruct{Name: "helen"},
						testDefaultStruct{Name: "jo"},
					},
				},
				map[string]interface{}{
					"name": "Jin",
					"parent": []interface{}{
						map[string]interface{}{
							"name": "helen",
			},
						map[string]interface{}{
							"name": "jo",
						},
					},
				},
			},
			false,
		},
		{
			"test_9",
			args{
				&testStructWithNestedArray{
					Name: "jin",
					Parent: []testDefaultStruct{
						testDefaultStruct{Name: "helen"},
						testDefaultStruct{Age: 20},
					},
				},
				map[string]interface{}{
					"name": "Jin",
					"parent": []interface{}{
						map[string]interface{}{
							"name": "helen",
			},
						map[string]interface{}{
							"age": 20,
						},
					},
				},
			},
			true,
		},
		{
			"test_10",
			args{
				&testStructWithNestedMap{
					Name: "jin",
					Parent: map[string]string{
						"mam": "Lisa",
						"pa":  "jo",
					},
				},
				map[string]interface{}{
					"name": "Jin",
			},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkeRequiredFields(tt.args.target, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("checkeRequiredFields() error = %v, wantErr %v", err, tt.wantErr)
				log.Println(tt.name, err)
			} else {
				log.Println(tt.name, err)
			}

		})
	}
}

type testStructDefaultFields struct {
	F1 int8   `default:"127"`
	F2 string `default:"str"`
	F3 byte   `default:"16"`
}

type testStructDefaultFieldsWrongInt struct {
	F1 int    `default:"10a"`
	F2 string `default:"str"`
	F3 byte   `default:"16"`
}

type testStructDefaultFieldsWrongByte struct {
	F1 int    `default:"10"`
	F2 string `default:"str"`
	F3 byte   `default:"16a"`
}

type testStructDefaultFieldsWrongArr struct {
	F1 []int `default:"10"`
}

type testStructDefaultFieldsWrongMap struct {
	F1 map[string]string `default:"10"`
}

type testStructDefaultFieldsWrongStruct struct {
	F1 testStructDefaultFields `default:"10"`
}

type testStructDefaultFieldsStruct struct {
	F1 testStructDefaultFields
}

type testStructDefaultFieldsMap struct {
	F1 map[string]*testStructDefaultFields
}

func Test_setDefaultFields(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test_1",
			args{
				&testStructDefaultFields{},
			},
			false,
		},
		{
			"test_2",
			args{
				&testStructDefaultFieldsWrongInt{},
			},
			true,
		},
		{
			"test_3",
			args{
				&testStructDefaultFieldsWrongByte{},
			},
			true,
		},
		{
			"test_4",
			args{
				&testStructDefaultFieldsWrongArr{},
			},
			true,
		},
		{
			"test_5",
			args{
				&testStructDefaultFieldsWrongMap{},
			},
			true,
		},
		{
			"test_6",
			args{
				&testStructDefaultFieldsWrongStruct{},
			},
			true,
		},
		{
			"test_7",
			args{
				&testStructDefaultFieldsStruct{},
			},
			false,
		},
		{
			"test_8",
			args{
				&testStructDefaultFieldsMap{
					F1: map[string]*testStructDefaultFields{
						"foo": &testStructDefaultFields{},
						"bar": &testStructDefaultFields{},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setDefaultFields(tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("setDefaultFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type testStruct struct {
	F1 byte                           `json:"byte-field"`                           // обычное поле без новых тэгов; опциональное, так как нет ключа required
	F2 string                         `json:"string-field,required"`                // строковое поле; обязательное, так как есть ключ required
	F3 int                            `json:"int1-field" default:"123"`             // числовое поле со значением по умолчанию 123
	F4 int                            `json:"int2-field" default:"foo"`             // числовое поле со значением по умолчанию foo (ошибка, некорректное значение)
	F5 []string                       `json:"slice-field,required" default:"what?"` // обязательное поле со значением по умолчанию (ошибка, срезы не могут иметь default)
	F6 anotherTestStruct1             `json:"struct-field,required"`                // подструктура, которая тоже может иметь тэги и т.д.
	F7 map[string]int                 `json:"primitive-map-field"`                  // обычное отображение без новых тэгов
	F8 map[string]*anotherTestStruct2 `json:"struct-map-field,required"`            // отображение на указатели подструктур, которые тоже могут иметь тэги и т.д.
}

type anotherTestStruct1 struct {
	F1 byte   `json:"byte-field"`            // обычное поле без новых тэгов; опциональное, так как нет ключа required
	F2 string `json:"string-field,required"` // строковое поле; обязательное, так как есть ключ required
	F3 int    `json:"int1-field" default:"123"`
}

type anotherTestStruct2 struct {
	F1 byte   `json:"byte-field"`            // обычное поле без новых тэгов; опциональное, так как нет ключа required
	F2 string `json:"string-field,required"` // строковое поле; обязательное, так как есть ключ required
	F3 int    `json:"int1-field" default:"123"`
}

type anotherTestStruct3 struct {
	F5 []anotherTestStruct1 `json:"slice-field,required"`
}

func TestParce(t *testing.T) {
	type args struct {
		filepath string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test_1 string into int",
			args{
				filepath: "test1.json",
				target:   &testStruct{},
			},
			true,
		},
		{
			"test_2 required field F5",
			args{
				filepath: "test2.json",
				target:   &testStruct{},
			},
			true,
		},
		{
			"test_3",
			args{
				filepath: "test3.json",
				target:   &testStruct{},
			},
			false,
		},
		{
			"test_4 F4 wrong default",
			args{
				filepath: "test4.json",
				target:   &testStruct{},
			},
			true,
		},
		{
			"test_5",
			args{
				filepath: "test5.json",
				target:   &anotherTestStruct3{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Parce(tt.args.filepath, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_readJSON(t *testing.T) {
	type args struct {
		filepath string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test_3",
			args{
				filepath: "test3.json",
				target:   &testStruct{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := readJSON(tt.args.filepath, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("readJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
