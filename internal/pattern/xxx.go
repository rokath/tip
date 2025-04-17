package pattern

/*

*/
/*
// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses the key positions.
func (p *Histogram) ReduceFromSmallerSide() {
	var lastSmallerKeys []string
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
			p.ReduceOverlappingKeys(biggerKeys, smallerKeys, lastSmallerKeys)
			lastSmallerKeys = smallerKeys
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}
*/
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

// DeletePositionsOfKey removes positions from key.
func (p *Histogram) DeletePositionsOfKey(key string, positions []int, lastSmallerKeys []string) {
	if len(positions) == 0 {
		return
	}

	// Before we delete the positions of key, we need to check if lastSmallerKeys
	// contain some deleted matching positions or in other words, if current key
	// caused a position deletion in lastSmallerKeys and we need to restore that position.
	//  for i := range lastSmallerKeys {
	//  	p.RestoreSubKey(lastSmallerKeys[i], key )
	//  }

	slices.Sort(positions)
	positions = slices.Compact(positions) // uniq
	v := p.Hist[key]
	n := 0
	for _, x := range v.Pos {
		if !slices.Contains(positions, x) {
			v.Pos[n] = x // keep
			n++
		} else {
			v.DeletedPos = append(v.DeletedPos, x)
		}
	}
	v.Pos = v.Pos[:n]
	p.Hist[key] = v
}

func (p *Histogram) RestorePositionsOfKey(key string, positions []int) {
	if len(positions) == 0 {
		return
	}

	v := p.Hist[key]
	for i, x := range v.DeletedPos {
		if slices.Contains(positions, x) {
			v.Pos = append(v.Pos, x)
		} else {
			v.DeletedPos[i] = -1 // invalidate deleted position
		}
	}
	p.Hist[key] = v
}

func (p *Histogram) RestoreSubKey(bkey, subKey string) {
	p.mu.Lock()
	pos := p.getMatchingSubKeyPositions(bkey, subKey)
	p.RestorePositionsOfKey(subKey, pos)
	p.mu.Unlock()
}

// ReduceSubKey checks if subKey is inside bkey and removes the subKey internal positions,
// if they match with the bkey positions. Example:
//
//	Xabc  Xabc         Yabc
//
// subkey:        abc found @       14    18           42
// bkey:         Xabc found @ 10
// subkey index inside bkey @        4     8
// then the subkey positions 14, 18 are removed.
func (p *Histogram) ReduceSubKey(bkey, subKey string, lastSmallerKeys []string) {
	p.mu.Lock()
	pos := p.getMatchingSubKeyPositions(bkey, subKey)
	p.DeletePositionsOfKey(subKey, pos, lastSmallerKeys)
	p.mu.Unlock()
}

// ReduceOverlappingKeys checks for all biggerKeys if the smallerKeys are part of them
// and removes the subkey internal positions, if the positions are matching.
func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys, lastSmallerKeys []string) {
	//var wg sync.WaitGroup
	for _, bkey := range biggerKeys {
		//wg.Add(1)
		go func(bigKey string) {
			//defer wg.Done()
			for _, subKey := range smallerKeys {
				p.ReduceSubKey(bkey, subKey, lastSmallerKeys)
			}
		}(bkey)
	}
	//wg.Wait()
}
*/

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

/*
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

/*
// SortKeysDescSize sorts p.Keys by decending size and alphabetical.
func (p *Histogram) SortKeysByDescSize() {
	compareFn := func(a, b string) int {
		if len(a) < len(b) {
			return 1
		}
		if len(a) > len(b) {
			return -1
		}
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
		return 0
	}
	slices.SortFunc(p.Keys, compareFn)
}
*/
/*
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
*/

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
/*
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
*/

/*
// countOverlapping returns sub count in s.
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping(s, sub string) int {
	var c int
	for i := range s {
		if strings.HasPrefix(s[i:], sub) {
			c++
		}
	}
	return c
}

// countOverlapping returns sub count in s.
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping(s, sub string) int {
	var c int
	for d := range s {
		if strings.HasPrefix(s[d:], sub) {
			c++
		}
	}
	return c
}

// SortByIncrLength returns list ordered for increasing pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncrLength(list []Patt) []Patt {
	compareFn := func(a, b Patt) int {
		if len(a.Bytes) > len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) < len(b.Bytes) {
			return -1
		}
		if a.Key > b.Key {
			return 1
		}
		if a.Key < b.Key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// merge copies all keys with their values from src into p.
// If p contains a key already, the values are added.
func (p *PatternHistogram) merge(src map[string]int) {
	for k, v := range src {
		p.mu.Lock()
		p.Hist[k] = p.Hist[k] + v
		p.mu.Unlock()
	}
}

// reduceSubCounts searches for p[i].Bytes being a part of an other p[k].Bytes with i < k.
// Example: If a pattern A is 3 times in pattern B, the pattern A.Cnt value is decreased by 3.
// Algorithm: check from small to big
func reduceSubCounts(p []Patt) []Patt {
	if Verbose {
		fmt.Println("Reducing sub pattern counts...")
	}
	if len(p) <= 1 {
		return p // nothing to do
	}
	list := SortByIncLength(p) // smallest pattern first

	count := getCounts(list) // get a copy to work on
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, x := range list[:len(list)-1] { // last list element is longest pattern
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			if Verbose {
				fmt.Println(k, "...")
			}
			sub := x.Bytes                 // sub is the next (smaller) pattern we want to check.
			for _, y := range list[k+1:] { // range over the next patterns
				n := slice.Count(y.Bytes, sub)
				if n > 0 {
					mu.Lock()
					count[k] -= n * y.Cnt
					mu.Unlock()
				}
			}
			if Verbose {
				fmt.Println(k, "...done")
			}
		}(i)
	}
	wg.Wait()
	setCounts(list, count)
	if Verbose {
		fmt.Println("Reducing sub pattern counts...done")
	}
	return list
}

func getCounts(list []Patt) []int {
	count := make([]int, len(list))
	for i, x := range list {
		count[i] = x.Cnt
	}
	return count
}

func setCounts(list []Patt, count []int) {
	for i := range list {
		list[i].Cnt = count[i]
	}
}

func GenerateDescendingCountSortedList(data []byte, maxPatternSize int) []Patt {
	m := BuildHistogram(data, maxPatternSize)
	list := histogramToList(m)
	//rList := list // reduceSubCounts(list)
	//sList := SortByDescCountDescLength(rList)
	return list // biggest cnt first, biggest length first on equal cnt
}

// BuildHistogram searches data for any 2-to-max bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func BuildHistogram(data []byte, max int) map[string]int {
	if Verbose {
		fmt.Println("Building histogram...")
	}
	subMap := make([]map[string]int, max) // maps slice
	var wg sync.WaitGroup
	for i := 0; i < max-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			subMap[k] = scan_ForRepetitions(data, k+2)
		}(i)
	}
	wg.Wait()
	m := make(map[string]int, 100000)
	for i := 0; i < max; i++ { // loop over pattern sizes
		maps.Copy(m, subMap[i])
	}
	if Verbose {
		fmt.Println("Building histogram...done. Length is", len(m))
	}
	return m
}
*/
