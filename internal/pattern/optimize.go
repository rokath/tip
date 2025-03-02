package pattern

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

func (p *Histogram) PrintInfo(message string) {
	var (
		smallest float64 = math.MaxFloat32
		biggest          = float64(0)
		sum      float64
		count    int
	)

	for _, v := range p.Hist {
		smallest = min(smallest, v.Weight)

		biggest = max(biggest, v.Weight)
		sum += v.Weight
		count++
	}
	fmt.Println("average=\t", sum/float64(count), "\tsmallest:", smallest, "\tbiggest:", biggest, "\t", message)
}

// BalanceByteUsage multiplies each key value with maxPatternSize / len(key)
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
func (p *Histogram) BalanceByteUsage(maxPatternSize int) {
	for k, v := range p.Hist {
		factor := float64(maxPatternSize) / float64((maxPatternSize+1)-(len(k)>>1))
		v.Weight *= factor
		p.Hist[k] = v
	}
}

// AddWeigths multiplies weight values with key len.
func (p *Histogram) AddWeigths() {
	for k, v := range p.Hist {
		v.Weight *= float64(len(k) >> 1)
		p.Hist[k] = v
	}
}

// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses
func (p *Histogram) Reduce() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Hist) < 2 { // less than 2 keys
		return
	}
	p.GetKeys()
	p.SortKeysByIncrSize()
	for i := 0; i < len(p.Key)-1; { // iterate over by increasing length sorted keys
		var smallerKeys []string
		smallerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key)-1 && smallerLength == len(p.Key[i]) {
			smallerKeys = append(smallerKeys, p.Key[i])
			i++
		}
		k := i // Keep position
		var biggerKeys []string
		biggerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key) && biggerLength == len(p.Key[i]) {
			biggerKeys = append(biggerKeys, p.Key[i])
			i++
		}
		if smallerLength < biggerLength { // == on last item
			p.ReduceOverlappingKeys(biggerKeys, smallerKeys)
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
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

//  // positionMatch return pos if a and b have one value common or -1.
//  func positionMatch(a, b []int) int {
//  	for _, x := range a {
//  		for _, y := range b {
//  			if x == y {
//  				return x
//  			}
//  		}
//  	}
//  	return -1
//  }

//  // DeletePosition removes position key and reduces its weight by 1.
//  func (p *Histogram)DeletePosition(key string, position int){
//  	v := p.Hist[key]
//  	v.Pos = slices.DeleteFunc(v.Pos, func(position int) bool {
//  		return slices.Contains(v.Pos, position)
//  	})
//  	v.Weight--
//  	p.Hist[key] = v
//  }

// DeletePosition removes position key and reduces its weight by 1.
func (p *Histogram) DeletePosition(key string, position int) {
	v := p.Hist[key]
	for i, x := range v.Pos {
		if x == position {
			// if Verbose{
			// 	fmt.Println("DeletePosition:", key, v.Weight)
			// }
			v.Pos[i] = v.Pos[len(v.Pos)-1]
			v.Pos = v.Pos[:len(v.Pos)-1]
			v.Weight -= 1.0
			p.Hist[key] = v
			return
		}
	}
}

// ReduceSubKey checks if subKey is inside bkey and removes the subKey internal positions,
// if they match with the bkey positions. Example: if subkey has position 21 and bkey has
// position 18 and idx is 2, then the subkey position 20 is removed (18+2==20).
func (p *Histogram) ReduceSubKey(bkey, subKey string) {
	// We need to pay attention, that the .Pos values are 1-step values, but the keys are double sized!
	var offset int
repeat:
	s := bkey[offset:]
	idx := strings.Index(s, subKey)
	if idx == -1 { // subKey not inside bkey
		return
	}
	if idx&1 == 1 { // odd not allowed
		if offset+idx+1 < len(bkey) {
			offset += idx + 1
			goto repeat
		} else {
			return
		}
	}
	p.mu.Lock()
	bkeyPos := p.Hist[bkey].Pos
	subKeyPos := p.Hist[subKey].Pos
	p.mu.Unlock()
	pos := positionIndexMatch(bkeyPos, (offset+idx)>>1, subKeyPos)
	if pos >= 0 {
		p.mu.Lock()
		p.DeletePosition(subKey, pos)
		p.mu.Unlock()
	}
	if offset+idx+2 < len(bkey) {
		offset += idx + 2
		goto repeat
	}
}

// ReduceOverlappingKeys checks for all biggerKeys if the smallerKeys are part of them
// and removes their internal positions, if the positions are matching.
func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys []string) {
	//  if Verbose {
	//  	fmt.Println("ReduceOverlappingKeys")
	//  	fmt.Println(len(biggerKeys), "\tbiggerKeys:", biggerKeys[:3], "...")
	//  	fmt.Println(len(smallerKeys), "\tsmallerKeys:", smallerKeys[:3], "...")
	//  }
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

// countOverlapping2 returns sub count in s, assuming s & sub are hex-encoded byte buffers (even length).
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping2(s, sub string) int {
	var c int
	for i := 0; i < len(s); i += 2 {
		if strings.HasPrefix(s[i:], sub) {
			c++
		}
	}
	return c
}
