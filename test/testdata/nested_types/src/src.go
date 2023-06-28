package src

import "time"

type Payslip struct {
	ID       string
	GrossPay int
	Tax      int
	NetPay   int

	Payroll  Payroll
	Employee Employee
}

type Payroll struct {
	ID     string
	Payday time.Time

	Approver Approver
}

type Approver struct {
	ID   string
	Name string
}

type Employee struct {
	ID   string
	Name string
}
