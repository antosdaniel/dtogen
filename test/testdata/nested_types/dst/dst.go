package dst

import "time"

type Payslip struct {
	ID      string
	Payroll Payroll
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
