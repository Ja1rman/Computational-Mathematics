package main

import (
    "fmt"
    "math"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/plotutil"
)

const (
    tol   = 0.01
    maxIt = 100
)

type Func func(x []float64) []float64
type Jacobian func(x []float64, jac [][]float64) [][]float64

func Newton(f Func, J Jacobian, x []float64) ([]float64, int, error) {
    n := len(x)

    for it := 0; it < maxIt; it++ {
        fx := f(x)
        fmt.Println("Вектор погрешностей:", fx)
        // Проверить сходимость
        rms := 0.0
        for i := range fx {
            rms += fx[i] * fx[i]
        }
        rms = math.Sqrt(rms / float64(n))
        if rms < tol {
            return x, it, nil
        }

        // Вычисление ябиана
        jac := make([][]float64, n)
        for i := range jac {
            jac[i] = make([]float64, n)
        }
        J(x, jac)

        // Решение линейной системы
        dx, err := SolveLinearSystem(jac, fx)
        if err != nil {
            return nil, it, err
        }

        // Обновление x
        for j := range x {
            x[j] -= dx[j]
        }
    }

    return nil, 0, fmt.Errorf("Newton method failed to converge")
}

func SolveLinearSystem(A [][]float64, b []float64) ([]float64, error) {
    n := len(A)
    if len(b) != n {
        return nil, fmt.Errorf("invalid dimensions")
    }

    // прямой ход Гаусса
    for k := 0; k < n; k++ {
        // Найти точку опоры
        maxRow := k
        for i := k + 1; i < n; i++ {
            if math.Abs(A[i][k]) > math.Abs(A[maxRow][k]) {
                maxRow = i
            }
        }

        // Поменять местами строки
        A[k], A[maxRow] = A[maxRow], A[k]
        b[k], b[maxRow] = b[maxRow], b[k]

        for i := k + 1; i < n; i++ {
            factor := A[i][k] / A[k][k]
            for j := k + 1; j < n; j++ {
                A[i][j] -= factor * A[k][j]
            }
            b[i] -= factor * b[k]
        }
    }

    // Обратный ход
    x := make([]float64, n)
    for i := n - 1; i >= 0; i-- {
        sum := b[i]
        for j := i + 1; j < n; j++ {
            sum -= A[i][j] * x[j]
        }
        if math.Abs(A[i][i]) < tol {
            return nil, fmt.Errorf("matrix is singular")
        }
        x[i] = sum / A[i][i]
    }

    return x, nil
}

func drawPlot(type_func int) {
    // Создаем новый график
    p := plot.New()
    p.Title.Text = "График функций"
    p.X.Label.Text = "X"
    p.Y.Label.Text = "Y"

    
    var f1, f2, g1 *plotter.Function

    if type_func == 1 {
        f1 = plotter.NewFunction(func(x float64) float64 { 
            if 4 - x*x > 0 {
                return math.Sqrt(math.Abs(4 - x*x)) 
            }else {
                return 0
            }})
        g1 = plotter.NewFunction(func(x float64) float64 {
            if 4 - x*x > 0 {
                return -math.Sqrt(4 - x*x)
            }else {
                return 0
            }
        })
        f2 = plotter.NewFunction(func(x float64) float64 { return 3 * x * x })
    }else {
        f1 = plotter.NewFunction(func(x float64) float64 { 
            if 2*x*x-1 > 0 {
                return math.Sqrt(2*x*x - 1) 
            }else {
                return 0
            }
        })
        g1 = plotter.NewFunction(func(x float64) float64 {
            if 2*x*x-1 > 0 {
                return -math.Sqrt(2*x*x-1)
            }else {
                return 0
            }
        })
        f2 = plotter.NewFunction(func(x float64) float64 { return math.Log(math.Abs(x+1)) })
    }
    // Создаем данные для графика первой функции
    f1.Samples = 10000 // количество точек на графике
    f1.Color = plotutil.Color(0) // цвет графика

    // Создаем данные для графика второй функции 
    f2.Samples = 10000 // количество точек на графике
    f2.Color = plotutil.Color(1) // цвет графика

    g1.Samples = 10000 // количество точек на графике
    g1.Color = plotutil.Color(0) // цвет графика

    // Добавляем данные на график
    p.Add(f1, f2, g1)

    // Задаем промежуток по оси X
    p.X.Min = -10
    p.X.Max = 10
    p.Y.Min = -10
    p.Y.Max = 10

    // Сохраняем график в файл
    if err := p.Save(10*vg.Inch, 10*vg.Inch, "data/plot.png"); err != nil {
        panic(err)
    }
}

func main() {
    var f Func
    var J Jacobian
    var system_type int
    fmt.Println("Выберите номер системы линейных уравнений:")
    fmt.Println("1 - x^2 + y^2 = 4; y = 3*x^2\n" +
                "2 - 2*x^2 - y^2 = 1; x - e^(y) = -1")
    fmt.Scanln(&system_type)
    switch system_type {
    case 1:
        f = func(x []float64) []float64 {
            y := make([]float64, 2)
            y[0] = x[0]*x[0] + x[1]*x[1] - 4
            y[1] = x[1] -3*x[0]*x[0]
            return y
        }
        J = func(x []float64, jac [][]float64) [][]float64 {
            jac[0][0] = 2 * x[0]
            jac[0][1] = 2 * x[1]
            jac[1][0] = -6*x[0]
            jac[1][1] = 1
            return jac
        }
    case 2:
        f = func(x []float64) []float64 {
            y := make([]float64, 2)
            y[0] = 2*x[0]*x[0] - x[1]*x[1] - 1
            y[1] = x[0] - math.Exp(x[1]) + 1
            return y
        }
        J = func(x []float64, jac [][]float64)[][]float64 {
            jac[0][0] = 4 * x[0]
            jac[0][1] = -2 * x[1]
            jac[1][0] = 1
            jac[1][1] = -math.Exp(x[1])
            return jac
        }
    default:
        panic("Неизвестная система")
    }

    x0 := []float64{0.5, 0.5}
    fmt.Println("Введите начальное приближение:")
    fmt.Scanln(&x0[0], &x0[1])
    drawPlot(system_type)
    x, n, err := Newton(f, J, x0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Solution:", x, " Итераций:", n)
    }
}
