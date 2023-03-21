package main

import (
    "fmt"
    "math"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "os/exec"
)

const (
    eps   = 0.01
    maxIt = 1000
)

type Func func(x []float64) []float64
type Jacobian func(x []float64, jac [][]float64) [][]float64


func Newton(f Func, J Jacobian, x []float64) ([]float64, int, error) {
    n := len(x)
    for it := 0; it < maxIt; it++ {
        // Вычисление ябиана
        jac := make([][]float64, n)
        for i := range jac {
            jac[i] = make([]float64, n)
        }
        J(x, jac)

        // Решение линейной системы
        fx := f(x)
        dx, err := Solve(jac, fx)
        if err != "ok" {
            return nil, it, fmt.Errorf(err)
        }
        totalDX := dx[0]
        for j := range dx {
            if math.Abs(dx[j]) > totalDX {
                totalDX = math.Abs(dx[j])
            }
        }
        // Обновление x
        for j := range x {
            x[j] -= dx[j]
        }
        fx = f(x)
        
        for j := range x {
            fmt.Print(x[j], " ")
        }
        fmt.Println()
        if totalDX <= eps {
            return x, it+1, nil
        }
    }

    return nil, 0, fmt.Errorf("Метод Ньютона не сходится за %d итераций", maxIt)
}


func Solve(matrix [][]float64, answers []float64) ([]float64, string) {
    n := len(matrix)    
    index := make([]int, n)
	for i := range index {
		index[i] = i
	}
    per := 0.
    // прямой ход
    for i := 0; i < n; i++ {
        // главный элемент  по умолчанию
        main_elem := matrix[i][index[i]]
        // если главный элемент равен нулю, то нужно найти другой методом перестановки колонок в матрице
        if main_elem == 0 {
            var k int
            for j := i; j < n; j++ {
                if matrix[i][index[j]] != .0 {
                    k = j
                    per += 1
                    break
                }
            }
            if k > 0 {
                // если удалось найти главный элемент, меняем местами колонки, так чтобы главный элемент встал в диагональ матрицы
                index[i], index[k] = index[k], index[i]
            }
            // главный элемента текущей строки из диагонали
            main_elem = matrix[i][index[i]];
        }
        // если главный элемент строки всё ещё равен 0, то метод гаусса не работает (можно не проверять, так как считали определитель)
        if main_elem == 0 {
            if answers[i] == 0 {
                return answers, "Система имеет множество решений"
            } else {
                return answers, "Система не имеет решений"
            }
        }

        // вычитание текущей строки из всех ниже расположенных строк с занулением i-ого элемента в каждой из них
        for j := i + 1; j < n; j++ {
            main_elem = matrix[j][index[i]]
            p := main_elem/matrix[i][index[i]]
            for m := 0; m < n; m++ {
                matrix[j][index[m]] -= matrix[i][index[m]] * p
            }
            answers[j] -= answers[i]*p
        }
    }
    result := make([]float64, len(answers))
    for i := 0; i < n; i++ {
        main_elem := matrix[i][index[i]]
        for j := 0; j < n; j++ {
            matrix[i][index[j]] /= main_elem
        }
        answers[i] /= main_elem
    }
    // обратный ход
	for i := n - 1; i >= 0; i-- {
		// начальное значение элемента x[i]
		result[i] = answers[i]
		for j := i + 1; j < n; j++ {
			result[i] -= result[j] * matrix[i][index[j]];
		}
	}
    return result, "ok"
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
                return math.Sqrt(4 - x*x)
            } else {
                return 0
            } 
        })
        g1 = plotter.NewFunction(func(x float64) float64 {
            if 4 - x*x > 0 {
                return -math.Sqrt(4 - x*x)
            } else {
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

        f2 = plotter.NewFunction(func(x float64) float64 { return math.Log(x+1) })
        f2.XMin = -0.999999999999
        f2.XMax = 100
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
    cmd := exec.Command("open", "data/plot.png")
    if err := cmd.Run(); err != nil {
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
    drawPlot(system_type)
    x0 := []float64{0.5, 0.5}
    fmt.Println("Введите начальное приближение:")
    fmt.Scanln(&x0[0], &x0[1])
    
    x, n, err := Newton(f, J, x0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Решение:", x, " Итераций:", n)
        fmt.Printf("Проверка:\nf(x,y)) = %v;\n" +
        "g(x,y) = %v\n", f(x)[0], f(x)[1])
    }
}
