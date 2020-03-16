// This script illustrates Go concurrency for matrix multiplication.
//
// This is a demonstration, not for production use.
//
// Below we consider several approaches for multiplying matrices concurrently.
// The best-perfoming method can be seen to depend on the dimensions of the
// matrices being  multiplied.
//
// We are multiplying a p x q matrix by a q x r matrix. A summary of the
// results of this comparison is:
//
// If p, q, r are all small (less than ~100) then a non-concurrent approach is
// as fast or faster than a concurrent approach.
//
// If q is large but p and r are small, the inner-product concurrent approach
// is fastest.
//
// If p and r are large but q is small, the outer-product concurrent approach
// is fastest.
//
// Note that different approaches use different methods for data synchronization.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gonum.org/v1/gonum/floats"
)

const (
	// We will multiply a p x q matrix (mat1) by a q x r matrix (mat2').
	p = 1000
	q = 200
	r = 3000
)

// genmat produces a matrix filled with iid standard normal values.
func genmat(p, q int) []float64 {
	x := make([]float64, p*q)

	for i := range x {
		x[i] = rand.NormFloat64()
	}

	return x
}

// basicMul multiplies two matrices non-concurrently.
func basicMul(mat1, mat2 []float64) []float64 {

	prod := make([]float64, p*r)

	for i := 0; i < p; i++ {
		for j := 0; j < r; j++ {
			ii := i*r + j
			for k := 0; k < q; k++ {
				prod[ii] += mat1[i*q+k] * mat2[j*q+k]
			}
		}
	}

	return prod
}

// concurrentMulInner multiplies two matrices concurrently using
// an inner-product algorithm.  There is no data sharing here,
// so no synchronization is needed.
func concurrentMulInner(mat1, mat2 []float64) []float64 {

	var wg sync.WaitGroup
	prod := make([]float64, p*r)

	dot := func(i, j int) {
		ii := i*r + j
		for k := 0; k < q; k++ {
			prod[ii] += mat1[i*q+k] * mat2[j*q+k]
		}

		wg.Done()
	}

	for i := 0; i < p; i++ {
		for j := 0; j < r; j++ {
			wg.Add(1)
			go dot(i, j)
		}
	}

	wg.Wait()

	return prod
}

// concurrentMulOuterBad is a concurrent outer-product based
// matrix multiplication function that performs very badly due
// to extreme lock contention.
func concurrentMulOuterBad(mat1, mat2 []float64) []float64 {

	prod := make([]float64, p*r)
	var wg sync.WaitGroup
	var sy sync.Mutex

	outer := func(k int) {
		for i := 0; i < p; i++ {
			for j := 0; j < r; j++ {
				ii := i*r + j
				sy.Lock()
				prod[ii] += mat1[i*q+k] * mat2[j*q+k]
				sy.Unlock()
			}
		}
		wg.Done()
	}

	for k := 0; k < q; k++ {
		wg.Add(1)
		go outer(k)
	}

	wg.Wait()

	return prod
}

// concurrentMulOuter multiplies two matrices concurrently using outer
// products.  A mutex is used to synchronize data sharing.
func concurrentMulOuter(mat1, mat2 []float64) []float64 {

	prod := make([]float64, p*r)
	var wg sync.WaitGroup
	var sy sync.Mutex

	outer := func(k int) {
		mat := make([]float64, p*r)
		for i := 0; i < p; i++ {
			for j := 0; j < r; j++ {
				mat[i*r+j] = mat1[i*q+k] * mat2[j*q+k]
			}
		}
		sy.Lock()
		floats.Add(prod, mat)
		sy.Unlock()
		wg.Done()
	}

	for k := 0; k < q; k++ {
		wg.Add(1)
		go outer(k)
	}

	wg.Wait()

	return prod
}

// concurrentMulOuter multiplies two matrices concurrently using outer
// products.  Channels are used to synchronize data sharing.
func concurrentMulOuterChan(mat1, mat2 []float64) []float64 {

	prod := make([]float64, p*r)
	var wg sync.WaitGroup
	rc := make(chan []float64)

	go func() {
		for i := 0; i < q; i++ {
			x := <-rc
			floats.Add(prod, x)
		}
	}()

	outer := func(k int) {
		mat := make([]float64, p*r)
		for i := 0; i < p; i++ {
			for j := 0; j < r; j++ {
				mat[i*r+j] = mat1[i*q+k] * mat2[j*q+k]
			}
		}
		rc <- mat
		wg.Done()
	}

	for k := 0; k < q; k++ {
		wg.Add(1)
		go outer(k)
	}

	wg.Wait()

	return prod
}

func main() {

	mat1 := genmat(p, q)
	mat2 := genmat(r, q)

	var prod [][]float64

	//methods := []func() []float64{basicMul, concurrentMulInner, concurrentMulOuterBad, concurrentMulOuter,
	//concurrentMulOuterChan}
	methods := []func([]float64, []float64) []float64{basicMul, concurrentMulInner, concurrentMulOuter,
		concurrentMulOuterChan}

	for _, f := range methods {
		start := time.Now()
		x := f(mat1, mat2)
		duration := time.Now().Sub(start)
		fmt.Printf("Duration: %v\n", duration)
		prod = append(prod, x)
	}

	for i, x := range prod {
		for j, y := range prod {
			if !floats.EqualApprox(x, y, 1e-12) {
				fmt.Printf("Some results differ (%d, %d)\n", i, j)
			}
		}
	}
}
