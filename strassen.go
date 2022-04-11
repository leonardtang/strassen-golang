package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func timer(start time.Time, process string) time.Duration {
	totalTime := time.Since(start)
	log.Printf("%s took %s", process, totalTime)
	return totalTime
}

func makeMatrix(size int) [][]int {

	mat := make([][]int, size)
	for i := 0; i < size; i++ {
		mat[i] = make([]int, size)
	}
	return mat
}

func populateMatrix(size int) [][]int {

	mat := makeMatrix(size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			mat[i][j] = rand.Intn(1000000)
		}
	}
	return mat
}

func add(A [][]int, B [][]int) [][]int {
	sum := makeMatrix(len(A))
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A); j++ {
			sum[i][j] = A[i][j] + B[i][j]
		}
	}
	return sum
}

func subtract(A [][]int, B [][]int) [][]int {
	diff := makeMatrix(len(A))
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A); j++ {
			diff[i][j] = A[i][j] - B[i][j]
		}
	}
	return diff
}

func basicMultiply(A [][]int, B [][]int) [][]int {

	size := len(A)
	prod := makeMatrix(size)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			prod[i][j] = 0
			for k := 0; k < size; k++ {
				prod[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return prod
}

func strassen(A [][]int, B [][]int, threshold int) [][]int {

	if len(A) <= threshold && len(B) <= threshold {
		return basicMultiply(A, B)
	}

	var padded bool
	if len(A)%2 == 1 {

		var row []int
		for i := 0; i < len(A); i++ {
			A[i] = append(A[i], 0)
			row = append(row, 0)
		}
		row = append(row, 0)
		A = append(A, row)

		for i := 0; i < len(B); i++ {
			B[i] = append(B[i], 0)
		}
		B = append(B, row)
		padded = true

	}

	size := len(A)

	a := makeMatrix(size / 2)
	b := makeMatrix(size / 2)
	c := makeMatrix(size / 2)
	d := makeMatrix(size / 2)
	e := makeMatrix(size / 2)
	f := makeMatrix(size / 2)
	g := makeMatrix(size / 2)
	h := makeMatrix(size / 2)

	for i := 0; i < size/2; i++ {
		for j := 0; j < size/2; j++ {
			a[i][j] = A[i][j]
			b[i][j] = A[i][j+size/2]
			c[i][j] = A[i+size/2][j]
			d[i][j] = A[i+size/2][j+size/2]
			e[i][j] = B[i][j]
			f[i][j] = B[i][j+size/2]
			g[i][j] = B[i+size/2][j]
			h[i][j] = B[i+size/2][j+size/2]
		}
	}

	sub1 := subtract(f, h)
	sub2 := add(a, b)
	sub3 := add(c, d)
	sub4 := subtract(g, e)
	sub5 := add(a, d)
	sub6 := add(e, h)
	sub7 := subtract(b, d)
	sub8 := add(g, h)
	sub9 := subtract(a, c)
	sub10 := add(e, f)

	p1 := strassen(a, sub1, threshold)
	p2 := strassen(sub2, h, threshold)
	p3 := strassen(sub3, e, threshold)
	p4 := strassen(d, sub4, threshold)
	p5 := strassen(sub5, sub6, threshold)
	p6 := strassen(sub7, sub8, threshold)
	p7 := strassen(sub9, sub10, threshold)

	out11 := subtract(add(p5, add(p4, p6)), p2)
	out12 := add(p1, p2)
	out21 := add(p3, p4)
	out22 := subtract(subtract(add(p1, p5), p3), p7)

	out := makeMatrix(size)

	for i := 0; i < size/2; i++ {
		for j := 0; j < size/2; j++ {
			out[i][j] += out11[i][j]
			out[i][j+size/2] += out12[i][j]
			out[i+size/2][j] += out21[i][j]
			out[i+size/2][j+size/2] += out22[i][j]
		}
	}

	if padded == true {
		out = out[0 : size-1]
		for i := 0; i < size-1; i++ {
			out[i] = out[i][0 : size-1]
		}
	}

	return out
}

func triangle(p float64) int {
	n := 1024
	t := 32
	graph := makeMatrix(n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if j == i {
				graph[i][j] = 0
			} else {
				v := rand.Float64()
				if v <= p {
					graph[i][j] = 1
					graph[j][i] = 1
				} else {
					graph[i][j] = 0
					graph[j][i] = 0
				}
			}

		}
	}
	graphsq := strassen(graph, graph, t)
	triangle := strassen(graphsq, graph, t)

	count := 0
	for i := 0; i < n; i++ {
		count = count + triangle[i][i]
	}
	count = count / 6
	return count
}

// Triangle counting in random graphs via Strassen's
func triangles() {
	p := [5]float64{0.01, 0.02, 0.03, 0.04, 0.05}

	for i := 0; i < len(p); i++ {
		expected := 178433024 * math.Pow(p[i], 3)
		avgCount := 0
		for j := 0; j < 5; j++ {
			avgCount += triangle(p[i])
		}
		log.Print("p = ", p[i], " Actual: ", avgCount/5, " Expected: ", expected)
	}
}

func readText(dim int, path string) ([][]int, [][]int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileBytes, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		log.Fatal(err)
	}

	sliceData := strings.Split(string(fileBytes), "\n")
	scanner := bufio.NewScanner(file)
	sliceA := make([]int, len(sliceData)/2)
	sliceB := make([]int, len(sliceData)/2)

	lineNum := 0
	for scanner.Scan() {
		lineStr := scanner.Text()
		num, _ := strconv.Atoi(lineStr)
		if lineNum < len(sliceData)/2 {
			sliceA[lineNum%(len(sliceData)/2)] = num
		} else {
			sliceB[lineNum%(len(sliceData)/2)] = num
		}
		lineNum += 1
	}

	A := makeMatrix(dim)
	B := makeMatrix(dim)
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			A[i][j] = sliceA[i*dim+j]
			B[i][j] = sliceB[i*dim+j]
		}
	}

	return A, B
}

func printDiag(mat [][]int) {
	for i := 0; i < len(mat); i++ {
		fmt.Println(mat[i][i])
	}
}

func crossExperiments() {
	dims := []int{512, 1024, 2048, 513, 1025, 2049}
	crossovers := []int{8, 16, 32, 64, 128, 256}
	for _, dim := range dims {
		A := populateMatrix(dim)
		B := populateMatrix(dim)
		for i := 0; i < len(crossovers); i++ {
			start := time.Now()
			strassen(A, B, crossovers[i])
			timer(start, fmt.Sprintf("Crossover %d for dim %d", crossovers[i], dim))
		}
	}
}

func main() {

	args := os.Args
	if len(args) > 1 {
		dim, _ := strconv.Atoi(args[2])
		inputText := args[3]
		A, B := readText(dim, inputText)
		prod := strassen(A, B, 32)
		printDiag(prod)
	} else {
		crossExperiments()
		triangles()
	}
}
