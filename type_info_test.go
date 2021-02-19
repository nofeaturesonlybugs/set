package set_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

func typeinfo_Invalid(i set.TypeInfo) bool {
	return i.IsMap == false && i.IsScalar == false && i.IsSlice == false && i.IsStruct == false && i.Kind == reflect.Invalid && i.Type == nil && i.ElemType == nil
}

func TestTypeInfo(t *testing.T) {
	chk := assert.New(t)
	{ // Test case that a map[reflect.Type]interface{} returns same value when reflect.Type is nil
		var err error
		var rw http.ResponseWriter
		var s string
		var sp *string
		var spp **string
		var t struct{}
		var tp *struct{}
		var tpp **struct{}
		var m map[string]string
		var mp *map[string]string
		var mpp **map[string]string
		var sl []struct{}
		var slp *[]struct{}
		var slpp *[]struct{}
		//
		var info set.TypeInfo
		size := 5
		ch := make(chan struct{})
		signals := []chan struct{}{}
		//
		for k := 0; k < size; k++ {
			signals = append(signals, make(chan struct{}))
			go func(idx int) {
				<-ch
				//
				info = set.TypeCache.Stat(nil)
				chk.Equal(true, typeinfo_Invalid(info))
				info = set.TypeCache.Stat(err)
				chk.Equal(true, typeinfo_Invalid(info))
				info = set.TypeCache.Stat(rw)
				chk.Equal(true, typeinfo_Invalid(info))
				//
				info = set.TypeCache.Stat(s)
				chk.Equal(true, info.IsScalar)
				chk.Equal(nil, info.ElemType)
				info = set.TypeCache.Stat(sp)
				chk.Equal(true, info.IsScalar)
				chk.Equal(nil, info.ElemType)
				info = set.TypeCache.Stat(spp)
				chk.Equal(true, info.IsScalar)
				chk.Equal(nil, info.ElemType)
				//
				info = set.TypeCache.Stat(t)
				chk.Equal(true, info.IsStruct)
				chk.Equal(nil, info.ElemType)
				info = set.TypeCache.Stat(tp)
				chk.Equal(true, info.IsStruct)
				chk.Equal(nil, info.ElemType)
				info = set.TypeCache.Stat(tpp)
				chk.Equal(true, info.IsStruct)
				chk.Equal(nil, info.ElemType)
				//
				info = set.TypeCache.Stat(m)
				chk.Equal(true, info.IsMap)
				chk.Equal(reflect.TypeOf(""), info.ElemType)
				info = set.TypeCache.Stat(mp)
				chk.Equal(true, info.IsMap)
				info = set.TypeCache.Stat(mpp)
				chk.Equal(true, info.IsMap)
				//
				info = set.TypeCache.Stat(sl)
				chk.Equal(true, info.IsSlice)
				chk.Equal(reflect.TypeOf(struct{}{}), info.ElemType)
				info = set.TypeCache.Stat(slp)
				chk.Equal(true, info.IsSlice)
				info = set.TypeCache.Stat(slpp)
				chk.Equal(true, info.IsSlice)
				//
				close(signals[idx])
			}(k)
		}
		close(ch)
		for k := 0; k < size; k++ {
			<-signals[k]
		}
	}
}

func BenchmarkTypeInfo(b *testing.B) {
	var err error
	var rw http.ResponseWriter
	var s string
	var sp *string
	var spp **string
	var t struct{}
	var tp *struct{}
	var tpp **struct{}
	var m map[string]string
	var mp *map[string]string
	var mpp **map[string]string
	var sl []struct{}
	var slp *[]struct{}
	var slpp *[]struct{}
	for k := 0; k < b.N; k++ {
		set.TypeCache.Stat(nil)
		set.TypeCache.Stat(err)
		set.TypeCache.Stat(rw)
		//
		set.TypeCache.Stat(s)
		set.TypeCache.Stat(sp)
		set.TypeCache.Stat(spp)
		//
		set.TypeCache.Stat(t)
		set.TypeCache.Stat(tp)
		set.TypeCache.Stat(tpp)
		//
		set.TypeCache.Stat(m)
		set.TypeCache.Stat(mp)
		set.TypeCache.Stat(mpp)
		//
		set.TypeCache.Stat(sl)
		set.TypeCache.Stat(slp)
		set.TypeCache.Stat(slpp)
	}
}

func BenchmarkTypeInfoParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var err error
		var rw http.ResponseWriter
		var s string
		var sp *string
		var spp **string
		var t struct{}
		var tp *struct{}
		var tpp **struct{}
		var m map[string]string
		var mp *map[string]string
		var mpp **map[string]string
		var sl []struct{}
		var slp *[]struct{}
		var slpp *[]struct{}
		for pb.Next() {
			set.TypeCache.Stat(nil)
			set.TypeCache.Stat(err)
			set.TypeCache.Stat(rw)
			//
			set.TypeCache.Stat(s)
			set.TypeCache.Stat(sp)
			set.TypeCache.Stat(spp)
			//
			set.TypeCache.Stat(t)
			set.TypeCache.Stat(tp)
			set.TypeCache.Stat(tpp)
			//
			set.TypeCache.Stat(m)
			set.TypeCache.Stat(mp)
			set.TypeCache.Stat(mpp)
			//
			set.TypeCache.Stat(sl)
			set.TypeCache.Stat(slp)
			set.TypeCache.Stat(slpp)
		}
	})
}
