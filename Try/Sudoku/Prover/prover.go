package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type SudokuCircuit struct {
	IncompleteGrid [9][9]frontend.Variable `gnark:"IncompleteSudoku,public"`
	CompleteGrid   [9][9]frontend.Variable `gnark:"CompleteSudoku"`
}

type Sudoku struct {
	Grid [9][9]int `json:"grid"`
}

func (circuit *SudokuCircuit) Define(api frontend.API) error {
	// Constraint 1: Each cell value in the CompleteGrid must be between 1 and 9
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			api.AssertIsLessOrEqual(circuit.CompleteGrid[i][j], 9)
			api.AssertIsLessOrEqual(1, circuit.CompleteGrid[i][j])
		}
	}

	// Constraint 2: Each row in the CompleteGrid must contain unique values
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := j + 1; k < 9; k++ {
				api.AssertIsDifferent(circuit.CompleteGrid[i][j], circuit.CompleteGrid[i][k])
			}
		}
	}

	// Constraint 3: Each column in the CompleteGrid must contain unique values
	for j := 0; j < 9; j++ {
		for i := 0; i < 9; i++ {
			for k := i + 1; k < 9; k++ {
				api.AssertIsDifferent(circuit.CompleteGrid[i][j], circuit.CompleteGrid[k][j])
			}
		}
	}

	// Constraint 4: Each 3x3 sub-grid in the CompleteGrid must contain unique values
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			for i := 0; i < 9; i++ {
				for j := i + 1; j < 9; j++ {
					row1 := boxRow*3 + i/3
					col1 := boxCol*3 + i%3
					row2 := boxRow*3 + j/3
					col2 := boxCol*3 + j%3
					api.AssertIsDifferent(circuit.CompleteGrid[row1][col1], circuit.CompleteGrid[row2][col2])
				}
			}
		}
	}

	// Constraint 5: The values in the IncompleteGrid must match the CompleteGrid where provided
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			isCellGiven := api.IsZero(circuit.IncompleteGrid[i][j])
			api.AssertIsEqual(api.Select(isCellGiven, circuit.CompleteGrid[i][j], circuit.IncompleteGrid[i][j]), circuit.CompleteGrid[i][j])
		}
	}

	return nil
}

func main() {
	// Load the JSON files
	incompleteFile, err := ioutil.ReadFile("public.json")
	if err != nil {
		fmt.Println("Error reading incomplete file:", err)
		return
	}

	completeFile, err := ioutil.ReadFile("private.json")
	if err != nil {
		fmt.Println("Error reading complete file:", err)
		return
	}

	// Unmarshal the byte slices into the Sudoku structs
	var incompleteSudoku Sudoku
	err = json.Unmarshal(incompleteFile, &incompleteSudoku)
	if err != nil {
		fmt.Println("Error unmarshalling incomplete JSON:", err)
		return
	}

	var completeSudoku Sudoku
	err = json.Unmarshal(completeFile, &completeSudoku)
	if err != nil {
		fmt.Println("Error unmarshalling complete JSON:", err)
		return
	}

	// Create the circuit assignment
	assignment := &SudokuCircuit{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			assignment.IncompleteGrid[i][j] = frontend.Variable(incompleteSudoku.Grid[i][j])
			assignment.CompleteGrid[i][j] = frontend.Variable(completeSudoku.Grid[i][j])
		}
	}

	// Debugging: Print the assignment to ensure it's properly initialized
	fmt.Printf("Assignment: %+v\n", assignment)

	var myCircuit SudokuCircuit
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		fmt.Println("Error creating witness:", err)
		return
	}

	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &myCircuit)
	if err != nil {
		fmt.Println("Error compiling circuit:", err)
		return
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		fmt.Println("Error during setup:", err)
		return
	}

	proof, errproof := groth16.Prove(cs, pk, witness)
	if errproof != nil {
		fmt.Println("Error Proving: ", errproof)
		return
	}

	// Save the proof to a binary file
	proofFile, err := os.Create("proof.bin")
	if err != nil {
		fmt.Println("Error creating proof file:", err)
		return
	}
	defer proofFile.Close()

	_, err = proof.WriteTo(proofFile)
	if err != nil {
		fmt.Println("Error writing proof to file:", err)
		return
	}
	fmt.Println("Proof saved to proof.bin")

	// Save the verification key to a binary file
	vkFile, err := os.Create("vk.bin")
	if err != nil {
		fmt.Println("Error creating verification key file:", err)
		return
	}
	defer vkFile.Close()

	_, err = vk.WriteTo(vkFile)
	if err != nil {
		fmt.Println("Error writing verification key to file:", err)
		return
	}
	fmt.Println("Verification key saved to vk.bin")
}
