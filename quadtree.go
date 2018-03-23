package main



type Node struct {
	nw, ne, sw, se *Node
	level int
}

var BaseNodes [2]*Node = [2]*Node{&Node{nil, nil, nil, nil, 0}, &Node{nil, nil, nil, nil, 0}}

var GetNode func(nw, ne, sw, se *Node) *Node = MemoizeNode()

func Live(n *Node) int {
	if n==BaseNodes[1] { 
		return 1
	} 
	return 0
}

func MemoizeNode() func(nw, ne, sw, se *Node) *Node {
	NodeMemoized := make(map[[4]*Node]*Node)
	for _, n1 := range BaseNodes {
		for _, n2 := range BaseNodes {
			for _, n3 := range BaseNodes {
				for _, n4 := range BaseNodes {
					NodeMemoized[[4]*Node{n1,n2,n3,n4}] = &Node{n1,n2,n3,n4,1}
				}
			}
		}
	}
	return func(nw, ne, sw, se *Node) *Node {
		key := [4]*Node{nw,ne,sw,se}
		if n, ok := NodeMemoized[key]; ok {
			return n
		}
		n := &Node{nw,ne,sw,se,nw.level+1}
		NodeMemoized[key] = n
		return n
	}
}

func CenterSubnode(n *Node) *Node {
	return GetNode(n.nw.se, n.ne.sw, n.sw.ne, n.se.nw)
}

func CenterHorizontalSubnode(w, e *Node) *Node {
	return GetNode(w.ne.se, e.nw.sw, w.se.ne, e.sw.nw)
}

func CenterVerticalSubnode(n, s *Node) *Node {
	return GetNode(n.sw.se, n.se.sw, s.nw.ne, s.ne.nw)
}

func CenterSubSubnode(n *Node) *Node {
   return GetNode(n.nw.se.se, n.ne.sw.sw, n.sw.ne.ne, n.se.nw.nw)
}

//Cell is the same as Base Node
//m represent 3-by-3 with bits 012,456,89 10
func NextGenCell(m int) *Node { 
	if m==0 {
		return BaseNodes[0]
	}
	center := (m >> 5) & 1
	m = m & 0x757
	count := 0
	for m != 0 {
		count++
		m = m >> 1
	}
	if (count==3) || (count==2 && center!=0) {
		return BaseNodes[1]
	} else {
		return BaseNodes[0]
	}
	
}

func NextGenSubnode(n *Node) *Node {
	if n.level == 2 { //4-by-4
		ns := []*Node{
			n.nw.nw, n.nw.ne, n.ne.nw, n.ne.ne,
			n.nw.sw, n.nw.se, n.ne.sw, n.ne.se,
			n.sw.nw, n.sw.ne, n.se.nw, n.se.ne,
			n.sw.sw, n.sw.se, n.se.sw, n.se.se }
		m := 0 //represent 4-by-4 with 16 least significant bits from top to bot - left to right
		for i:=0; i < 16; i++ {
			m = (m << 1) + Live(ns[i])
		}
		nw, ne := NextGenCell(m >> 5), NextGenCell(m >> 4)
		sw, se := NextGenCell(m >> 1), NextGenCell(m)
		return GetNode(nw, ne, sw, se)
		
	}
	n00 := CenterSubnode(n.nw)
	n01 := CenterHorizontalSubnode(n.nw, n.ne)
	n02 := CenterSubnode(n.ne)
	n10 := CenterVerticalSubnode(n.nw, n.sw)
	n11 := CenterSubSubnode(n)
	n12 := CenterVerticalSubnode(n.ne, n.se)
	n20 := CenterSubnode(n.sw)
	n21 := CenterHorizontalSubnode(n.sw, n.se)
	n22 := CenterSubnode(n.se)
	return GetNode(
		NextGenSubnode(GetNode(n00, n01, n10, n11)),
		NextGenSubnode(GetNode(n01, n02, n11, n12)),
		NextGenSubnode(GetNode(n10, n11, n20, n21)),
		NextGenSubnode(GetNode(n11, n12, n21, n22)))
}