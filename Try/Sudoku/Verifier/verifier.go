package main

import (
	"os"
	"fmt"
	"encoding/json"

	"github.com/consensys/gnark/backend/groth16"
	
	

)



func main() {
	
	vkJSON, _ := os.ReadFile("gnark_vk.json");
	var vk []byte
	_ = json.Unmarshal(vkJSON, &vk);
	fmt.Println(vk)

	proofJSON, _ := os.ReadFile("gnark_proof.json");
	var proof groth16.Proof
	_ = json.Unmarshal(proofJSON, &proof)
	fmt.Println(proof);
}
