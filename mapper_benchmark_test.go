package set_test

import (
	"encoding/json"
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

func BenchmarkMapper(b *testing.B) {
	rowsDecoded, size := loadBenchmarkMapperData(b)
	var jsonRows [][]byte
	for _, decoded := range rowsDecoded {
		if encoded, err := json.Marshal(decoded); err != nil {
			b.Fatalf("During json.Marshal: %v", err.Error())
		} else {
			jsonRows = append(jsonRows, encoded)
		}

	}
	//
	b.ResetTimer()
	//
	b.Run("json", func(b *testing.B) {
		type CustomerContact struct {
			Id    int    `json:"customer_id"`
			First string `json:"customer_first"`
			Last  string `json:"customer_last"`
		}
		type VendorContact struct {
			Id    int    `json:"vendor_contact_id"`
			First string `json:"vendor_contact_first"`
			Last  string `json:"vendor_contact_last"`
		}
		type Vendor struct {
			Id          int    `json:"vendor_id"`
			Name        string `json:"vendor_name"`
			Description string `json:"vendor_description"`
			VendorContact
		}
		type T struct {
			Id           int    `json:"id"`
			CreatedTime  string `json:"created_time"`
			ModifiedTime string `json:"modified_time"`
			Price        int    `json:"price"`
			Quantity     int    `json:"quantity"`
			Total        int    `json:"total"`

			CustomerContact
			Vendor
		}
		var jsonRow []byte
		var k int
		dest := make([]T, 100)
		for n := 0; n < b.N; n++ {
			k = n % size
			jsonRow = jsonRows[k]
			if err := json.Unmarshal(jsonRow, &dest[k]); err != nil {
				b.Fatalf("json unmarshal error %v", err.Error())
			}
		}
	})
	//
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		Common
		Timestamps // Not used but present anyways
		First      string
		Last       string
	}
	type Vendor struct {
		Common
		Timestamps  // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Common
		Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	b.Run("bound no rebind", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var bound set.BoundMapping
		var row MapperBenchmarkJsonRow
		var k int
		dest := make([]T, 100)
		for n := 0; n < b.N; n++ {
			k = n % size
			row = rowsDecoded[k]
			bound = mapper.Bind(&dest[k])
			//
			bound.Set("Id", row.Id)
			bound.Set("CreatedTime", row.CreatedTime)
			bound.Set("ModifiedTime", row.ModifiedTime)
			bound.Set("Price", row.Price)
			bound.Set("Quantity", row.Quantity)
			bound.Set("Total", row.Total)
			//
			bound.Set("Customer_Id", row.CustomerId)
			bound.Set("Customer_First", row.CustomerFirst)
			bound.Set("Customer_Last", row.CustomerLast)
			//
			bound.Set("Vendor_Id", row.VendorId)
			bound.Set("Vendor_Name", row.VendorName)
			bound.Set("Vendor_Description", row.VendorDescription)
			bound.Set("Vendor_Contact_Id", row.VendorContactId)
			bound.Set("Vendor_Contact_First", row.VendorContactFirst)
			bound.Set("Vendor_Contact_Last", row.VendorContactLast)
			//
			if err := bound.Err(); err != nil {
				b.Fatalf("Unable to set: %v", err.Error())
			}
		}
	})
	//
	b.Run("bound rebind", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var row MapperBenchmarkJsonRow
		var k int
		dest := make([]T, 100)
		bound := mapper.Bind(&dest[0])
		for n := 0; n < b.N; n++ {
			k = n % size
			row = rowsDecoded[k]
			bound.Rebind(&dest[k])
			//
			bound.Set("Id", row.Id)
			bound.Set("CreatedTime", row.CreatedTime)
			bound.Set("ModifiedTime", row.ModifiedTime)
			bound.Set("Price", row.Price)
			bound.Set("Quantity", row.Quantity)
			bound.Set("Total", row.Total)
			//
			bound.Set("Customer_Id", row.CustomerId)
			bound.Set("Customer_First", row.CustomerFirst)
			bound.Set("Customer_Last", row.CustomerLast)
			//
			bound.Set("Vendor_Id", row.VendorId)
			bound.Set("Vendor_Name", row.VendorName)
			bound.Set("Vendor_Description", row.VendorDescription)
			bound.Set("Vendor_Contact_Id", row.VendorContactId)
			bound.Set("Vendor_Contact_First", row.VendorContactFirst)
			bound.Set("Vendor_Contact_Last", row.VendorContactLast)
			//
			if err := bound.Err(); err != nil {
				b.Fatalf("Unable to set: %v", err.Error())
			}
		}
	})
	//
	b.Run("prepared", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var row MapperBenchmarkJsonRow
		var k int
		dest := make([]T, 100)
		prepared := mapper.Prepare(&dest[0])
		err := prepared.Plan(
			"Id", "CreatedTime", "ModifiedTime",
			"Price", "Quantity", "Total",
			"Customer_Id", "Customer_First", "Customer_Last",
			"Vendor_Id", "Vendor_Name", "Vendor_Description", "Vendor_Contact_Id", "Vendor_Contact_First", "Vendor_Contact_Last")
		if err != nil {
			b.Fatalf("error preparing plan %v", err.Error())
		}
		for n := 0; n < b.N; n++ {
			k = n % size
			row = rowsDecoded[k]
			prepared.Rebind(&dest[k])
			//
			//
			prepared.Set(row.Id)
			prepared.Set(row.CreatedTime)
			prepared.Set(row.ModifiedTime)
			prepared.Set(row.Price)
			prepared.Set(row.Quantity)
			prepared.Set(row.Total)
			//
			prepared.Set(row.CustomerId)
			prepared.Set(row.CustomerFirst)
			prepared.Set(row.CustomerLast)
			//
			prepared.Set(row.VendorId)
			prepared.Set(row.VendorName)
			prepared.Set(row.VendorDescription)
			prepared.Set(row.VendorContactId)
			prepared.Set(row.VendorContactFirst)
			prepared.Set(row.VendorContactLast)
			//
			if err := prepared.Err(); err != nil {
				b.Fatalf("Unable to set: %v", err.Error())
			}
		}
	})
}
func BenchmarkMapperBaseline(b *testing.B) {
	rows, size := loadBenchmarkMapperData(b)
	//
	type Person struct {
		Id    int
		First string
		Last  string
	}
	type Vendor struct {
		Id          int
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		Id           int
		CreatedTime  string
		ModifiedTime string
		Price        int
		Quantity     int
		Total        int
		Customer     Person
		Vendor       Vendor
	}
	//
	b.ResetTimer()
	//
	for k := 0; k < b.N; k++ {
		row := rows[k%size]
		dest := new(T)
		//
		dest.Id = row.Id
		dest.CreatedTime = row.CreatedTime
		dest.ModifiedTime = row.ModifiedTime
		dest.Price = row.Price
		dest.Quantity = row.Quantity
		dest.Total = row.Total
		//
		dest.Customer.Id = row.CustomerId
		dest.Customer.First = row.CustomerFirst
		dest.Customer.Last = row.CustomerLast
		//
		dest.Vendor.Id = row.VendorId
		dest.Vendor.Name = row.VendorName
		dest.Vendor.Description = row.VendorDescription
		dest.Vendor.Contact.Id = row.VendorContactId
		dest.Vendor.Contact.First = row.VendorContactFirst
		dest.Vendor.Contact.Last = row.VendorContactLast
		// TODO RM
		// if err := bound.Err(); err != nil {
		// 	b.Fatalf("Unable to set: %v", err.Error())
		// }
	}
}

