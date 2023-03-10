package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func Solve(matrix [][]float64, answers []float64) (string, []float64) {
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
                fmt.Println("Система имеет множество решений")
            } else {
                fmt.Println("Система не имеет решений")
            }
            return "error", answers
        }

        // деление элементов строки на главный элемент
        //for j := 0; j < n; j++ {
        //    matrix[i][index[j]] /= main_elem
        //}
        //answers[i] /= main_elem

        // вычитание текущей строки из всех ниже расположенных строк с занулением i-ого элемента в каждой из них
        for j := i + 1; j < n; j++ {
            main_elem = matrix[j][index[i]]
            p := main_elem/matrix[i][index[i]]
            for m := 0; m < n; m++ {
                //matrix[j][index[m]] -= matrix[i][index[m]]*main_elem
                matrix[j][index[m]] -= matrix[i][index[m]] * p
            }
            answers[j] -= answers[i]*p
        }
    }
    fmt.Println("Треугольная Матрица:")
    for i := range matrix {
        for j := range matrix[i] {
            fmt.Printf("%f ", matrix[i][index[j]])
        }
        fmt.Printf("| %f\n", answers[i])
    }
    det := math.Pow(-1, per)
    for i := 0; i < n; i++ {
        det *= matrix[i][index[i]]
    }
    fmt.Println("Определитель по матрице", det)
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
    // вычисление невязок
    r := make([]float64, n)
    for i := range matrix {
        r[i] = answers[i]
        for j := range matrix[i] {
            r[i] -= matrix[i][index[j]] * result[j]
        }
    }
    fmt.Println("Вектор невязок:")
    for i := range matrix {
        fmt.Printf("%.20f, ", r[i])
    }
    fmt.Print("\n\n")
    return "ok", result
}


func FindDeterminant(matrix [][]float64) (float64) {
    n := len(matrix)
    if n == 1 {
        return matrix[0][0]
    }
    var det float64 = 0
    var sign int = 1
    for i := 0; i < n; i++ {
        det += float64(sign) * matrix[0][i] * FindDeterminant(FindMinor(matrix, i))
        sign *= -1
    }
    return det
}


func FindMinor(matrix [][]float64, i int) ([][]float64) {
    n := len(matrix)
    var res_matrix [][]float64
    for row := 1; row < n; row++ {
        var matrix_row []float64
        for col := 0; col < n; col++ {
            if col == i {
                continue
            }
            matrix_row = append(matrix_row, matrix[row][col])
        }
        res_matrix = append(res_matrix, matrix_row)
    } 
    return res_matrix
}


func InputFromFile() ([][]float64, []float64) {

    b, err := os.ReadFile("input.txt")
    if err != nil {
        panic(err)
    }
    data := string(b)
    
    input := strings.FieldsFunc(data, func (r rune) bool {
        return r == ' ' || r == '\n'
    })
    n, _ := strconv.Atoi(input[0])
    var matrix [][]float64
    answers := make([]float64, n)
    for i := 0; i < n; i++ {
        row := make([]float64, n)
        for j := 0; j < n; j++ {
            row[j], _ = strconv.ParseFloat(input[i*n+1+j], 64)
        }
        matrix = append(matrix, row)
    }
    for i := 0; i < n; i++ {
        answers[i], _ = strconv.ParseFloat(input[n*n+1+i], 64)
    }

    return matrix, answers
}


func InputFromKeyboard() ([][]float64, []float64) {
    n := 0
    var matrix [][]float64
    fmt.Println("Размер матрицы:")
    fmt.Scanln(&n)
    fmt.Println("Введите матрицу коэффициентов:")
    for i := 0; i < n; i++ {
        input := make([]float64, n)
        for j := 0; j < n; j++ {
            fmt.Scanf("%f", &input[j])
        }
        matrix = append(matrix, input)
    }
    fmt.Println("Введите матрицу ответов:")
    answers := make([]float64, n)
    for i := 0; i < n; i++ {
        
        fmt.Scanf("%f", &answers[i])
    }
    return matrix, answers
}


func main() {
    fmt.Println("Вариант 27")
    fmt.Println("Введите + для ввода из файла и любой другой символ для ввода в консоли:")
    input_type := ""
    fmt.Scanln(&input_type)
    var matrix [][]float64
    var answers []float64
    if input_type == "-" {
        matrix, answers = InputFromKeyboard()
    } else {
        matrix, answers = InputFromFile()
    }
    det := FindDeterminant(matrix)
    fmt.Println("Определитель равен:", det)
    if det == 0 {
        fmt.Println("Система является несовместной.")
        //return
    }
    state, result := Solve(matrix, answers)
    if state == "error" {
        return
    }
    fmt.Println("Ответ:")
    for i := range result {
        fmt.Printf("x%d: %f; ", i, result[i])
    }
}
