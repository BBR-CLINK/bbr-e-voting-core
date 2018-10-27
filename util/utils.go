package util

import (
	"reflect"
	"bbrHack/node"
)

func SliceExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func NodeExists(nodeList node.NodeList, check node.Node) bool {
	for _, node := range nodeList.NodeList {
		if node.IP == check.IP && node.Port == check.Port {
			return true
		}
	}
	return false
}