// TODO RM
// func BenchmarkValue(b *testing.B) { // TODO MOVE TO DIFFERENT FILE
// 	type Common struct {
// 		Id int
// 	}
// 	type Timestamps struct {
// 		CreatedTime  string
// 		ModifiedTime string
// 	}
// 	type Person struct {
// 		*Common
// 		*Timestamps // Not used but present anyways
// 		First       string
// 		Last        string
// 	}
// 	type Vendor struct {
// 		*Common
// 		*Timestamps // Not used but present anyways
// 		Name        string
// 		Description string
// 		Contact     Person
// 	}
// 	type T struct {
// 		*Common
// 		*Timestamps
// 		dest.Id = row.Id
// 		dest.CreatedTime = row.CreatedTime
// 		dest.ModifiedTime = row.ModifiedTime
// 		dest.Price = row.Price
// 		dest.Quantity = row.Quantity
// 		dest.Total = row.Total
// 		//
// 		dest.Customer.Id = row.CustomerId
// 		dest.Customer.First = row.CustomerFirst
// 		dest.Customer.Last = row.CustomerLast
// 		//
// 		dest.Vendor.Id = row.VendorId
// 		dest.Vendor.Name = row.VendorName
// 		dest.Vendor.Description = row.VendorDescription
// 		dest.Vendor.Contact.Id = row.VendorContactId
// 		dest.Vendor.Contact.First = row.VendorContactFirst
// 		dest.Vendor.Contact.Last = row.VendorContactLast
// 	}
// }

