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

func (p *Histogram) DeleteEmptyKeys() {
	for k, v := range p.Hist {
		if len(v.Pos) == 0 {
			delete(p.Hist, k)
		}
	}
}

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

// GetKeys extracts all p.Hist keys into p.Keys.
func (p *Histogram) GetKeys() {
	p.mu.Lock()
	for key := range p.Hist {
		p.Keys = append(p.Keys, key)
	}
	p.mu.Unlock()
}

// ComputeValues multiplies weight values with key len.
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

// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses the key positions.
func (p *Histogram) Reduce() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Hist) < 2 { // less than 2 keys
		return
	}
	p.SortKeysByIncrSize()
	for i := 0; i < len(p.Keys)-1; { // iterate over by increasing length sorted keys
		var smallerKeys []string
		smallerLength := len(p.Keys[i]) // smallerLength is multiple of 2 and the actual (smallest) size
		for i < len(p.Keys)-1 && smallerLength == len(p.Keys[i]) {
			smallerKeys = append(smallerKeys, p.Keys[i]) // collect all keys with the same smaller size
			i++
		}
		k := i // Keep position
		var biggerKeys []string
		biggerLength := len(p.Keys[i]) // biggerLength is multiple of 2 and the aktual next bigger size
		for i < len(p.Keys) && biggerLength == len(p.Keys[i]) {
			biggerKeys = append(biggerKeys, p.Keys[i]) // collect all keys with the same next bigger size
			i++
		}
		if smallerLength < biggerLength { // == on last item
			fmt.Println(i, biggerLength, smallerLength)
			p.ReduceOverlappingKeys(biggerKeys, smallerKeys)
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

// positionsMatch returns all sub positions where bpos + idx and sub equal.
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
}

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

// DeletePositionsOfKey removes positions from key and reduces its weight by len(positions).
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
			v.Weight -= 1
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
