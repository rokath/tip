package pattern



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
	slices.SortFunc(p.Key, compareFn)
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