// TODO RM
// func BenchmarkMapperPreparedMapping(b *testing.B) {
// 	rows, size := loadBenchmarkMapperData(b)
// 	//
// 	type Common struct {
// 		Id int
// 	}
// 	type Timestamps struct {
// 		CreatedTime  string
// 		ModifiedTime string
// 	}
// 	type Person struct {
// 		Common
// 		Timestamps // Not used but present anyways
// 		First      string
// 		Last       string
// 	}
// 	type Vendor struct {
// 		Common
// 		Timestamps  // Not used but present anyways
// 		Name        string
// 		Description string
// 		Contact     Person
// 	}
// 	type T struct {
// 		Common
// 		Timestamps
// 		//
// 		Price    int
// 		Quantity int
// 		Total    int
// 		//
// 		Customer Person
// 		Vendor   Vendor
// 	}
// 	//
// 	b.ResetTimer()
// 	//
// 	//
// 	dest := new(T)
// 	prepared := mapper.Prepare(&dest)
// 	err := prepared.Plan(
// 		"Id", "CreatedTime", "ModifiedTime",
// 		"Price", "Quantity", "Total",
// 		"Customer_Id", "Customer_First", "Customer_Last",
// 		"Vendor_Id", "Vendor_Name", "Vendor_Description", "Vendor_Contact_Id", "Vendor_Contact_First", "Vendor_Contact_Last")
// 	if err != nil {
// 		b.Fatalf("error preparing plan %v", err.Error())
// 	}
// 	//
// 	for k := 0; k < b.N; k++ {
// 		row := rows[k%size]
// 		dest = new(T)
// 		prepared.Rebind(&dest)
// 		//
// 		prepared.Set(row.Id)
// 		prepared.Set(row.CreatedTime)
// 		prepared.Set(row.ModifiedTime)
// 		prepared.Set(row.Price)
// 		prepared.Set(row.Quantity)
// 		prepared.Set(row.Total)
// 		//
// 		prepared.Set(row.CustomerId)
// 		prepared.Set(row.CustomerFirst)
// 		prepared.Set(row.CustomerLast)
// 		//
// 		prepared.Set(row.VendorId)
// 		prepared.Set(row.VendorName)
// 		prepared.Set(row.VendorDescription)
// 		prepared.Set(row.VendorContactId)
// 		prepared.Set(row.VendorContactFirst)
// 		prepared.Set(row.VendorContactLast)
// 		//
// 		if err := prepared.Err(); err != nil {
// 			b.Fatalf("Unable to set: %v", err.Error())
// 		}
// 	}
// }
