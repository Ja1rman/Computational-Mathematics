package main

import (
    "fmt"
    "math"
    "os"
    "io/ioutil"
    "bufio"
    "log"
    "strconv"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/gonum/diff/fd"
    "os/exec"
)

var type_of_func float64 = 1

func InputFromFile() (float64, float64, float64) {
    var arr []float64

    file, err := os.Open("data/input.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        inp, err := strconv.ParseFloat(scanner.Text(), 64)
        if err != nil {
            log.Fatal(err)
        }
        arr = append(arr, inp)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    if len(arr) != 3 {
        log.Panicln("Неверное количество аргументов в файле")
    }
    return arr[0], arr[1], arr[2]
}


func InputFromKeyboard() (float64, float64, float64) {
    var a, b, eps float64
    fmt.Print("Введите границы интервала через пробел: ")
    fmt.Scanln(&a, &b)
    fmt.Print("Введите погрешность вычислений: ")
    fmt.Scanln(&eps)
    return a, b, eps
}


func f(x float64) float64 {
    switch type_of_func {
    case 1:
        return 2.335*math.Pow(x, 3) + 3.98*math.Pow(x,2) - 4.52*x - 3.11
    case 2:
        return math.Pow(x, 3) - x + 4
    case 3:
        return math.Sin(x) + 0.1
    default:
        panic("System hacked")
    }
}

// Функция решения уравнения методом хорд
func chordMethod(a float64, b float64, eps float64) (float64, float64, int, float64) {
    x_prev := 9999999999.
    iters := 0
    x := a - (a*f(a) - b*f(a)) / (f(b) - f(a))
    for math.Abs(f(x)) > eps || math.Abs(x - x_prev) > eps {
        x_prev = x
        x = (a*f(b) - b*f(a)) / (f(b) - f(a))
        
        if f(x)*f(a) < 0 {
            b = x
        } else {
            a = x
        }
        iters++
    }

    return x, f(x), iters, x-x_prev
}

// Функция решения уравнения методом секущих
func secantMethod(a float64, b float64, eps float64) (float64, float64, int) {
    iters := 0
    x0 := 0.
    if f(a) * fd.Derivative(f, (fd.Derivative(f, a, nil)), nil) > 0 {
        x0 = a
    }else {
        x0 = b
    }
    var x2 float64
    x1 := x0 + 2*eps
    fmt.Printf("x0 = %v, x1 = %v\n", x0, x1)
    for math.Abs(x1-x0) > eps || math.Abs(f(x1)) > eps {
        x2 = x1 - f(x1)*(x1-x0)/(f(x1)-f(x0))
        x0 = x1
        x1 = x2
        iters++
    }
    return x2, f(x2), iters
}

// Функция решения уравнения методом простой итерации
func simpleIteration(a float64, b float64, eps float64) (float64, float64, int) {
    const maxIterations = 1000
    var x0 float64
    if fd.Derivative(f, a, nil) > fd.Derivative(f, b, nil){
        x0 = a
    }else {
        x0 = b 
    }
    lambda := math.Max(fd.Derivative(f, a, nil), fd.Derivative(f, b, nil))
    x := x0
    fmt.Printf("fi'(a) = %v\n", 1 - fd.Derivative(f, a, nil)/lambda)
    fmt.Printf("fi'(b) = %v\n", 1 - fd.Derivative(f, b, nil)/lambda)
    for i := 0; i < maxIterations; i++ {
        // Вычисляем следующее приближение
        xNext := x - f(x)/lambda
        // Проверяем достижение необходимой точности
        if math.Abs(xNext-x) <= eps && math.Abs(f(xNext)) <= eps {
            return xNext, f(xNext), i+1
        }
        x = xNext
    }
    // Если не достигли нужной точности за максимальное число итераций, выдаем ошибку
    panic(fmt.Sprintf("Метод простых итераций не смог дать точное решение за %d итераций", maxIterations))
}

func verifyInputs(a float64, b float64, eps float64) bool {
    // Проверяем, что a < b
    if a >= b {
        fmt.Println("Error: a должно быть больше b.")
        return false
    }
    // Проверяем, что функция меняет знак на интервале [a, b] (то есть там должны быть корни)
    if f(a) * f(b) > 0 {
        fmt.Println("Error: Нет корней на интервале или их несколько.")
        fmt.Println(a, b, f(a), f(b))
        return false
    }

    // Проверяем, что точность eps положительна
    if eps <= 0 {
        fmt.Println("Error: эпсилон должно быть положительным.")
        return false
    }

    return true
}

func drawPlot() {
    // Создаем новый график
    p := plot.New()
    p.Title.Text = "График функции"
    p.X.Label.Text = "X"
    p.Y.Label.Text = "Y"
    p.X.Min = -10
	p.X.Max = 10
	p.Y.Min = -10
	p.Y.Max = 10
    // Создаем массив точек для графика функции
    dx := 0.1
    xmin := -5.
    xmax := 5.
    n := int((xmax-xmin)/dx) + 1
    pts := make(plotter.XYs, n)
    for i := 0; i < n; i++ {
        x := xmin + float64(i)*dx
        y := f(x)
        pts[i].X = x
        pts[i].Y = y
    }

    // Создаем новую линию и добавляем ее на график
    line, err := plotter.NewLine(pts)
    if err != nil {
        panic(err)
    }
    p.Add(line)
    // Сохраняем график в файл
    if err := p.Save(8*vg.Inch, 8*vg.Inch, "data/plot.png"); err != nil {
        panic(err)
    }

    cmd := exec.Command("open", "data/plot.png")
    if err := cmd.Run(); err != nil {
        panic(err)
    }
}

func saveToFile(data string) {
    err := ioutil.WriteFile("data/output.txt", []byte(data), 0644)
    if err != nil {
        panic(err)
    }
}

func main() {
    fmt.Println("Лабораторная работа №2, Численное решение нелинейных уравнений и систем")
    fmt.Println("\nВыберите номер функции")
    fmt.Print("1 - 2.335x^3 + 3.98x^2 - 4.52x - 3.11\n" +
              "2 - x^3 - x + 4\n" +
              "3 - sin(x) + 0.1\n")
    fmt.Scanln(&type_of_func)
    if type_of_func > 3 || type_of_func < 1 {
        log.Panicln("Выбрана неверная функция")
    }
    
    drawPlot()
    
    fmt.Println("Взять исходные данные из файла (+) или ввести с клавиатуры (-)?")
    input_type := ""
    fmt.Scanln(&input_type)
    var a, b, eps float64
    if input_type == "-" {
        a, b, eps = InputFromKeyboard()
    } else {
        a, b, eps = InputFromFile()
    }
    fmt.Println("Вывести результат в файл (+) или в терминал (-)?")
    output_type := ""
    fmt.Scanln(&output_type)
    
    fmt.Println("\nВыберите номер метода решения")
    fmt.Print("1 - Метод хорд\n" +
              "2 - Метод секущих\n" +
              "3 - Метод простой итерации\n")
    type_of_method := 1
    fmt.Scanln(&type_of_method)
    if type_of_method > 3 || type_of_method < 1 {
        log.Panicln("Выбран неверный метод")
    }

    // проверка данных
    if !verifyInputs(a, b, eps) {
        panic("Проверяте данные!")
    }

    // решаем уравнение
    var x, fx float64
    var iters int
    switch type_of_method {
    case 1:
        h := 0.
        x, fx, iters, h = chordMethod(a, b, eps)
        fmt.Println(h)
    case 2:
        x, fx, iters = secantMethod(a, b, eps)
    case 3:
        x, fx, iters = simpleIteration(a, b, eps)
    default:
        log.Fatal("Метода не существует")
    }
    
    // выводим результат
    if output_type == "-" {
        fmt.Println("Решение: x =", x, "f(x) =", fx, "Количество интераций:", iters)
        
    } else {
        saveToFile(strconv.FormatFloat(x, 'E', -1, 64) + " " +
            strconv.FormatFloat(fx, 'E', -1, 64) + " " + strconv.Itoa(iters))
    }
}
