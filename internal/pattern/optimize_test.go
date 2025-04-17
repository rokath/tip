package pattern

import (
	"fmt"
	"sync"
	"testing"

	"github.com/tj/assert"
)

func TestHistogram_Reduce_2(t *testing.T) {
	var m sync.Mutex
	type fields struct {
		Hist map[string]Pattern
		mu   *sync.Mutex
		Keys []string
	}
	tests := []struct {
		name   string
		fields fields
		data   []byte
		maxLen int
		ring   bool
		reduce bool
		exp    map[string]Pattern
	}{ // test cases:
		{ // name fields                               data    01234567  max, ring, reduce,
			"", fields{map[string]Pattern{}, &m, nil}, []byte("xxAAAAyy"), 4, false, false,
			map[string]Pattern{ // expected
				s2h("xxAA"): {Pos: []int{0}},    //// xxAA     @ 0
				s2h("xAAA"): {Pos: []int{1}},    ////  xAAA    @ 1       <- when xAAA is big key, xAA@1 is removed but AA@2 needs to be restored then
				s2h("AAAA"): {Pos: []int{2}},    ////   AAAA   @ 2
				s2h("AAAy"): {Pos: []int{3}},    ////    AAAy  @ 3       <- when AAAy is big key, AAy@4 is removed but AA@4 needs to be restored then
				s2h("AAyy"): {Pos: []int{4}},    ////     AAyy @ 4
				s2h("xxA"):  {Pos: []int{0}},    //// xxA      @ 0
				s2h("xAA"):  {Pos: []int{1}},    ////  xAA     @ 1        <- when xAA is big key AA@2 is removed
				s2h("AAA"):  {Pos: []int{2, 3}}, ////   AAA    @ 2 = 2 3
				//                                 //    AAA   @ 3 = 2 3
				s2h("AAy"): {Pos: []int{4}},       //     AAy  @ 4        <- when AAy is big key AA@4 is removed
				s2h("Ayy"): {Pos: []int{5}},       //      Ayy @ 5
				s2h("xx"):  {Pos: []int{0}},       // xx       @ 0
				s2h("xA"):  {Pos: []int{1}},       //  xA      @ 1
				s2h("AA"):  {Pos: []int{2, 3, 4}}, //   AA     @ 2 = 2 3 4
				s2h("Ay"):  {Pos: []int{5}},       //      Ay  @ 5
				s2h("yy"):  {Pos: []int{6}},       //       yy @ 6
			},
		},

		//  // WITHOUT RESTORE
		//  { // name fields                               data    01234567  max, ring, reduce,
		//  	"", fields{map[string]Pattern{}, &m, nil}, []byte("xxAAAAyy"), 4, false, true,
		//  	map[string]Pattern{ // expected                              - removes -
		//  		s2h("xxAA"): {Pos: []int{0}}, // xxAA     @ 0            xxA@0 xAA@1
		//  		s2h("xAAA"): {Pos: []int{1}}, //  xAAA    @ 1            xAA@1 AAA@2                                 <- when xAAA is big key, xAA@1 is removed but AA@2 needs to be restored then
		//  		s2h("AAAA"): {Pos: []int{2}}, //   AAAA   @ 2            AAA@2 AAA@3
		//  		s2h("AAAy"): {Pos: []int{3}}, //    AAAy  @ 3            AAA@3 AAy@4                                 <- when AAAy is big key, AAy@4 is removed but AA@4 needs to be restored then
		//  		s2h("AAyy"): {Pos: []int{4}}, //     AAyy @ 4            AAy@4 Ayy@6
		//  		s2h("xxA"):  {Pos: []int{}},  // xxA      @ 0            xx@0  xA@1
		//  		s2h("xAA"):  {Pos: []int{}},  //  xAA     @ 1            xA@1  AA@2                                  <- when xAA is big key AA@2 is removed
		//  		s2h("AAA"):  {Pos: []int{}},  //   AAA    @ 2 = 2 3      AA@2  AA@3  AA@4
		//  		//                            //    AAA   @ 3 = 2 3
		//  		s2h("AAy"): {Pos: []int{}}, ////     AAy  @ 4            AA@4  Ay@5                                  <- when AAy is big key AA@4 is removed
		//  		s2h("Ayy"): {Pos: []int{}}, ////      Ayy @ 5            Ay@5  yy@6
		//  		s2h("xx"):  {Pos: []int{}}, //// xx       @ 0
		//  		s2h("xA"):  {Pos: []int{}}, ////  xA      @ 1
		//  		s2h("AA"):  {Pos: []int{}}, ////   AA     @ 2 = 2 3 4
		//  		s2h("Ay"):  {Pos: []int{}}, ////      Ay  @ 5
		//  		s2h("yy"):  {Pos: []int{}}, ////       yy @ 6
		//  	},
		//  },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Histogram{
				Hist: tt.fields.Hist,
				mu:   tt.fields.mu,
				Keys: tt.fields.Keys,
			}
			p.ScanData(tt.data, tt.maxLen, tt.ring)
			if tt.reduce {
				//p.ReduceFromSmallerSide()
			}

			for k, v := range tt.exp {
				fmt.Println(k, v.Pos, tt.fields.Hist[k].Pos)
			}
			for k, v := range tt.exp {
				assert.Equal(t, v.Pos, tt.fields.Hist[k].Pos)
			}
		})
	}
}

// generated: ////////////////////////////////
