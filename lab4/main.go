package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)


var x1 float64 = 0
var x2 float64 = 0
var x3 float64 = 0
var x4 float64 = 0
var x5 float64 = 0
var x6 float64 = 0
var y float64 = 0
var xy float64 = 0
var x2y float64 = 0
var x3y float64 = 0
var lnX float64 = 0
var ln2X float64 = 0
var lnY float64 = 0
var lnXlnY float64 = 0
var xLnY float64 = 0
var yLnX float64 = 0

var INF float64 = 9999999999
var POINTS_COUNT int = 500
var MIN_POINT float64 = INF
var MAX_POINT float64 = -INF
var MIN_EPS float64 = INF
var BEST_FUNC string
var FUNC = map[string][]float64{}
var FUNC_NAME = map[string]string{}


func InputFromFile() ([][]float64) {
    var arr [][]float64

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
        arr = append(arr, []float64{dot1, dot2, 0, 0})
        if dot1 < MIN_POINT {
            MIN_POINT = dot1
        }
        if dot1 > MAX_POINT {
            MAX_POINT = dot1
        }
    }
    MIN_POINT -= 2
    MAX_POINT += 2

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    //if len(arr) < 8 || len(arr) > 12 {
    //    log.Panicln("Точек в файле должно быть от 8 до 12")
    //}
    return arr
}


func SumAll(dots [][]float64) {
    n := len(dots)
    for i := 0; i < n; i++ {
        x1 += dots[i][0]
        x2 += dots[i][0]*dots[i][0]
        x3 += dots[i][0]*dots[i][0]*dots[i][0]
        x4 += dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]
        x5 += dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]
        x6 += dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]*dots[i][0]
        y += dots[i][1]
        xy += dots[i][0]*dots[i][1]
        x2y += dots[i][0]*dots[i][0]*dots[i][1]
        x3y += dots[i][0]*dots[i][0]*dots[i][0]*dots[i][1]
    
        lnX += math.Log(dots[i][0])
        ln2X += math.Pow(math.Log(dots[i][0]), 2)
        lnY += math.Log(dots[i][1])
        lnXlnY += math.Log(dots[i][0]) * math.Log(dots[i][1])
        xLnY += dots[i][0] * math.Log(dots[i][1])
        yLnX += dots[i][1] * math.Log(dots[i][0])
    }
}


func saveToFile(data string) {
    f, err := os.OpenFile("data/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString(data); err != nil {
        panic(err)
    }
}


