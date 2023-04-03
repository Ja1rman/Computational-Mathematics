package main

import (
    "fmt"
    "math"
    "log"
)

var type_of_func float64 = 1
var START_N int = 4
var MAX_ITERS int = 100000
var INF float64 = 999999999

func InputFromKeyboard() (float64, float64, float64, int) {
    fmt.Println("\nВыберите номер функции")
    fmt.Print("1 - 2x^3 - 3x^2 - 5x + 27\n" +
              "2 - x^2\n" +
              "3 - 1/x\n")
    fmt.Scanln(&type_of_func)
    if type_of_func > 3 || type_of_func < 1 {
        log.Panicln("Выбрана неверная функция")
    }

    var a, b, eps float64
    fmt.Print("Введите границы интервала через пробел: ")
    fmt.Scanln(&a, &b)
    fmt.Print("Введите погрешность вычислений: ")
    fmt.Scanln(&eps)

    fmt.Println("\nВыберите номер метода решения")
    fmt.Print("1 - Метод прямоугольников\n" +
              "2 - Метод трапеций\n" +
              "3 - Метод Симпсона\n")
    type_of_method := 1
    fmt.Scanln(&type_of_method)
    if type_of_method > 3 || type_of_method < 1 {
        log.Panicln("Выбран неверный метод")
    }

    return a, b, eps, type_of_method
}


func rectangleMethod(a float64, b float64, eps float64, k int) (float64, int) {
    // средние
    n := START_N
    res := INF
    for n <= MAX_ITERS {
        temp := 0.
        x := a
        h := (b-a) / float64(n)
        for i := 0; i < n; i++ {
            if k == 0 {
                temp += f(x)
            } else if k == 1 {
                temp += f(x + h)
            } else {
                temp += f(x + h/2)
            }
            x += h
        }
        temp *= h

        nowEps := math.Abs(res - temp)
        res = temp
        if nowEps / (math.Pow(2, 2)-1) <= eps {
            break
        } else {
            n *= 2
        }
    }
    if n > MAX_ITERS {
        log.Panicln("Превышено число разбиений")
    }
    return res, n
}


func trapezoidalMethod(a float64, b float64, eps float64) (float64, int) {
    n := START_N
    res := 0.
    for n <= MAX_ITERS {
        temp := (f(a) + f(b)) / 2
        h := (b-a) / float64(n)
        x := a + h
        for i := 0; i < n-1; i++ {
            temp += f(x)
            x += h
        }
        temp *= h

        nowEps := math.Abs(res - temp)
        res = temp
        if nowEps / (math.Pow(2, 2)-1) <= eps {
            break
        } else {
            n *= 2
        }
    }
    if n > MAX_ITERS {
        log.Panicln("Превышено число разбиений")
    }
    return res, n
}


func simpsonMethod(a float64, b float64, eps float64) (float64, int) {
    if START_N % 2 != 0 {
        log.Panicln("Установлено нечётное число разбиений!")
    }
    n := START_N
    res := 0.
    for n <= MAX_ITERS {
        temp := f(a) + f(b)
        h := (b-a) / float64(n)
        x := a + h
        for i := 0; i < n-1; i++ {
            if i % 2 == 0 {
                temp += 4 * f(x)
            } else {
                temp += 2 * f(x)
            }
            x += h
        }
        temp *= h/3

        nowEps := math.Abs(res - temp)
        res = temp
        if nowEps / (math.Pow(2, 4)-1) <= eps {
            break
        } else {
            n *= 2
        }
    }
    if n > MAX_ITERS {
        log.Panicln("Превышено число разбиений")
    }
    return res, n
}


func f(x float64) float64 {
    switch type_of_func {
    case 1:
        return 2*math.Pow(x, 3) - 3*math.Pow(x, 2) - 5*x + 27
    case 2:
        return x*x
    case 3:
        return 1/x
    default:
        panic("System hacked")
    }
}


func verifyInputs(a float64, b float64, eps float64) bool {
    // Проверяем, что a < b
    if a >= b {
        fmt.Println("Error: Правая граница должна быть больше левой.")
        return false
    }

    // Проверяем, что точность eps положительна
    if eps <= 0 {
        fmt.Println("Error: эпсилон должно быть положительным.")
        return false
    }

    // Проверяем, что определённый интеграл существует на првмежутке
    return isContinuous(a, b)
}

// функция для проверки непрерывности функции на интервале [a, b]
func isContinuous(a, b float64) bool {
    for x := a; x <= b; x += 0.1 {
        lim := limit(x, 1)
        if lim != math.NaN() && (lim + 0.001 < f(x) || lim - 0.001 > f(x))  {
            fmt.Println(x, lim, f(x))
            return false
        }
    }
    return true
}

// функция для вычисления предела функции в точке x
func limit(x, p float64) float64 {
    const delta = 0.0001
    h := delta * math.Pow(10, p)

    l_y := f(x - h)
    r_y := f(x + h)

    if math.IsNaN(l_y) || math.IsNaN(r_y) {
        return math.NaN()
    }

    return (l_y + r_y) / 2
}


func main() {
    fmt.Println("Лабораторная работа №3, Вариант 27, Численное интегрирование")
    // ввод данных
    a, b, eps, type_of_method := InputFromKeyboard()

    // проверка данных
    if !verifyInputs(a, b, eps) {
        panic("Проверяте данные!")
    }

    // решаем уравнение
    var x float64
    var n int
    switch type_of_method { 
    case 1:
        x, n = rectangleMethod(a, b, eps, 1)
        fmt.Println("Значение интеграла методом левых прямоугольников: ", x, "Число разбиений =", n)
        x, n = rectangleMethod(a, b, eps, 2)
        fmt.Println("Значение интеграла методом правых прямоугольников: ", x, "Число разбиений =", n)
        x, n = rectangleMethod(a, b, eps, 3)
        fmt.Println("Значение интеграла методом средних прямоугольников: ", x, "Число разбиений =", n)
    case 2:
        x, n = trapezoidalMethod(a, b, eps)
        fmt.Println("Значение интеграла методом трапеций: ", x, "Число разбиений =", n)
    case 3:
        x, n = simpsonMethod(a, b, eps)
        fmt.Println("Значение интеграла методом симпсона: ", x, "Число разбиений =", n)
    default:
        log.Fatal("Метода не существует")
    }
    
}
