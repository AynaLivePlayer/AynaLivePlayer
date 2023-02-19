package provider

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLocal_Read(t *testing.T) {
	items, _ := ioutil.ReadDir(".")
	for _, item := range items {
		if item.IsDir() {
			subitems, _ := ioutil.ReadDir(item.Name())
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					// handle file there
					fmt.Println(item.Name() + "/" + subitem.Name())
				}
			}
		} else {
			// handle file there
			fmt.Println(item.Name())
		}
	}
}
