package main

import (
	"fmt"
	"math"
	"math/cmplx"
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

func getValue() (complex128, int, error) {
	var realRet float64
	var imgRet float64
	var ret complex128
	qtdRead, err := fmt.Scanf("%f %f", &realRet, &imgRet)
	ret = complex(realRet, imgRet)
	return ret, qtdRead, err
}

func readData() []complex128 {
	var data []complex128
	for {
		value, qtdRead, _ := getValue()
		if qtdRead == 0 {
			break
		} else if qtdRead < 2 {
			fmt.Println("input format invalid")
			os.Exit(1)
		}
		data = append(data, value)
	}
	return data
}

func printData(data []complex128) {
	for i := 0; i < len(data); i++ {
		fmt.Printf("data[%d]: ", i)
		fmt.Println(data[i])
	}
}

func dct() {
	data := readData()
	printData(data)
	X := dct_calc(data)
	printData(X)
}

func dct_calc(data []complex128) []complex128 {
	N := len(data)
	//var WN complex128 = cmplx.Exp((-2i * math.Pi) / N)
	//var X []complex128

	for k := 0; k < N; k++ {
		var acc complex128 = 0 + 0i
		for n := 0; n < N; n++ {
			//X[k] += data[n] * cmplx.Pow(WN, n*k)
			value := 2i
			cPI := complex128(math.Pi)
			nk := complex(float64(n*k), 0)
			NComplex := complex(float64(N), 0)
			acc += data[n] * cmplx.Exp((value*cPI*nk)/NComplex)
		}
		fmt.Println(acc)
		//append(X, acc)
	}
	return []complex128{}
}
