package main

import (
  "os"
	"encoding/csv"
	"fmt"
	"io"
	// "log"

)

type Item struct{
  id string
  producto string
  valor string
  tienda string
  destino string
}

func newItem(id_item string, item string, value string, store string, dest string) *Item{
  n_item := Item{id: id_item, producto: item, valor:value, tienda:store,destino:dest}
  return &n_item
}

func RetailReader(arr []*Item) {
	// Open the file
	recordFile, err := os.Open("retail.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	reader := csv.NewReader(recordFile)
	reader.Read() //saltar primera linea
	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("Error ::", err)
			return
		}
    prod := newItem(record[0],record[1],record[2],record[3],record[4])
    arr= append(arr, prod)
  	}
}



func main() {
  productos:=[]*Item{}
  RetailReader(productos)
  fmt.Println(productos[0])

}
