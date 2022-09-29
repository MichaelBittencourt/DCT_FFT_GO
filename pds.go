package main

import (
	"fmt"
	"os"
)

var Version = "0.0.0"

func help() {
	fmt.Println("Usage:")
	fmt.Printf("    %s [ help | version | dct | fft ]\n", os.Args[0])
	fmt.Println("\t -h | help\tShow help message")
	fmt.Println("\t -v | version\tShow current version")
	fmt.Println("\t dct\t\tSet to use dct")
	fmt.Println("\t fft\t\tSet to use fft")
	fmt.Println()
	fmt.Println("    Eg:")
	fmt.Printf("\t%s help\n", os.Args[0])
	fmt.Printf("\t%s -h\n", os.Args[0])
	fmt.Printf("\t%s version\n", os.Args[0])
	fmt.Printf("\t%s -v\n", os.Args[0])
	fmt.Printf("\t%s dct\n", os.Args[0])
	fmt.Printf("\t%s fft\n", os.Args[0])
}

func version() {
	fmt.Printf("Go Test Log Version: %s\n", Version)
}

func invalidParam(param string) {
	fmt.Printf("Invalid Param: %s\n", param)
	fmt.Printf("Try: %s help\n", os.Args[0])
}

func main() {
	dct_flag := false
	fft_flag := false
	argv := os.Args[1:]
	for _, item := range argv {
		switch item {
		case "version", "-v":
			version()
			os.Exit(0)
		case "help", "-h":
			help()
			os.Exit(0)
		case "dct":
			dct_flag = true
		case "fft":
			fft_flag = true
		default:
			invalidParam(item)
			os.Exit(1)
		}
	}
	if dct_flag {
		fmt.Println("Run dct function")
		dct()
	}
	if fft_flag {
		fmt.Println("Run fft function")
	}
}

func getValue() (float64, int, error) {
	var ret float64
	qtdRead, err := fmt.Scanf("%f", &ret)
	return ret, qtdRead, err
}

func readData() []float64 {
	var data []float64
	for {
		value, qtdRead, _ := getValue()
		if qtdRead == 0 {
			break
		}
		data = append(data, value)
	}
	return data
}

func printData(data []float64) {
	for i := 0; i < len(data); i++ {
		fmt.Printf("data[%d]: %.2f\n", i, data[i])
	}
}

func dct() {
	data := readData()
	printData(data)
}
