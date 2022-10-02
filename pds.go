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
        case "dct":
            dct_flag = true
        case "inv":
            invert_flag = true
        case "freq":
            frequency_decimation_flag = true
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

func fft_prof(data []complex128, inverse bool) []complex128{

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
    return getDecimatedVector(data)
}



func t_fft_prof(data []complex128, inverse bool) []complex128{

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

func fft_calc(data []complex128) []complex128 {
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
        PI := complex(math.Pi, 0)
        /*
        var sign complex128
        constInverse := complex(1, 0)
        if inverse {
            sign = 1
            constInverse /= N
        } else {
            sign = -1
        }
        var WN complex128 = cmplx.Exp((sign * 2i * PI) / N)
        */
        var WN complex128 = cmplx.Exp((-2i * PI) / N)
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
            G := fft_calc(data1)
            if verbosity {
                fmt.Printf("G: %d\n", len(G))
            }
            H := fft_calc(data2)
            if verbosity {
                fmt.Printf("H: %d\n", len(H))
            }
            GSize := len(G)
            for k := 0; k < arraySize; k++ {
                xk := G[k%GSize] + cmplx.Pow(WN, complex(float64(k), 0)*H[k%GSize])
                X = append(X, xk)
            }
        } else {
            X = data
        }
    }
    return X
}

// ifft does the actual work for IFFT
func ifft(data []complex128) []complex128 {
    N := len(data)
    X := data
    // Reverse the input vector
    printData(data)
    for i := 1; i < N/2; i++ {
        j := N - i
        X[i], X[j] = X[j], X[i]
    }

    fmt.Println("Data after change")
    printData(data)
    // Do the transform.
    X = fft_calc(X)
    fmt.Println("Data after FFT")
    printData(X)

    // Scale the output by 1/N
    invN := complex(1.0/float64(N), 0)
    fmt.Printf("InvN: ")
    fmt.Println(invN)
    for i := 0; i < N; i++ {
        fmt.Printf("X[%d]: ", i)
        fmt.Println(X[i])
        X[i] *= invN
    }
    return X
}

func fft(data []complex128, inverse bool, frequency_decimation bool) []complex128 {
    lgSize := math.Log2(float64(len(data)))
    var ret []complex128
    if lgSize == float64(int(lgSize)) {
        if frequency_decimation {
            ret = fft_prof(data, inverse)
        } else {
            ret = t_fft_prof(getDecimatedVector(data), inverse)
        }
    } else {
        os.Stderr.WriteString("Use a vector that is power of 2")
        os.Exit(2)
    }
    return ret
}
