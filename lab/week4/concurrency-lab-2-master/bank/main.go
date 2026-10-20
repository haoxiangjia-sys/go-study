package main

import (
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var debug *bool

// An executor is a type of a worker goroutine that handles the incoming transactions.
// It only receives transactions whose two accounts are already locked by the manager.
func executor(bank *bank, executorId int, ready <-chan transaction, finished chan<- transaction, done chan<- bool) {
	for {
		t := <-ready

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		fmt.Println("Executor\t", executorId, "executing transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		bank.execute(t, executorId)

		bank.removeCompleted(e, executorId) // Removing this line will break visualisations.
		finished <- t                       // tell the manager both accounts can be released
		done <- true
	}
}

// manager is the only goroutine that locks and unlocks accounts.
// It hands a transaction to the executors only when both of its accounts are free.
func manager(bank *bank, transactionQueue <-chan transaction, ready chan<- transaction, finished <-chan transaction, total int) {
	busy := make([]bool, len(bank.accounts)) // only the manager touches this
	pending := []transaction{}
	inFlight := 0
	received := 0

	for received < total || len(pending) > 0 || inFlight > 0 {
		// 1. Pull every transaction currently waiting in the queue.
		for received < total && len(transactionQueue) > 0 {
			pending = append(pending, <-transactionQueue)
			received++
		}

		// 2. Dispatch every pending transaction whose two accounts are free.
		remaining := pending[:0]
		for _, t := range pending {
			if !busy[t.from] && !busy[t.to] {
				busy[t.from], busy[t.to] = true, true
				bank.lockAccount(t.from, "M")
				bank.lockAccount(t.to, "M")
				ready <- t
				inFlight++
			} else {
				remaining = append(remaining, t)
			}
		}
		pending = remaining

		if inFlight == 0 {
			// Nothing running and nothing dispatchable: wait for a new transaction
			// instead of spinning in the loop (no busy waiting).
			if received < total {
				pending = append(pending, <-transactionQueue)
				received++
			}
			continue
		}

		// 3. Nothing else can start: block (no busy waiting) until a transaction finishes.
		t := <-finished
		bank.unlockAccount(t.from, "M")
		bank.unlockAccount(t.to, "M")
		busy[t.from], busy[t.to] = false, false
		inFlight--
	}
}

func toChar(i int) rune {
	return rune('A' + i)
}

// main creates a bank and executors that will be handling the incoming transactions.
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	debug = flag.Bool("debug", false, "generate DOT graphs of the state of the bank")
	flag.Parse()

	bankSize := 6 // Must be even for correct visualisation.
	transactions := 1000

	accounts := make([]*account, bankSize)
	for i := range accounts {
		accounts[i] = &account{name: string(toChar(i)), balance: 1000}
	}

	bank := bank{
		accounts:               accounts,
		transactionsInProgress: list.New(),
		gen:                    newGenerator(),
	}

	startSum := bank.sum()

	transactionQueue := make(chan transaction, transactions)
	expectedMoneyTransferred := 0
	for i := 0; i < transactions; i++ {
		t := bank.getTransaction()
		expectedMoneyTransferred += t.amount
		transactionQueue <- t
	}

	done := make(chan bool)
	ready := make(chan transaction, transactions)    // manager -> executors
	finished := make(chan transaction, transactions) // executors -> manager

	go manager(&bank, transactionQueue, ready, finished, transactions)

	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, ready, finished, done)
	}

	for total := 0; total < transactions; total++ {
		fmt.Println("Completed transactions\t", total)
		<-done
	}

	fmt.Println()
	fmt.Println("Expected transferred", expectedMoneyTransferred)
	fmt.Println("Actual transferred", bank.moneyTransferred)
	fmt.Println("Expected sum", startSum)
	fmt.Println("Actual sum", bank.sum())
	if bank.sum() != startSum {
		panic("sum of the account balances does not much the starting sum")
	} else if len(transactionQueue) > 0 {
		panic("not all transactions have been executed")
	} else if bank.moneyTransferred != expectedMoneyTransferred {
		panic("incorrect amount of money was transferred")
	} else {
		fmt.Println("The bank works!")
	}
}
