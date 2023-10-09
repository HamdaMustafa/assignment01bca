package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction represents a single transaction.
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Timestamp string
}

// Block represents a single block in the blockchain.
type Block struct {
	Index         int
	Timestamp     string
	Transactions  []Transaction
	Nonce         int
	PrevHash      string
	Hash          string
}

// Blockchain represents the blockchain as a slice of blocks.
var Blockchain []Block

const difficulty = 2 // Number of leading zeros required in the hash (difficulty level)

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.PrevHash
	for _, transaction := range block.Transactions {
		record += transaction.Sender + transaction.Receiver + fmt.Sprintf("%.2f", transaction.Amount) + transaction.Timestamp
	}
	record += string(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func createBlock(transactions []Transaction, nonce int, previousHash string) *Block {
	newBlock := &Block{
		Index:        len(Blockchain),
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		Nonce:        nonce,
		PrevHash:     previousHash,
	}
	newBlock.Hash = calculateHash(*newBlock)
	return newBlock
}

func mineBlock(transactions []Transaction, previousHash string) *Block {
	var newBlock *Block
	nonce := 0
	for {
		newBlock = createBlock(transactions, nonce, previousHash)
		hash := calculateHash(*newBlock)

		// Simulate mining by requiring the hash to start with leading zeros
		if hash[:difficulty] == "00" {
			break
		}

		nonce++
	}
	return newBlock
}

func DisplayBlocks() {
	for _, block := range Blockchain {
		fmt.Printf("Block #%d:\n", block.Index)
		fmt.Printf("  Transactions:\n")
		for _, transaction := range block.Transactions {
			fmt.Printf("    Sender: %s\n", transaction.Sender)
			fmt.Printf("    Receiver: %s\n", transaction.Receiver)
			fmt.Printf("    Amount: %.2f\n", transaction.Amount)
		}
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
	}
}

func VerifyChain() bool {
	for i := 1; i < len(Blockchain); i++ {
		if Blockchain[i].PrevHash != Blockchain[i-1].Hash {
			return false
		}
		if calculateHash(Blockchain[i]) != Blockchain[i].Hash {
			return false
		}
	}
	return true
}

func ChangeBlock(index int, newTransaction Transaction) {
	if index >= 0 && index < len(Blockchain) {
		Blockchain[index].Transactions = []Transaction{newTransaction}
		Blockchain[index].Hash = calculateHash(Blockchain[index])
	}
}

func main() {
	genesisBlock := createBlock([]Transaction{Transaction{Sender: "Genesis", Receiver: "Alice", Amount: 100.00}}, 0, "")
	Blockchain = append(Blockchain, *genesisBlock)

	fmt.Println("Genesis Block added to the blockchain")

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Add a Block")
		fmt.Println("2. Display the whole chain")
		fmt.Println("3. Change a Block")
		fmt.Println("4. Verify Chain")
		fmt.Println("5. Quit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var newTransactions []Transaction
			fmt.Print("Enter sender name: ")
			var sender, receiver string
			var amount float64

			fmt.Scanln(&sender)

			fmt.Print("Enter receiver name: ")
			fmt.Scanln(&receiver)

			fmt.Print("Enter amount: ")
			fmt.Scanln(&amount)

			transaction := Transaction{
				Sender:    sender,
				Receiver:  receiver,
				Amount:    amount,
				Timestamp: time.Now().String(),
			}

			newTransactions = append(newTransactions, transaction)

			newBlock := mineBlock(newTransactions, Blockchain[len(Blockchain)-1].Hash)
			Blockchain = append(Blockchain, *newBlock)

			fmt.Printf("Block #%d added to the blockchain\n", newBlock.Index)
			fmt.Printf("Hash: %s\n", newBlock.Hash)

		case 2:
			fmt.Println("\nBlockchain:")
			DisplayBlocks()

		case 3:
			fmt.Print("Enter the index of the block to change: ")
			var index int
			fmt.Scanln(&index)

			if index >= 0 && index < len(Blockchain) {
				var sender, receiver string
				var amount float64

				fmt.Print("Enter new sender name: ")
				fmt.Scanln(&sender)

				fmt.Print("Enter new receiver name: ")
				fmt.Scanln(&receiver)

				fmt.Print("Enter new amount: ")
				fmt.Scanln(&amount)

				newTransaction := Transaction{
					Sender:    sender,
					Receiver:  receiver,
					Amount:    amount,
					Timestamp: time.Now().String(),
				}

				ChangeBlock(index, newTransaction)
				fmt.Printf("Block #%d changed\n", index)
			} else {
				fmt.Println("Invalid block index")
			}

		case 4:
			if VerifyChain() {
				fmt.Println("Blockchain is valid.")
			} else {
				fmt.Println("Blockchain is invalid.")
			}

		case 5:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
