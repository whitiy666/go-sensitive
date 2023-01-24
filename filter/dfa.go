package filter

type node struct {
	children map[rune]*node
	isLeaf   bool
}

func newNode() *node {
	return &node{
		children: make(map[rune]*node),
		isLeaf:   false,
	}
}

type DfaModel struct {
	root *node
}

func NewDfaModel() *DfaModel {
	return &DfaModel{
		root: newNode(),
	}
}

func (m *DfaModel) AddWords(words ...string) {
	for _, word := range words {
		m.AddWord(word)
	}
}

func (m *DfaModel) AddWord(word string) {
	now := m.root
	runes := []rune(word)

	for _, r := range runes {
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			next = newNode()
			now.children[r] = next
			now = next
		}
	}

	now.isLeaf = true
}

func (m *DfaModel) DelWords(words ...string) {
	for _, word := range words {
		m.DelWords(word)
	}
}

func (m *DfaModel) DelWord(word string) {
	var lastLeaf *node
	var lastLeafNextRune rune
	now := m.root
	runes := []rune(word)

	for _, r := range runes {
		if next, ok := now.children[r]; !ok {
			return
		} else {
			if now.isLeaf {
				lastLeaf = now
				lastLeafNextRune = r
			}
			now = next
		}
	}

	delete(lastLeaf.children, lastLeafNextRune)
}

func (m *DfaModel) Listen(addChan, delChan <-chan string) {
	go func() {
		for word := range addChan {
			m.AddWord(word)
		}
	}()

	go func() {
		for word := range delChan {
			m.DelWord(word)
		}
	}()
}

func (m *DfaModel) Filter(text string) string {
	var found bool
	var now *node

	start := 0 // 从文本的第几个文字开始匹配
	parent := m.root
	runes := []rune(text)
	length := len(runes)
	filtered := make([]rune, 0, length)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			filtered = append(filtered, runes[start])
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf {
			start = pos + 1
			parent = m.root
		} else {
			parent = now
		}
	}

	filtered = append(filtered, runes[start:]...)

	return string(filtered)
}

func (m *DfaModel) Replace(text string, repl rune) string {
	var found bool
	var now *node

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf && start <= pos {
			for i := start; i <= pos; i++ {
				runes[i] = repl
			}
		}

		parent = now
	}

	return string(runes)
}

func (m *DfaModel) IsSensitive(text string) bool {
	return m.FindOne(text) != ""
}

func (m *DfaModel) FindOne(text string) string {
	var found bool
	var now *node

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found || (!now.isLeaf && pos == length-1) {
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf && start <= pos {
			return string(runes[start : pos+1])
		}

		parent = now
	}

	return ""
}

func (m *DfaModel) FindAll(text string) []string {
	var matches []string
	var found bool
	var now *node

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf && start <= pos {
			matches = append(matches, string(runes[start:pos+1]))
		}

		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}

		parent = now
	}

	var res []string
	set := make(map[string]struct{})

	for _, word := range matches {
		if _, ok := set[word]; !ok {
			set[word] = struct{}{}
			res = append(res, word)
		}
	}

	return res
}

func (m *DfaModel) FindAllCount(text string) map[string]int {
	res := make(map[string]int)
	var found bool
	var now *node

	start := 0
	parent := m.root
	runes := []rune(text)
	length := len(runes)

	for pos := 0; pos < length; pos++ {
		now, found = parent.children[runes[pos]]

		if !found {
			parent = m.root
			pos = start
			start++
			continue
		}

		if now.isLeaf && start <= pos {
			res[string(runes[start:pos+1])]++
		}

		if pos == length-1 {
			parent = m.root
			pos = start
			start++
			continue
		}

		parent = now
	}

	return res
}
