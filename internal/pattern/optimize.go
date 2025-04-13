package pattern

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
)

func (p *Histogram) PrintInfo(message string) {
	smallest := math.MaxInt
	biggest := 0
	sum := 0
	count := 0
	for _, v := range p.Hist {
		size := len(v.Pos)
		smallest = min(smallest, size)
		biggest = max(biggest, size)
		sum += size
		count++
	}
	avr := float64(sum) / float64(count)
	fmt.Println("average=\t", avr, "\tsmallest:", smallest, "\tbiggest:", biggest, "\t", message)
}

/*
// ComputeBalance multiplies key length value with maxPatternSize / len(key)
// to achieve a balance in byte usage for pattern of different length.
// pattern   | counts | sum | factor | weight
// abcd      | 1	  | 1	| 4/1    | 4
// abc bcd   | 1+1	  | 2   | 4/2    | 4
// ab bc cd  | 1+1+1  | 3   | 4/3    | 4
// a b c d   | 1+1+1+1|	4   | 4/4    | 4
// aaaa      | 1	  | 1	| 4/1    | 4
// aaa aaa   | 1+1	  | 2   | 4/2    | 4
// aa aa aa  | 1+1+1  | 3   | 4/3    | 4
// a a a a   | 1+1+1+1|	4   | 4/4    | 4
// sum = max - (size - 1) = max+1-size
// weight = max/(max+1-size)
func (p *Histogram) ComputeBalance(maxPatternSize int) {
	for k, v := range p.Hist {
		x := maxPatternSize + 1 - len(k)/2 // longest pattern: x==1, 2-bytes pattern: x == maxPatternSize-1
		// Longer pattern get a bigger factor, shorter once a smaller factor.
		// This is the same as with simply using the pattern length but weaker.
		factor := float64(maxPatternSize) / float64(x)
		v.Balance = factor * float64(len(v.Bytes))
		p.Hist[k] = v
	}
}
*/

// GetKeys extracts all p.Hist keys into p.Keys.
func (p *Histogram) GetKeys() {
	p.mu.Lock()
	for key := range p.Hist {
		p.Keys = append(p.Keys, key)
	}
	p.mu.Unlock()
}

/*
// ComputeValues computes various values:
// weight: multiply position count with key len
// ...
func (p *Histogram) ComputeValues(maxPatternSize int) {
	p.GetKeys()
	p.ComputeBalance(maxPatternSize)
	for k, v := range p.Hist {
		v.Weight = float64(len(v.Pos) * len(v.Bytes))
		v.RateDirect = 1 / v.Weight
		v.RateIndirect = 2 / v.Weight
		p.Hist[k] = v
	}
}
*/

