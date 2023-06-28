package nested_types

import (
	"github.com/antosdaniel/mappergen/test/testdata/nested_types/dst"
	"github.com/antosdaniel/mappergen/test/testdata/nested_types/src"
)

func ToPayslip(src src.Payslip) dst.Payslip {
	return dst.Payslip{
		ID: src.ID,
		Payroll: dst.Payroll{
			ID:     src.Payroll.ID,
			Payday: src.Payroll.Payday,
			Approver: dst.Approver{
				ID:   src.Payroll.Approver.ID,
				Name: src.Payroll.Approver.Name,
			},
		},
	}
}