func LinearApproximation(dots [][]float64) () {
    n := len(dots)
    d := x2*float64(n) - x1*x1
    d1 := xy*float64(n) - x1*y
    d2 := x2*y - x1*xy
    a := d1/d
    b := d2/d
    S := 0. // мера отклонения
    r1 := 0. // числитель коэффициента корреляции
    r2 := 0. // знаменатель x коэффициента корреляции
    r3 := 0. // знаменатель y коэффициента корреляции
    for i := 0; i < n; i++ {
        dots[i][2] = a*dots[i][0]+b // значение аппроксимирующей функции ax+b
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
        r1 += (dots[i][0]-x1/float64(n)) * (dots[i][1]-y/float64(n))
        r2 += (dots[i][0]-x1/float64(n)) * (dots[i][0]-x1/float64(n))
        r3 += (dots[i][1]-y/float64(n)) * (dots[i][1]-y/float64(n))
    }
    r := r1 / math.Sqrt(r2*r3) // коэффициент корреляции
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Линейная"
    }
    FUNC["lin"] = []float64{a, b}
    FUNC_NAME["lin"] = fmt.Sprintf("%fx + %f\n", a, b)
    saveToFile("Линейная аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %fx + %f\n", a, b))
    saveToFile(fmt.Sprintf("S = %f, r = %f, Среднеквадратичное отклонение = %f\n", S, r, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


func GaussSolver(matrix [][]float64, lines int) ([]float64) {
	columns := len(matrix[0]) - 1
	a := make([][]float64, lines)
	for i := range a {
		a[i] = make([]float64, columns)
	}

	for i := 0; i < lines; i++ {
		for j := 0; j < columns; j++ {
			a[i][j] = matrix[i][j]
		}
	}

	b := make([]float64, lines)
	for i := range b {
		b[i] = matrix[i][columns]
	}

	index := make([]int, len(a))
	for i := range index {
		index[i] = i
	}
    
	for i := 0; i < len(a); i++ {

		r := a[i][index[i]]
		if r == 0 {
			var kk int

			for k := i; k < len(a); k++ {
				if math.Abs(a[i][index[k]]) > r {
					kk = k
				}
			}

			if kk > 0 {
				index[i], index[kk] = index[kk], index[i]
			}
			r = a[i][index[i]]
		}

		if r == 0 {
			if b[i] == 0 {
				panic("Система имеет множество решений")
			} else {
				panic("Система не имеет решений")
			}
		}

		for j := 0; j < len(a[i]); j++ {
			a[i][index[j]] /= r
		}
		b[i] /= r

		for k := i + 1; k < len(a); k++ {
			r = a[k][index[i]]
			for j := 0; j < len(a[i]); j++ {
				a[k][index[j]] = a[k][index[j]] - a[i][index[j]]*r
			}
			b[k] = b[k] - b[i]*r
		}
	}

	var x []float64 = make([]float64, len(b))

	for i := len(a) - 1; i >= 0; i-- {
		x[i] = b[i]

		for j := i + 1; j < len(a); j++ {
			x[i] = x[i] - (x[j] * a[i][index[j]])
		}
	}

	return x
}


func QuadraticApproximation(dots [][]float64) () {
    n := len(dots)
    
    matrix := [][]float64{{float64(n), x1, x2, y}, 
                    {x1, x2, x3, xy},
                    {x2, x3, x4, x2y}}
    a := GaussSolver(matrix, 3)
    S := 0.
    for i := 0; i < n; i++ {
        dots[i][2] = a[0] + a[1]*dots[i][0] + a[2]*dots[i][0]*dots[i][0] // значение аппроксимирующей функции a0 + a1x + a2x^2
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
    }
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Квадратичная"
    }
    FUNC["quad"] = []float64{a[0], a[1], a[2]}
    FUNC_NAME["quad"] = fmt.Sprintf("%f + %fx + %fx^2", a[0], a[1], a[2])
    saveToFile("\n\nКвадратичная аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %f + %fx + %fx^2\n", a[0], a[1], a[2]))
    saveToFile(fmt.Sprintf("S = %f, Среднеквадратичное отклонение = %f\n", S, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


func CubicApproximation(dots [][]float64) () {
    n := len(dots)
    
    matrix := [][]float64{{float64(n), x1, x2, x3, y}, 
                        {x1, x2, x3, x4, xy},
                        {x2, x3, x4, x5, x2y},
                        {x3, x4, x5, x6, x3y}}
    a := GaussSolver(matrix, 4)
    S := 0.
    for i := 0; i < n; i++ {
        dots[i][2] = a[0] + a[1]*dots[i][0] + a[2]*dots[i][0]*dots[i][0] +
            a[3]*dots[i][0]*dots[i][0]*dots[i][0] // значение аппроксимирующей функции a0 + a1x + a2x^2 + a3x^3
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
    }
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Кубическая"
    }
    FUNC["cub"] = []float64{a[0], a[1], a[2], a[3]}
    FUNC_NAME["cub"] = fmt.Sprintf("%f + %fx + %fx^2 + %fx^3", a[0], a[1], a[2], a[3])
    saveToFile("\n\nКубическая аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %f + %fx + %fx^2 + %fx^3\n", a[0], a[1], a[2], a[3]))
    saveToFile(fmt.Sprintf("S = %f, Среднеквадратичное отклонение = %f\n", S, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


func ExponentApproximation(dots [][]float64) {
    n := len(dots)
    if (math.IsNaN(lnY)) {
        saveToFile("\n\nЭкспоненциальная аппроксимация невозможна из-за отрицательного числа Y\n")
        return
    }
    matrix := [][]float64{{x2, x1, xLnY},
                          {x1, float64(n), lnY}}
    a := GaussSolver(matrix, 2)
    a[1] = math.Exp(a[1])
    S := 0.
    for i := 0; i < n; i++ {
        dots[i][2] = a[1] * math.Exp(dots[i][0] * a[0]) // значение аппроксимирующей функции a * e^bx
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
    }
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Экспоненциальная"
    }
    FUNC["exp"] = []float64{a[1], a[0]}
    FUNC_NAME["exp"] = fmt.Sprintf("%f * e^%fx", a[1], a[0])
    saveToFile("\n\nЭкспоненциальная аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %f * e^%fx\n", a[1], a[0]))
    saveToFile(fmt.Sprintf("S = %f, Среднеквадратичное отклонение = %f\n", S, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


func LogApproximation(dots [][]float64) {
    n := len(dots)
    if (math.IsNaN(lnX)) {
        saveToFile("\n\nЛогарифмическая аппроксимация невозможна из-за отрицательного числа X\n")
        return
    }
    matrix := [][]float64{{ln2X, lnX, yLnX},
                          {lnX, float64(n), y}}
    a := GaussSolver(matrix, 2)
    S := 0.
    for i := 0; i < n; i++ {
        dots[i][2] = a[0] * math.Log(dots[i][0]) + a[1] // значение аппроксимирующей функции a * lnx + b
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
    }
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Логарифмическая"
    }
    FUNC["log"] = []float64{a[1], a[0]}
    FUNC_NAME["log"] = fmt.Sprintf("%f * lnx + %f", a[1], a[0])
    saveToFile("\n\nЛогарифмическая аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %f * lnx + %f\n", a[1], a[0]))
    saveToFile(fmt.Sprintf("S = %f, Среднеквадратичное отклонение = %f\n", S, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


func PowApproximation(dots [][]float64) {
    n := len(dots)
    if (math.IsNaN(lnX)) {
        saveToFile("\n\nСтепенная аппроксимация невозможна из-за отрицательного числа X\n")
        return
    }
    if (math.IsNaN(lnY)) {
        saveToFile("\n\nСтепенная аппроксимация невозможна из-за отрицательного числа Y\n")
        return
    }
    matrix := [][]float64{{ln2X, lnX, lnXlnY},
                          {lnX, float64(n), lnY}}
    a := GaussSolver(matrix, 2)
    a[1] = math.Exp(a[1])
    S := 0.
    for i := 0; i < n; i++ {
        dots[i][2] = a[1] * math.Pow(dots[i][0], a[0]) // значение аппроксимирующей функции a * x^b
        dots[i][3] = dots[i][2]-dots[i][1] // эпсилоен(отклонение)
        S += (dots[i][2]-dots[i][1])*(dots[i][2]-dots[i][1])
    }
    eps := math.Sqrt(S/float64(n))
    if eps < MIN_EPS {
        MIN_EPS = eps
        BEST_FUNC = "Степенная"
    }
    FUNC["pow"] = []float64{a[1], a[0]}
    FUNC_NAME["pow"] = fmt.Sprintf("%f * x^%f", a[1], a[0])
    saveToFile("\n\nСтепенная аппроксимация:\n")
    saveToFile(fmt.Sprintf("P(x) = %f * x^%f\n", a[1], a[0]))
    saveToFile(fmt.Sprintf("S = %f, Среднеквадратичное отклонение = %f\n", S, eps))
    saveToFile("Значения:\n")
    saveToFile("x\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile("\ny\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    saveToFile("\nP(x)\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][2]))
    }
    saveToFile("\ne\t")
    for i := 0; i < n; i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][3]))
    }
}


// generate random data for line chart
func generateLineItems(name string, args []float64) []opts.LineData {
	var items []opts.LineData
	for i := MIN_POINT; i < MAX_POINT; i += (MAX_POINT-MIN_POINT) / float64(POINTS_COUNT) {
        val := 0.
        switch name {
            case "lin": val = args[0]*i + args[1]
            case "quad": val = args[0] + args[1]*i + args[2]*i*i
            case "cub": val = args[0] + args[1]*i + args[2]*i*i + args[3]*i*i*i
            case "exp": val = args[0] * math.Exp(args[1]*i)
            case "log": 
                if i > 0 {
                    val = args[0] * math.Log(i) + args[1]
                    items = append(items, opts.LineData{Value: val})
                } 
            case "pow": 
            val = args[0] * math.Pow(i, args[1])
            if !math.IsNaN(val) {
                items = append(items, opts.LineData{Value: val})
            }
            default: val = 0
        }
		if name != "log" && name != "pow" {
            items = append(items, opts.LineData{Value: val})
        }
	}
	return items
}


func httpserver(w http.ResponseWriter, _ *http.Request) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}))
    ab := []float64{}
    for i := MIN_POINT; i < MAX_POINT; i += (MAX_POINT-MIN_POINT) / float64(POINTS_COUNT) {
        ab = append(ab, i)
    }
	line.SetXAxis(ab).
         AddSeries(FUNC_NAME["lin"], generateLineItems("lin", FUNC["lin"])).
		 AddSeries(FUNC_NAME["quad"], generateLineItems("quad", FUNC["quad"])).
         AddSeries(FUNC_NAME["cub"], generateLineItems("cub", FUNC["cub"]))
    if _, ok := FUNC_NAME["exp"]; ok{
        line.AddSeries(FUNC_NAME["exp"], generateLineItems("exp", FUNC["exp"])) 
    } 
    if _, ok := FUNC_NAME["log"]; ok{
        line.AddSeries(FUNC_NAME["log"], generateLineItems("log", FUNC["log"])) 
    }
    if _, ok := FUNC_NAME["pow"]; ok{
        line.AddSeries(FUNC_NAME["pow"], generateLineItems("pow", FUNC["pow"])) 
    }
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)
}


func main() {
    fmt.Println("Лабораторная работа №4, Вариант 27, АППРОКСИМАЦИЯ ФУНКЦИИ МЕТОДОМ НАИМЕНЬШИХ КВАДРАТОВ")
    f, err := os.OpenFile("data/output.txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err = f.WriteString(""); err != nil {
        panic(err)
    }
    dots := InputFromFile()
    SumAll(dots)
    LinearApproximation(dots)
    QuadraticApproximation(dots)
    CubicApproximation(dots)
    ExponentApproximation(dots)
    LogApproximation(dots)
    PowApproximation(dots)
    saveToFile("\n\nЛучшая аппроксимация: " + BEST_FUNC + "\n")

    http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8080", nil)       
}