/*
ms@MacBook-Pro TRY % time ti_generate -i t1.txt -v -o idTable -u 7 -n 126 -z 3
67979 of 68290 patterns removed; 311 remaining,
OHNE REDUCE:
SortByDescCount
cnt:20041 w:  40082.0 b:    3.00 rateD:  0.0249 rateI:  0.0499 hex:            6162, ascci:'ab'
cnt:10231 w:  20462.0 b:    3.00 rateD:  0.0489 rateI:  0.0977 hex:            6263, ascci:'bc'
cnt:10185 w:  30555.0 b:    9.00 rateD:  0.0327 rateI:  0.0655 hex:          616263, ascci:'abc'
cnt: 628 w:   1256.0 b:    3.00 rateD:  0.7962 rateI:  1.5924 hex:            6261, ascci:'ba'
cnt: 601 w:   1202.0 b:    3.00 rateD:  0.8319 rateI:  1.6639 hex:            6361, ascci:'ca'

ms@MacBook-Pro TRY % time ti_generate -i t1.txt -v -o idTable -u 7 -n 126 -z 3
67979 of 68290 patterns removed; 311 remaining,
MIT REDUCE:
Reducing histogram with length 311 ...
622 6 4
Reducinging histogram...done. New length is 311
SortByDescCount
cnt:10185 w:  30555.0 b:    9.00 rateD:  0.0327 rateI:  0.0655 hex:          616263, ascci:'abc'
cnt: 430 w:   1290.0 b:    9.00 rateD:  0.7752 rateI:  1.5504 hex:          486162, ascci:'Hab'
cnt: 416 w:   1248.0 b:    9.00 rateD:  0.8013 rateI:  1.6026 hex:          616162, ascci:'aab'
cnt: 412 w:   1236.0 b:    9.00 rateD:  0.8091 rateI:  1.6181 hex:          6e6162, ascci:'nab'
cnt: 410 w:   1230.0 b:    9.00 rateD:  0.8130 rateI:  1.6260 hex:          4a6162, ascci:'Jab'

PROBLEM:
Pattern 'ab' count is 20041 and should get like 10000 after REDUCE and not like 0!!!!
The reason seems to be:
i: 81, cnt:   592, ascci:'xab'
i: 82, cnt:   589, ascci:'fab'
i: 83, cnt:   589, ascci:'Jab'
i: 84, cnt:   589, ascci:'Iab'
i: 85, cnt:   588, ascci:'Eab'
i: 86, cnt:   588, ascci:'aab'
These pattern do not survive but reduce ab!

Pattern 'bc' count is 10231 and should get like 0 after REDUCE, what is ok.
*/
// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses the key positions.
func (p *Histogram) ReduceFromSmallerSide() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Hist) < 2 { // less than 2 keys
		return
	}
	p.GetKeys()
	p.SortKeysByIncrSize()
	for i := 0; i < len(p.Keys)-1; { // iterate over by increasing length sorted keys
		var smallerKeys []string
		smallerLength := len(p.Keys[i]) // smallerLength is multiple of 2 and the actual (smallest) size.
		for i < len(p.Keys)-1 && smallerLength == len(p.Keys[i]) {
			smallerKeys = append(smallerKeys, p.Keys[i]) // collect all keys with the same smaller size.
			i++
		}
		k := i // Keep position
		var biggerKeys []string
		biggerLength := len(p.Keys[i]) // biggerLength is multiple of 2 and the aktual next bigger size.
		for i < len(p.Keys) && biggerLength == len(p.Keys[i]) {
			biggerKeys = append(biggerKeys, p.Keys[i]) // collect all keys with the same next bigger size.
			i++
		}
		if smallerLength < biggerLength { // == on last item
			fmt.Println(i, biggerLength, smallerLength, "(i, biggerLength, smallerLength)")
			p.ReduceOverlappingKeys(biggerKeys, smallerKeys)
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

/*
Without reduction:
SortByDescCount
i:  0, cnt: 30121, ascci:'ab'
i:  1, cnt: 20301, ascci:'bc'
i:  2, cnt: 20178, ascci:'abc'   removed
i:  3, cnt: 10315, ascci:'cd'
i:  4, cnt: 10200, ascci:'bcd'   removed
i:  5, cnt: 10197, ascci:'abcd'
i:  6, cnt: 10126, ascci:'RS'
i:  7, cnt: 10121, ascci:'TU'
i:  8, cnt: 10114, ascci:'ST'
i:  9, cnt: 10102, ascci:'xy'
i: 10, cnt: 10003, ascci:'RST'   removed
i: 11, cnt: 10000, ascci:'RSTU'
i: 12, cnt: 10000, ascci:'STU'   removed
i: 13, cnt: 10000, ascci:'123'   removed
i: 14, cnt: 10000, ascci:'12'
i: 15, cnt: 10000, ascci:'23'
i: 16, cnt:   923, ascci:'da'
i: 17, cnt:   911, ascci:'ba'
i: 18, cnt:   889, ascci:'Ua'
i: 19, cnt:   887, ascci:'ya'


If we reduce from larger size we get:

i:  0, cnt: 20185, ascci:'bc' stable
i:  1, cnt: 20178, ascci:'ab' -10000
i:  2, cnt: 10202, ascci:'cd' statble
i:  3, cnt: 10197, ascci:'abcd' stable
i:  4, cnt: 10114, ascci:'ST' stable
i:  5, cnt: 10008, ascci:'TU'
i:  6, cnt: 10008, ascci:'RS'
i:  7, cnt: 10000, ascci:'RSTU'
i:  8, cnt: 10000, ascci:'23'
i:  9, cnt: 10000, ascci:'12'
i: 10, cnt:   719, ascci:'da'
i: 11, cnt:   679, ascci:'Ua'
i: 12, cnt:   650, ascci:'ca'
i: 13, cnt:   548, ascci:'Ia'

So, the 3-byte sub-pattern are removed but not the 2-byte sub-pattern.

*/
func (p *Histogram) ReduceFromLargerSide() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Hist) < 2 { // less than 2 keys
		return
	}
	p.GetKeys()
	p.SortKeysByDescSize()
	for i := 0; i < len(p.Keys)-1; { // iterate over by decreasing length sorted keys
		var largerKeys []string
		largerLength := len(p.Keys[i]) // largerLength is multiple of 2 and the actual (smallest) size.
		for i < len(p.Keys)-1 && largerLength == len(p.Keys[i]) {
			largerKeys = append(largerKeys, p.Keys[i]) // collect all keys with the same larger size.
			i++
		}
		k := i // Keep position
		var smallerKeys []string
		smallerLength := len(p.Keys[i]) // smallerLength is multiple of 2 and the aktual next smaller size.
		for i < len(p.Keys) && smallerLength == len(p.Keys[i]) {
			smallerKeys = append(smallerKeys, p.Keys[i]) // collect all keys with the same next smaller size.
			i++
		}
		if largerLength > smallerLength { // == on last item
			fmt.Println(i, "maller:", smallerLength, "larger:", largerLength, "(i, smallerLength, largerLength)")
			p.ReduceOverlappingKeys(largerKeys, smallerKeys)
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

/*/ positionsMatch returns all sub positions where bpos + idx and sub equal.
func positionsMatch(bpos []int, idx int, sub []int) []int {
	pos := make([]int, 0)
	for _, x := range bpos {
		for _, y := range sub {
			if x+idx == y {
				pos = append(pos, y)
			}
		}
	}
	return pos
}*/

// positionIndexMatch return b pos if a + idx and b have one value common or -1.
func positionIndexMatch(a []int, idx int, b []int) int {
	for _, x := range a {
		for _, y := range b {
			if x+idx == y {
				return y
			}
		}
	}
	return -1
}

// DeletePositionsOfKey removes positions from key.
func (p *Histogram) DeletePositionsOfKey(key string, positions []int) {
	slices.Sort(positions)
	positions = slices.Compact(positions) // uniq
	v := p.Hist[key]
	n := 0
	for _, x := range v.Pos {
		if !slices.Contains(positions, x) {
			v.Pos[n] = x // keep
			n++
		} else {
			
		}
	}
	v.Pos = v.Pos[:n]
	p.Hist[key] = v
}

// getMatchingSubKeyPositions returns those subKey positions, which match appropriate bkey positions.
func (p *Histogram) getMatchingSubKeyPositions(bkey, subKey string) []int {
	var offset int
	subKeyPositions := make([]int, 0)
	for {
		s := bkey[offset:]
		idx := strings.Index(s, subKey) // get next subkey location inside bkey
		if idx == -1 {                  // not found
			return subKeyPositions
		}
		// "aaaaa"
		//subKey: "aa", 5,[]int{0,1,2,3,4}
		// bKey: "aaa", 5,[]int{0,1,2,3,4}
		// aa:0,aa:1                <- aaa:0
		//      aa:1,aa:2           <- aaa:1
		//           aa:2,aa:3      <- aaa:2
		//                aa:3,aa:4 <- aaa:3
		// aa:0,               aa:4 <- aaa:4
		// For each bkey are to remove 2 sunKey positions.
		// When subKey idx==0, we need to remove
		// We need to pay attention, that the .Pos values are 1-step values, but the keys are double sized!
		sKeyLoc := (offset + idx) / 2
		for _, x := range p.Hist[bkey].Pos { // x=01234
			for _, y := range p.Hist[subKey].Pos {
				if x+sKeyLoc == y {
					subKeyPositions = append(subKeyPositions, y)
					// break
				}
			}
		}
		offset += 2*sKeyLoc + 2
	}
}

// ReduceSubKey checks if subKey is inside bkey and removes the subKey internal positions,
// if they match with the bkey positions. Example: if subkey has positions 14, 18, 42 and bkey has
// position 10 and subkey is at index is 4 and 8 and, then the subkey positions 14, 18 are removed.
func (p *Histogram) ReduceSubKey(bkey, subKey string) {
	p.mu.Lock()
	pos := p.getMatchingSubKeyPositions(bkey, subKey)
	if len(pos) > 0 {
		p.DeletePositionsOfKey(subKey, pos)
	}
	p.mu.Unlock()
}

// ReduceOverlappingKeys checks for all biggerKeys if the smallerKeys are part of them
// and removes the subkey internal positions, if the positions are matching.
func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys []string) {
	var wg sync.WaitGroup
	for _, bkey := range biggerKeys {
		wg.Add(1)
		go func(bigKey string) {
			defer wg.Done()
			for _, subKey := range smallerKeys {
				p.ReduceSubKey(bkey, subKey)
			}
		}(bkey)
	}
	wg.Wait()
}
