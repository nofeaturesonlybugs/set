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
	b.Run("Bind no Rebind", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var bound set.BoundMapping
		var row MapperBenchmarkJsonRow
		var k int
		var err error
		dest := make([]T, 100)
		for n := 0; n < b.N; n++ {
			k = n % size
			row = rowsDecoded[k]
			bound, err = mapper.Bind(&dest[k])
			if err != nil {
				b.Fatalf("Unable to bind: %v", err.Error())
			}
			//
			_ = bound.Set("Id", row.Id)
			_ = bound.Set("CreatedTime", row.CreatedTime)
			_ = bound.Set("ModifiedTime", row.ModifiedTime)
			_ = bound.Set("Price", row.Price)
			_ = bound.Set("Quantity", row.Quantity)
			_ = bound.Set("Total", row.Total)
			//
			_ = bound.Set("Customer_Id", row.CustomerId)
			_ = bound.Set("Customer_First", row.CustomerFirst)
			_ = bound.Set("Customer_Last", row.CustomerLast)
			//
			_ = bound.Set("Vendor_Id", row.VendorId)
			_ = bound.Set("Vendor_Name", row.VendorName)
			_ = bound.Set("Vendor_Description", row.VendorDescription)
			_ = bound.Set("Vendor_Contact_Id", row.VendorContactId)
			_ = bound.Set("Vendor_Contact_First", row.VendorContactFirst)
			_ = bound.Set("Vendor_Contact_Last", row.VendorContactLast)
			//
			if err := bound.Err(); err != nil {
				b.Fatalf("Unable to set: %v", err.Error())
			}
		}
	})
	//
	b.Run("Bind Rebind", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var row MapperBenchmarkJsonRow
		var k int
		dest := make([]T, 100)
		bound, err := mapper.Bind(&dest[0])
		if err != nil {
			b.Fatalf("Unable to bind: %v", err.Error())
		}
		for n := 0; n < b.N; n++ {
			k = n % size
			row = rowsDecoded[k]
			bound.Rebind(&dest[k])
			//
			_ = bound.Set("Id", row.Id)
			_ = bound.Set("CreatedTime", row.CreatedTime)
			_ = bound.Set("ModifiedTime", row.ModifiedTime)
			_ = bound.Set("Price", row.Price)
			_ = bound.Set("Quantity", row.Quantity)
			_ = bound.Set("Total", row.Total)
			//
			_ = bound.Set("Customer_Id", row.CustomerId)
			_ = bound.Set("Customer_First", row.CustomerFirst)
			_ = bound.Set("Customer_Last", row.CustomerLast)
			//
			_ = bound.Set("Vendor_Id", row.VendorId)
			_ = bound.Set("Vendor_Name", row.VendorName)
			_ = bound.Set("Vendor_Description", row.VendorDescription)
			_ = bound.Set("Vendor_Contact_Id", row.VendorContactId)
			_ = bound.Set("Vendor_Contact_First", row.VendorContactFirst)
			_ = bound.Set("Vendor_Contact_Last", row.VendorContactLast)
			//
			if err := bound.Err(); err != nil {
				b.Fatalf("Unable to set: %v", err.Error())
			}
		}
	})
	//
	b.Run("Prepare Rebind", func(b *testing.B) {
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Common{}, Timestamps{}),
			Join:     "_",
		}
		//
		var row MapperBenchmarkJsonRow
		var k int
		dest := make([]T, 100)
		prepared, err := mapper.Prepare(&dest[0])
		if err != nil {
			b.Fatalf("error preparing %v", err.Error())
		}
		err = prepared.Plan(
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
			_ = prepared.Set(row.Id)
			_ = prepared.Set(row.CreatedTime)
			_ = prepared.Set(row.ModifiedTime)
			_ = prepared.Set(row.Price)
			_ = prepared.Set(row.Quantity)
			_ = prepared.Set(row.Total)
			//
			_ = prepared.Set(row.CustomerId)
			_ = prepared.Set(row.CustomerFirst)
			_ = prepared.Set(row.CustomerLast)
			//
			_ = prepared.Set(row.VendorId)
			_ = prepared.Set(row.VendorName)
			_ = prepared.Set(row.VendorDescription)
			_ = prepared.Set(row.VendorContactId)
			_ = prepared.Set(row.VendorContactFirst)
			_ = prepared.Set(row.VendorContactLast)
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
	}
}
