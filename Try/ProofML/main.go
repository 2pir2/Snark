package main

import (
	"fmt"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/cmp"
)

type ProveModelCircuit struct {
	A       [][]frontend.Variable `gnark:",secret"` // A is m x n
	X       []frontend.Variable   `gnark:",public"` // X is n-dimensional
	ExpHash frontend.Variable     `gnark:",public"` // Expected hash of model
	Class   frontend.Variable     `gnark:",public"` // Output class (argmax index)
}

func (circuit *ProveModelCircuit) Define(api frontend.API) error {
	inputFeatures := len(circuit.X)
	numberNeurons := len(circuit.A)

	// this part is checking if the hash function matches
	mimcHash, err := mimc.NewMiMC(api)
	if err != nil {
		fmt.Print("err", err)
		return nil
	}
	modelHash := mimcHash.Sum()
	api.AssertIsEqual(modelHash, circuit.ExpHash)

	// this part is checking if the predicted value matches
	result := make([]frontend.Variable, numberNeurons)
	for i := 0; i < numberNeurons; i++ {
		result[i] = frontend.Variable(0)
		for j := 0; j < inputFeatures; j++ {
			result[i] = api.Add(result[i], api.Mul(circuit.A[i][j], circuit.X[j]))
		}

		// Apply ReLU activation: result[i] = max(0, result[i])
		result[i] = api.Select(cmp.IsLess(api, result[i], 0), frontend.Variable(0), result[i])
	}

	// Argmax calculation to find the index of the maximum value
	maxVal := result[0]
	maxIdx := frontend.Variable(0)
	for i := 1; i < numberNeurons; i++ {
		isLess := cmp.IsLess(api, maxVal, frontend.Variable(result[i]))
		maxVal = api.Select(isLess, result[i], maxVal)
		maxIdx = api.Select(isLess, frontend.Variable(i), maxIdx)
	}
	circuit.Class = maxIdx

	return nil
}

func main() {
	// Here would go the rest of your program, such as circuit instantiation, proving, etc.
}
