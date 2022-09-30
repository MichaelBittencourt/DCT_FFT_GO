package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strconv"
)

var Version = "0.0.0"
var verbosity = false

func help() {
	fmt.Println("Usage:")
	fmt.Printf("    %s [ help | version | dct | fft ]\n", os.Args[0])
	fmt.Println("\t -h | help\tShow help message")
	fmt.Println("\t version\tShow current version")
	fmt.Println("\t -v\tSet verbosity")
	fmt.Println("\t dct\t\tSet to use dct")
	fmt.Println("\t inv\t\tSet to use invert transformation")
	fmt.Println()
	fmt.Println("    Eg:")
	fmt.Printf("\t%s help\n", os.Args[0])
	fmt.Printf("\t%s -h\n", os.Args[0])
	fmt.Printf("\t%s version\n", os.Args[0])
	fmt.Printf("\t%s -v\n", os.Args[0])

	fmt.Printf("\t%s dct\n", os.Args[0])
	fmt.Printf("\t%s inv\n", os.Args[0])
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
	invert_flag := false
	argv := os.Args[1:]
	for _, item := range argv {
		switch item {
		case "version":
			version()
			os.Exit(0)
		case "help", "-h":
			help()
			os.Exit(0)
		case "-v":
			verbosity = true
		case "dct":
			dct_flag = true
		case "inv":
			invert_flag = true
		default:
			invalidParam(item)
			os.Exit(1)
		}
	}
	data := readData()
	if verbosity {
		os.Stderr.WriteString("Data from stdin\n")
		printData(data)
	}
	var X []complex128
	if dct_flag {
		if verbosity {
			fmt.Println("Run dct function")
		}
		X = dct(data, invert_flag)
	} else {
		if verbosity {
			fmt.Println("Run fft function")
		}
		X = fft(data, invert_flag)
	}
	printOutput(X)
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
		fmt.Fprintf(os.Stderr, "data[%d]: ", i)
		os.Stderr.WriteString(strconv.FormatComplex(data[i], 'f', 0, 64) + "\n")
	}
}

func printOutput(data []complex128) {
	for i := 0; i < len(data); i++ {
		fmt.Printf("%f %f\n", real(data[i]), imag(data[i]))
	}
}

func dct(data []complex128, inverse bool) []complex128 {
	arraySize := len(data)
	N := complex(float64(arraySize), 0)
	var sign complex128
	PI := complex(math.Pi, 0)
	constInverse := complex(1, 0)
	if inverse {
		sign = 1
		constInverse /= N
	} else {
		sign = -1
	}

	var WN complex128 = cmplx.Exp((sign * 2i * PI) / N)
	var X = make([]complex128, arraySize)

	for k := 0; k < arraySize; k++ {
		X[k] = 0 + 0i
		for n := 0; n < arraySize; n++ {
			exponent := complex(float64(n*k), 0)
			X[k] += data[n] * cmplx.Pow(WN, exponent)
		}
		X[k] *= constInverse
	}
	return X
}

func fft(data []complex128, inverse bool) []complex128 {
	var X []complex128
	arraySize := len(data)
	if verbosity {
		fmt.Println("arraySize: %d", arraySize)
	}
	lgSize := math.Log2(float64(arraySize))
	if lgSize == float64(int(lgSize)) { //If lg(arraySize) is a integer value.
		if verbosity {
			fmt.Println("Is base 2")
		}
		N := complex(float64(arraySize), 0)
		var sign complex128
		PI := complex(math.Pi, 0)
		constInverse := complex(1, 0)
		if inverse {
			sign = 1
			constInverse /= N
		} else {
			sign = -1
		}
		var WN complex128 = cmplx.Exp((sign * 2i * PI) / N)
		if arraySize != 1 {
			var data1 []complex128
			var data2 []complex128
			for k := 0; k < arraySize; k++ {
				if k%2 == 0 {
					data1 = append(data1, data[k])
				} else {
					data2 = append(data2, data[k])
				}
			}
			G := fft(data1, inverse)
			if verbosity {
				fmt.Printf("G: %d\n", len(G))
			}
			H := fft(data2, inverse)
			if verbosity {
				fmt.Printf("H: %d\n", len(H))
			}
			GSize := len(G)
			for k := 0; k < arraySize; k++ {
				xk := constInverse * (G[k%GSize] + cmplx.Pow(WN, complex(float64(k), 0))*H[k%GSize])
				X = append(X, xk)
			}
		} else {
			X = data
		}
	}
	return X
}
