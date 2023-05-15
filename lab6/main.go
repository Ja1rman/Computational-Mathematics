package main

import (
	"fmt"
	"log"
	"math"
    "os"
    "github.com/go-echarts/go-echarts/charts"
    "net/http"
)


var INF float64 = 999999999
var FUNC_TYPE = 1
var eps = 0.01
var H = 0.01
var x0 = 0.
var xn = 0.
var y0 = 0.
var methodsDots = map[string][][]float64{}


func InputFromKeyboard() {
    fmt.Println("Выберите функцию:\n" +
              "1 - y' = y + (1 + x)y^2\n" +
              "2 - y' = x^2 - 2y\n" +
              "3 - y' = -y + (x + 1)^3")
    fmt.Scanln(&FUNC_TYPE)
    if FUNC_TYPE < 1 || FUNC_TYPE > 3 {
        log.Fatal("Неизвестная функция")
    }

    fmt.Print("\nВведите начальное условие y0: ")
    fmt.Scanln(&y0)

    fmt.Print("\nВведите интервал дифференцирования x0, xn: ")
    fmt.Scanln(&x0, &xn)

    fmt.Print("\nВведите шаг h: ")
    fmt.Scanln(&H)

    fmt.Print("\nВведите точность ε: ")
    fmt.Scanln(&eps)
}


func f(x float64, y float64) (float64) {
    switch FUNC_TYPE {
    case 1:
        return y + (1+x) * y*y
    case 2:
        return math.Pow(x+1, 3) - y
    case 3:
        return 6*x*x + 5*y
    default:
        log.Fatal("Неизвестная функция")
        return 0.
    }
}


func f2(x float64) (float64) {
    
    switch FUNC_TYPE {
    case 1:
        return -math.Pow(math.E, x) / (x*math.Pow(math.E, x) + (math.E - math.E))
    case 2:
        const_2 := (y0 - x0*x0*x0 - 3*x0 + 2) * math.Pow(math.E, x0)
        return const_2 * math.Pow(math.E, -x) + x*x*x + 3*x - 2
    case 3:
        const_3 := (y0 + 12/125 + (12*x0)/25 + (6*x0*x0)/5) / math.Pow(math.E, 5*x0)
        return const_3 * math.Pow(math.E, 5*x) - (6*x*x)/5 - (12*x)/25 - 12/125
    default:
        log.Fatal("Неизвестная функция")
        return 0.
    }
}


func EulerMethod(a float64, b float64, h float64) ([][]float64) {
    dots := [][]float64{{a, y0}}
    n := int((b - a) / h)
    for ;; {
        y := dots[0][1] + h * f(dots[0][0], dots[0][1])
        y2 := dots[0][1] + h/2 * f(dots[0][0], dots[0][1])

        if math.Abs(y-y2)/(math.Pow(2, 4) - 1) <= eps {
            break
        } 
        h /= 2
    }
    for i := 0; i < n; i++ {
        dots = append(dots, []float64{dots[i][0] + h,
                dots[i][1] + h * f(dots[i][0], dots[i][1])})
    }
    saveToFile(fmt.Sprintf("Шаг метода Эйлера: %f\n", h))
    return dots
}


func RungeKuttaMethod(a float64, b float64, h float64, needPrint bool) ([][]float64) {
    dots := [][]float64{{a, y0}}
    n := int((b - a) / h)
    for ;; {
        k1 := h*f(dots[0][0], dots[0][1])
        k2 := h*f(dots[0][0] + h/2, dots[0][1] + k1/2)
        k3 := h*f(dots[0][0] + h/2, dots[0][1] + k2/2)
        k4 := h*f(dots[0][0] + h, dots[0][1] + k3)
        y := dots[0][1] + (k1 + 2*k2 + 2*k3 + k4)/6
        h /= 2
        k2 = h*f(dots[0][0] + h/2, dots[0][1] + k1/2)
        k3 = h*f(dots[0][0] + h/2, dots[0][1] + k2/2)
        k4 = h*f(dots[0][0] + h, dots[0][1] + k3)
        y2 := dots[0][1] + (k1 + 2*k2 + 2*k3 + k4)/6
        if math.Abs(y-y2)/(math.Pow(2, 4) - 1) <= eps {
            h *= 2
            break
        }
    }
    for i := 0; i < n; i++ {
        k1 := h*f(dots[i][0], dots[i][1])
        k2 := h*f(dots[i][0] + h/2, dots[i][1] + k1/2)
        k3 := h*f(dots[i][0] + h/2, dots[i][1] + k2/2)
        k4 := h*f(dots[i][0] + h, dots[i][1] + k3)
        dots = append(dots, []float64{dots[i][0] + h,
                dots[i][1] + (k1 + 2*k2 + 2*k3 + k4)/6})
    }
    if needPrint {
        saveToFile(fmt.Sprintf("\nШаг метода Рунге-Кутта: %f\n", h))
    }
    return dots
}


