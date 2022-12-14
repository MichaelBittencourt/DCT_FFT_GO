package main

import (
    "fmt"
    "math"
    "math/cmplx"
    "os"
    "strconv"
    "unsafe"
)

var Version = "0.0.0"
var verbosity = false
var padding = false

func help() {
    fmt.Println("Usage:")
    fmt.Printf("    %s [ help | version | dft | fft ]\n", os.Args[0])
    fmt.Println("\t -h | help\tShow help message")
    fmt.Println("\t version\tShow current version")
    fmt.Println("\t -v\tSet verbosity")
    fmt.Println("\t dft\t\tSet to use dft")
    fmt.Println("\t inv\t\tSet to use invert transformation")
    fmt.Println()
    fmt.Println("    Eg:")
    fmt.Printf("\t%s help\n", os.Args[0])
    fmt.Printf("\t%s -h\n", os.Args[0])
    fmt.Printf("\t%s version\n", os.Args[0])
    fmt.Printf("\t%s -v\n", os.Args[0])

    fmt.Printf("\t%s dft\n", os.Args[0])
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
    dft_flag := false
    invert_flag := false
    frequency_decimation_flag := false
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
        case "dft":
            dft_flag = true
        case "inv":
            invert_flag = true
        case "freq":
            frequency_decimation_flag = true
        case "padding":
            padding = true
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
    if dft_flag {
        if verbosity {
            fmt.Println("Run dft function")
        }
        X = dft(data, invert_flag)
    } else {
        if verbosity {
            fmt.Println("Run fft function")
        }
        X = fft(data, invert_flag, frequency_decimation_flag)
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

func dft(data []complex128, inverse bool) []complex128 {
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
    WN_iterator := complex(1, 0) //Using this variable I can optimize the operation on the line cmplx.Pow(WN, n*k) because I can consider WN^k = WN^(k-1) * WN
    var X = make([]complex128, arraySize)

    for k := 0; k < arraySize; k++ {
        X[k] = 0 + 0i
        WNK := complex(1, 0) // With the same intention of WN_iterator I can calculate WNK^n = WNK^(n-1) * WN_iterator and I can remove the instruction cmplx.Pow(WN, n*k)
        for n := 0; n < arraySize; n++ {
            //exponent := complex(float64(n*k), 0)
            //X[k] += data[n] * cmplx.Pow(WN, exponent) // I will not remove this command to understand why I used WN_iterator and WNK to optimize operations
            X[k] += data[n] * WNK // this is as equal cmplx.Pow(WN, n*k) because WNk = WN^k and in the next instruction I update WNK that start as WNK^0 to WNK^1 = WNK^0 * WN_iterator WNK^n = WN_iterator^n = WNK^(n-1) * WN_iterator
            WNK *= WN_iterator
        }
        X[k] *= constInverse
        WN_iterator *= WN
    }
    return X
}

func bit_invert(n int, pot int) uintptr {
    qtd_bits := unsafe.Sizeof(n) * 8
    ret := uintptr(0)
    for i := 0; uintptr(i) < qtd_bits && pot != 0; i++ {
        ret <<= 1
        if n % 2 == 1 {
            ret += 1
        }
        n >>= 1
        pot--
    }
    return ret
}

func getDecimatedVector(data []complex128) []complex128 {
    N := len(data)
    X := make([]complex128, N)
    for i := 0; i < N; i++ {
        X[i] = data[bit_invert(i, int(math.Log2(float64(N))))]
    }
    return X
}

func fft_dec_freq(data []complex128, inverse bool) []complex128{

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
    //var WN complex128 = cmplx.Exp((sign * 2i * PI) / N)
    NU := int(math.Log2(float64(arraySize)))
    for L := 1; L < (NU+1); L++ {
        LE := math.Pow(2, float64(NU + 1 - L))
        LE1 := LE / 2
        u := complex(1, 0)
        w := cmplx.Exp((sign * 2i * PI) / complex(LE, 0))
        for J := 1; J < int(LE1) + 1; J++ {
            for I := J; I < arraySize+1; I+=int(LE) {
                IP := I + int(LE1)
                temp := data[I-1] + data[IP-1]
                data[IP-1] = u * (data[I-1] - data[IP-1])
                data[I-1] = temp
            }
            u *= w
        }
    }
    if inverse {
        for i := 0; i < arraySize; i++ {
            data[i] *= constInverse
        }
    }
    return data
}

func fft_dec_time(data []complex128, inverse bool) []complex128{

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
    //var WN complex128 = cmplx.Exp((sign * 2i * PI) / N)
    NU := int(math.Log2(float64(arraySize)))
    for L := 1; L < (NU+1); L++ {
        LE := math.Pow(2, float64(L))
        LE1 := LE / 2
        u := complex(1, 0)
        w := cmplx.Exp((sign * 2i * PI) / complex(LE, 0))
        for J := 1; J < int(LE1) + 1; J++ {
            for I := J; I < arraySize+1; I+=int(LE) {
                IP := I + int(LE1)
                temp := data[IP-1] * u
                data[IP-1] = data[I-1] - temp
                data[I-1] = data[I-1] + temp
            }
            u *= w
        }
    }
    if inverse {
        for i := 0; i < arraySize; i++ {
            data[i] *= constInverse
        }
    }
    return data
}

func fft(data []complex128, inverse bool, frequency_decimation bool) []complex128 {
    var ret []complex128
    qtdZerosAdded := 0
    if verbosity {
        printData(data)
    }
    data, qtdZerosAdded = zeroPadding(data)
    if frequency_decimation {
        ret = getDecimatedVector(fft_dec_freq(data, inverse))
    } else {
        ret = fft_dec_time(getDecimatedVector(data), inverse)
    }
    if padding {
        ret = ret[:len(ret)-qtdZerosAdded]
    }
    if verbosity {
        fmt.Fprintf(os.Stderr, "Quantity of zeros added: %d\n", qtdZerosAdded)
    }
    return ret
}

func zeroPadding(data []complex128) ([]complex128, int) {
    N := len(data)
    lgSize := math.Log2(float64(N))
    if lgSize != float64(int(lgSize)) {
        lgSize = float64(int(lgSize + 1))
    }
    arraySize := math.Pow(2, lgSize)
    for i := N; i < int(arraySize); i++ {
        data = append(data, complex(0, 0))
    }
    return data, (int(arraySize) - N)
}
