package main

import (
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var debug *bool

type manager struct {
	bank              *bank
	transactionQueue  <-chan transaction
	readyQueue        chan transaction // was <-chan transaction (cannot send/close)
	finished          <-chan transaction
	lockedAccounts    map[int]bool
	pending           []transaction
	totalTransactions int // was totalTransaction
	completed         int
}

func (m *manager) run() {
	for {
		select {
		case t, ok := <-m.transactionQueue:
			if !ok {
				m.transactionQueue = nil
				continue
			}
			m.schedule(t)
		case t := <-m.finished:
			m.unlock(t)
			m.completed++
			if m.completed == m.totalTransactions {
				close(m.readyQueue)
				return
			}
			m.schedulePending()
		}
	}
}

func (m *manager) schedule(t transaction) {
	if !m.lockedAccounts[t.from] && !m.lockedAccounts[t.to] {
		m.lockAndSend(t)
	} else {
		m.pending = append(m.pending, t)
	}
}

func (m *manager) lockAndSend(t transaction) {
	m.bank.lockAccount(t.from, "6")
	m.bank.lockAccount(t.to, "6")
	m.lockedAccounts[t.from] = true
	m.lockedAccounts[t.to] = true
	m.readyQueue <- t
}

func (m *manager) unlock(t transaction) {
	m.bank.unlockAccount(t.from, "6")
	m.bank.unlockAccount(t.to, "6")
	delete(m.lockedAccounts, t.from)
	delete(m.lockedAccounts, t.to)
}

func (m *manager) schedulePending() {
	var remaining []transaction
	for _, t := range m.pending {
		if !m.lockedAccounts[t.from] && !m.lockedAccounts[t.to] {
			m.lockAndSend(t)
		} else {
			remaining = append(remaining, t)
		}
	}
	m.pending = remaining
}

func executor(bank *bank, executorId int, readyQueue <-chan transaction, finished chan<- transaction, done chan<- bool) {
	for t := range readyQueue {
		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		fmt.Println("Executor\t", executorId, "attempting transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		bank.execute(t, executorId)

		bank.removeCompleted(e, executorId) // Removing this line will break visualisations.
		finished <- t
		done <- true
	}
}

func toChar(i int) rune {
	return rune('A' + i)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	debug = flag.Bool("debug", false, "generate DOT graphs of the state of the bank")
	flag.Parse()

	bankSize := 6
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
	close(transactionQueue)

	readyQueue := make(chan transaction, transactions)
	finished := make(chan transaction, transactions)
	done := make(chan bool)

	mgr := &manager{
		bank:              &bank,
		transactionQueue:  transactionQueue,
		readyQueue:        readyQueue,
		finished:          finished,
		lockedAccounts:    make(map[int]bool),
		pending:           []transaction{},
		totalTransactions: transactions, // matches struct field
	}
	go mgr.run()

	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, readyQueue, finished, done)
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
