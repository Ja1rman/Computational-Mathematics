package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
    "github.com/go-echarts/go-echarts/charts"
    "net/http"
)


var INF float64 = 999999999
var MIN_POINT float64 = INF
var MAX_POINT float64 = -INF
var POINTS_COUNT int = 500
var DOTS [][]float64
var DELTA float64 = 1


func CheckFrame(dot float64) {
    if dot < MIN_POINT {
        MIN_POINT = dot
    }
    if dot > MAX_POINT {
        MAX_POINT = dot
    }
}


func InputFromKeyboard() () {
    fmt.Print("\nВведите количество точек:")
    n := 0
    fmt.Scanln(&n)
    fmt.Println("\nВведите x,y каждой точки через пробел с новой строки:")
    for i := 0; i < n; i++ {
        x := 0.
        y := 0.
        fmt.Scanln(&x, &y)
        DOTS = append(DOTS, []float64{x,y})
        CheckFrame(x)
    }
    POINTS_COUNT = n
}


func InputFromFile() () {
    file, err := os.Open("./data/input.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        strDots := strings.Split(scanner.Text(), " ")
        dot1, err := strconv.ParseFloat(strDots[0], 64)
        if err != nil { log.Fatal(err) }
        dot2, err := strconv.ParseFloat(strDots[1], 64)
        if err != nil { log.Fatal(err) }
        DOTS = append(DOTS, []float64{dot1, dot2, 0, 0})
        CheckFrame(dot1)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    POINTS_COUNT = len(DOTS)
}


func InputFromFunc() () {
    fmt.Println("\nВведите функцию:\n" +
              "1 - sin(x)\n" +
              "2 - x^2")
    func_type := 1
    fmt.Scanln(&func_type)
    var f = func(x float64) (float64) {return math.Sin(x)}
    switch func_type {
    case 1:
        f = func(x float64) (float64) {return math.Sin(x)}
    case 2:
        f = func(x float64) (float64) {return x*x}
    default:
        log.Fatal("Такой функции не существует")
    }
    fmt.Print("\nВведите исследуемый интервал через пробел:")
    left := 0.
    right := 0.
    fmt.Scanln(&left, &right)
    fmt.Print("\nВведите количество точек:")
    n := 0
    fmt.Scanln(&n)
    for i := 0; i < n; i++ {
        x := left + (right-left)/float64(n-1)*float64(i)
        y := f(x)
        DOTS = append(DOTS, []float64{x,y})
        CheckFrame(x)
    }
    POINTS_COUNT = n
}


func Input() () {
    fmt.Println("Выберите способ ввода данных")
    fmt.Println("1 - ввод таблицы x,y с клавиатуры\n" +
                "2 - ввод таблицы x,y из файла\n" +
                "3 - ввод на основе функции")
    type_of_method := 1
    fmt.Scanln(&type_of_method)
    switch type_of_method { 
    case 1:
        InputFromKeyboard()
    case 2:
        InputFromFile()
    case 3:
        InputFromFunc()
    default:
        log.Fatal("Такого способа ввода не существует")
    }
}


func TableOfDifferences() {
    fmt.Println("\nТаблица разностей:")
    fmt.Print("x\t")
    for i := 0; i < POINTS_COUNT; i++ {
        fmt.Printf("%f\t", DOTS[i][0])
    }
    var prevDelta []float64
    fmt.Print("\ny\t")
    for i := 0; i < POINTS_COUNT; i++ {
        fmt.Printf("%f\t", DOTS[i][1])
        prevDelta = append(prevDelta, DOTS[i][1])
    }
    
    for i := 0; i < POINTS_COUNT-1; i++ {
        fmt.Printf("\nΔ^%dy\t", i+1)
        var newDelta []float64
        for i := 0; i < len(prevDelta)-1; i++ {
            newDelta = append(newDelta, prevDelta[i+1]-prevDelta[i])
            fmt.Printf("%f\t", newDelta[i])
        }
        prevDelta = newDelta
    }
}


func LagrangePolynomial(arg float64) (float64) {
    result := 0.
    for i := 0; i < POINTS_COUNT; i++ {
        numerator := 1.
        denominator := 1.
        for j := 0; j < POINTS_COUNT; j++ {
            if i != j {
                numerator *= arg - DOTS[j][0]
                denominator *= DOTS[i][0] - DOTS[j][0]
            }
        }
        result += DOTS[i][1] * numerator / denominator
    }
    return result
}


// разделённые разности
func f(k int) (float64) {
    k += 1
    result := 0.
    for i := 0; i < k; i++ {
        x := 1.
        for j := 0; j < k; j++ {
            if j != i {
                x *= DOTS[i][0] - DOTS[j][0]
            }
        }
        result += DOTS[i][1]/x
    }
    return result
}


func NewtonPolynomial(arg float64) (float64) {
    p := DOTS[0][1]
    for i := 1; i < POINTS_COUNT; i++ {
        x := 1.
        for j := 0; j < i; j++ {
            x *= arg - DOTS[j][0]
        }
        p += f(i) * x
    }
    return p
}


func httpserver(w http.ResponseWriter, _ *http.Request) {
    xValues := []float64{}
    yLagrangeValues := []float64{}
    yNewtonValues := []float64{}
    for i := MIN_POINT-DELTA; i < MAX_POINT+DELTA; i += 0.001 {
        xValues = append(xValues, i)
        yLagrangeValues = append(yLagrangeValues, LagrangePolynomial(i))
        yNewtonValues = append(yNewtonValues, NewtonPolynomial(i))
    }
	line := charts.NewLine()
	line.AddXAxis(xValues)
    line.AddYAxis("Лагранж", yLagrangeValues, charts.LineOpts{Smooth: true})
    line.AddYAxis("Ньютон", yNewtonValues, charts.LineOpts{Smooth: true})
	line.Render(w)

    xValues = []float64{}
    yValues := []float64{}
    for i := 0; i < POINTS_COUNT; i++ {
        xValues = append(xValues, DOTS[i][0])
        yValues = append(yValues, DOTS[i][1])
    }
    line2 := charts.NewLine()
    line2.AddXAxis(xValues)
    line2.AddYAxis("Исходные", yValues)
    line2.Render(w)
}


func main() {
    fmt.Println("Лабораторная работа №5, Вариант 27, Интерполяция функций")
    // ввод данных
    Input()
    DELTA = (MAX_POINT - MIN_POINT)/float64(POINTS_COUNT)

    // таблица разностей
    TableOfDifferences()
    
    // ввод аргумента
    fmt.Println("\nВведите значение аргумента для интерполирования")
    arg := 0.
    fmt.Scanln(&arg)
    
    fmt.Printf("По Лагранжу: %.9f\n", LagrangePolynomial(arg))
    fmt.Printf("По Ньютону: %.9f\n", NewtonPolynomial(arg))

    // график
    http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8080", nil)  
}
