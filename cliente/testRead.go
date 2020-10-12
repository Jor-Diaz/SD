package main

import (
  "os"
	"encoding/csv"
	"fmt"
	"io"
)

// Creación de Item struc
type Item struct{
  id string
  producto string
  valor string
  tienda string
  destino string
}

type Pyme struct{
  id string
  producto string
  valor string
  tienda string
  destino string
  prioritario string
}

var productos []*Item // Arreglo de Items (retail.csv)
var pymes []*Pyme // Arreglo de pymes (pymes.csv)

func RetailReader(){
	//Abir archivo
  recordFile, err := os.Open("retail.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		os.Exit(0)
	}
	reader := csv.NewReader(recordFile)
	reader.Read() //saltar primera linea
	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // final del archivo
		} else if err != nil {
			fmt.Println("Error ::", err)
			break
		}
    prod := Item{id: record[0], producto: record[1], valor:record[2], tienda:record[3],destino:record[4]}
    productos = append(productos, &prod)
  	}

}

func PymeReader(){
	//Abir archivo
  recordFile, err := os.Open("pymes.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		os.Exit(0)
	}
	reader := csv.NewReader(recordFile)
	reader.Read() //saltar primera linea
	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // final del archivo
		} else if err != nil {
			fmt.Println("Error ::", err)
			break
		}
    sme := Pyme{id:record[0],producto:record[1],valor:record[2],tienda:record[3],destino:record[4],prioritario:record[5]}
    pymes = append(pymes, &sme)
  	}

}

func main() {
  //Función para crear el array de estructuras Item
  RetailReader()
  PymeReader()
}
