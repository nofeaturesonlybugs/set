package set_test

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/nofeaturesonlybugs/set"
)

// This example demonstrates how Mapper might be used to create a general
// purpose CSV unmarshaler.

// CSVUnmarshaler is a general purpose CSV unmarshaler.
type CSVUnmarshaler struct {
	Mapper *set.Mapper
}

// Load reads CSV data from r and deserializes each row into dst.
func (u CSVUnmarshaler) Load(r io.Reader, dst interface{}) error {
	// Expect dst to be a *[]T or pointer chain to []T where T is struct or pointer chain to struct.
	slice, err := set.Slice(dst)
	if err != nil {
		return err
	} else if slice.ElemEndType.Kind() != reflect.Struct {
		return fmt.Errorf("dst elements should be struct or pointer to struct")
	}
	//
	m := u.Mapper
	if m == nil {
		// When the Mapper member is nil use a default mapper.
		m = &set.Mapper{
			Tags: []string{"csv"},
		}
	}
	//
	// Now the CSV can be read, processed, and deserialized into dst.
	c := csv.NewReader(r)
	//
	// Load the column names; these will be used later when calling BoundMapping.Set
	headers, err := c.Read()
	if err != nil {
		return err
	}
	c.ReuseRecord = true // From this point forward the csv reader can reuse the slice.
	//
	// Create a BoundMapping for element's type.
	b, err := m.Bind(slice.Elem())
	if err != nil {
		return err
	}
	for {
		row, err := c.Read()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		// Create a new element and bind to it.
		elemValue := slice.Elem() // Elem() returns reflect.Value and Rebind() conveniently
		b.Rebind(elemValue)       // allows reflect.Value as an argument.
		// Row is a slice of data and headers has the mapped names in corresponding indexes.
		for k, columnName := range headers {
			_ = b.Set(columnName, row[k]) // err will be checked after iteration.
		}
		// b.Err() returns the first error encountered from b.Set()
		if err = b.Err(); err != nil {
			return err
		}
		// Append to dst.
		slice.Append(elemValue)
	}
}

func Example_cSVUnmarshaler() {
	type Address struct {
		ID     int    `json:"id"`
		Street string `json:"street" csv:"address"`
		City   string `json:"city"`
		State  string `json:"state"`
		Zip    string `json:"zip" csv:"postal"`
	}
	PrintAddresses := func(all []Address) {
		for _, a := range all {
			fmt.Printf("ID=%v %v, %v, %v  %v\n", a.ID, a.Street, a.City, a.State, a.Zip)
		}
	}
	//
	loader := CSVUnmarshaler{
		Mapper: &set.Mapper{
			// Tags are listed in order of priority so csv has higher priority than json
			// when generating mapped names.
			Tags: []string{"csv", "json"},
		},
	}

	var data = `id,address,city,state,postal
1,06 Hoepker Court,Jacksonville,Florida,32209
2,92 Cody Hill,Falls Church,Virginia,22047
3,242 Burning Wood Terrace,Fort Worth,Texas,76105
4,41 Clarendon Pass,Fort Myers,Florida,33913
`
	var addresses []Address
	err := loader.Load(strings.NewReader(data), &addresses)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintAddresses(addresses)

	// Notice here this data has the columns in a different order
	fmt.Println()
	data = `city,address,postal,id,state
Tuscaloosa,2607 Hanson Junction,35487,1,Alabama
Bakersfield,2 Sherman Place,93305,2,California
Kansas City,4 Porter Place,64199,3,Missouri
New York City,23 Sachtjen Alley,10160,4,New York`
	addresses = addresses[0:0]
	err = loader.Load(strings.NewReader(data), &addresses)
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintAddresses(addresses)

	// Output: ID=1 06 Hoepker Court, Jacksonville, Florida  32209
	// ID=2 92 Cody Hill, Falls Church, Virginia  22047
	// ID=3 242 Burning Wood Terrace, Fort Worth, Texas  76105
	// ID=4 41 Clarendon Pass, Fort Myers, Florida  33913
	//
	// ID=1 2607 Hanson Junction, Tuscaloosa, Alabama  35487
	// ID=2 2 Sherman Place, Bakersfield, California  93305
	// ID=3 4 Porter Place, Kansas City, Missouri  64199
	// ID=4 23 Sachtjen Alley, New York City, New York  10160
}
