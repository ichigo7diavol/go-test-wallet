package wallet

// Wallet represents a simple wallet
type Wallet struct {
	balance int
}

// NewWallet creates a new wallet with zero balance
func NewWallet() *Wallet {
	return &Wallet{balance: 0}
}

// Balance returns the current balance
func (w *Wallet) Balance() int {
	return w.balance
}

// Deposit adds money to the wallet
func (w *Wallet) Deposit(amount int) {
	w.balance += amount
}

// Withdraw removes money from the wallet if sufficient balance
func (w *Wallet) Withdraw(amount int) bool {
	if w.balance >= amount {
		w.balance -= amount
		return true
	}
	return false
}
