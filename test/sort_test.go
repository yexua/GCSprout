// @Author : Lik
// @Time   : 2021/7/8
package main

import (
	"fmt"
	"testing"
)

type LinkNode struct {
	val  int
	next *LinkNode
}

func TestMerge(t *testing.T) {

	var l1 LinkNode
	l1.val = 1
	var l3 LinkNode
	l3.val = 3
	var l5 LinkNode
	l5.val = 5
	var l7 LinkNode
	l7.val = 7

	l1.next = &l3
	l3.next = &l5
	l5.next = &l7

	var l2 LinkNode
	l2.val = 2
	var l4 LinkNode
	l4.val = 4
	var l6 LinkNode
	l6.val = 6
	var l8 LinkNode
	l8.val = 8

	l2.next = &l4
	l4.next = &l6
	l6.next = &l8

	res := merge(&l1, &l2)
	printList(res)
}

func printList(list *LinkNode) {
	node := list
	for node != nil {
		fmt.Printf("%d", node.val)
		node = node.next
	}
}

func merge(node1, node2 *LinkNode) *LinkNode {

	if node1 == nil {
		return node2
	}
	if node2 == nil {
		return node1
	}
	var head *LinkNode
	if node1.val >= node2.val {
		head = node2
		head.next = merge(node1, node2.next)
	} else {
		head = node1
		head.next = merge(node1.next, node2)
	}
	return head
}
