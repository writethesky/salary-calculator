package calculator

type Tax struct {
	From                      int
	To                        int
	Rate                      int
	QuickCalculationDeduction int
}

func getTaxTable() []Tax {
	return []Tax{
		{
			From:                      0,
			To:                        36000,
			Rate:                      30,
			QuickCalculationDeduction: 0,
		},
		{
			From:                      36000,
			To:                        144000,
			Rate:                      100,
			QuickCalculationDeduction: 2520,
		},
		{
			From:                      144000,
			To:                        300000,
			Rate:                      200,
			QuickCalculationDeduction: 16920,
		},
		{
			From:                      300000,
			To:                        420000,
			Rate:                      250,
			QuickCalculationDeduction: 31920,
		},
		{
			From:                      420000,
			To:                        660000,
			Rate:                      300,
			QuickCalculationDeduction: 52920,
		},
		{
			From:                      660000,
			To:                        960000,
			Rate:                      350,
			QuickCalculationDeduction: 85920,
		},
		{
			From:                      960000,
			To:                        96000000000000000,
			Rate:                      450,
			QuickCalculationDeduction: 181920,
		},
	}
}