func AdamsMethod(a float64, b float64, h float64) ([][]float64) {
    n := int((b - a) / h)
    b1 := math.Min(b, a + 3 * h)
    dots := RungeKuttaMethod(a, b1, h, false)
    epsAdams := 0.
    for i := 3; i < n; i++ {
        df := f(dots[i][0], dots[i][1]) - f(dots[i-1][0], dots[i-1][1])
        d2f := f(dots[i][0], dots[i][1]) - 2 * f(dots[i-1][0], dots[i-1][1]) + 
            f(dots[i-2][0], dots[i-2][1])
        d3f := f(dots[i][0], dots[i][1]) - 3 * f(dots[i-1][0], dots[i-1][1]) + 
            3 * f(dots[i-2][0], dots[i-2][1]) - f(dots[i-3][0], dots[i-3][1])
        dots = append(dots, []float64{dots[i][0] + h,
                      dots[i][1] + h * f(dots[i][0], dots[i][1]) +
                      (h*h) * df / 2 + 5 * (h*h*h) * d2f / 12 + 3 * (h*h*h*h) * d3f / 8})
        epsAdams = math.Max(epsAdams, math.Abs(f2(dots[i+1][0])-dots[i+1][1]))
    }
    saveToFile(fmt.Sprintf("Погрешность метода Адамса: %f\n", epsAdams))
    return dots
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


func httpserver(w http.ResponseWriter, _ *http.Request) {
    xValues := []float64{}
    yValues := []float64{}
    yEulerValues := []float64{}
    yRungeKuttaValues := []float64{}
    yAdamsValues := []float64{}
    for i := 0; i < len(methodsDots["euler"]); i++ {
        xValues = append(xValues, methodsDots["euler"][i][0])
        yValues = append(yValues, methodsDots["true"][i][1])
        yEulerValues = append(yEulerValues, methodsDots["euler"][i][1])
        yRungeKuttaValues = append(yRungeKuttaValues, methodsDots["runge-kutta"][i][1])
        yAdamsValues = append(yAdamsValues, methodsDots["adams"][i][1])
    }
	line := charts.NewLine()
	line.AddXAxis(xValues)
    line.AddYAxis("Точные", yValues, charts.LineOpts{Smooth: true})
    line.AddYAxis("Эйлер", yEulerValues, charts.LineOpts{Smooth: true})
    line.AddYAxis("Рунге-Кутта", yRungeKuttaValues, charts.LineOpts{Smooth: true})
    line.AddYAxis("Адамс", yAdamsValues, charts.LineOpts{Smooth: true})
	line.Render(w)
}


func main() {
    fmt.Println("Лабораторная работа №6, Вариант 27, ЧИСЛЕННОЕ РЕШЕНИЕ ОБЫКНОВЕННЫХ ДИФФЕРЕНЦИАЛЬНЫХ УРАВНЕНИЙ»")
    f, err := os.OpenFile("data/output.txt", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    if _, err = f.WriteString(""); err != nil {
        panic(err)
    }
    // ввод данных
    InputFromKeyboard()

    // вычисления точек
    dots := EulerMethod(x0, xn, H)
    methodsDots["euler"] = dots
    saveToFile(fmt.Sprintf("x:\t\t\t\t\t"))
    for i := 0; i < len(dots); i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][0]))
    }
    saveToFile(fmt.Sprintf("\ny:\t\t\t\t\t"))
    for i := 0; i < len(dots); i++ {
        methodsDots["true"] = dots
        methodsDots["true"][i][1] = f2(dots[i][0])
        saveToFile(fmt.Sprintf("%f\t", methodsDots["true"][i][1]))
    }
    saveToFile(fmt.Sprintf("\nМетод Эйлера:\t\t"))
    for i := 0; i < len(dots); i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    dots = RungeKuttaMethod(x0, xn, H, true)
    methodsDots["runge-kutta"] = dots
    saveToFile(fmt.Sprintf("\nМетод Рунге-Кутта:\t"))
    for i := 0; i < len(dots); i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }
    dots = AdamsMethod(x0, xn, H)
    methodsDots["adams"] = dots
    saveToFile(fmt.Sprintf("\nМетод Адамса:\t\t"))
    for i := 0; i < len(dots); i++ {
        saveToFile(fmt.Sprintf("%f\t", dots[i][1]))
    }

    // график
    http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8080", nil)  
}