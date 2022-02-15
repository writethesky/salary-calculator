package calculator

type Rule struct {
	Endowment               int // 养老保险比例
	Medical                 int // 医疗保险比例
	MedicalInput            int // 医疗保险进入个人账户比例
	MedicalPlus             int // 大病保险附加
	Birth                   int // 生育保险比例
	Unemployment            int // 失业保险比例
	IndustrialInjury        int // 工伤保险比例
	HousingAccumulationFund int // 住房公积金比例
	BaseLimit               int // 五险一金缴存基数上限
	FreeTax                 int // 每年免税额度
}

type Result struct {
	Person  PersonResult
	Company CompanyResult
}

type PersonResult struct {
	Expenditure
	MedicalInput            int   // 医疗保险进入个人账户
	InsuranceSalary         int   // 五险一金后工资
	MonthlyAverageTax       int   // 应缴纳个税(月均)
	MonthlyTax              []int // 应缴纳个税(每月)
	MonthlyAverageTaxSalary int   // 税后工资(月均)
	MonthlyTaxSalary        []int // 税后工资(每月)
	RealIncome              int   // 实际收入
}

type CompanyResult struct {
	Expenditure
}

type Expenditure struct {
	Endowment               int // 养老保险
	Medical                 int // 医疗保险
	Birth                   int // 生育保险
	Unemployment            int // 失业保险
	IndustrialInjury        int // 工伤保险
	HousingAccumulationFund int // 住房公积金
	Total                   int // 总计
}

func Calculate(salary int, personRule, companyRule Rule) Result {
	personExpenditure := calculateExpenditure(salary, personRule)
	companyExpenditure := calculateExpenditure(salary, companyRule)
	baseMoney := getBaseMoney(salary, personRule)

	insuranceSalary := salary - personExpenditure.Endowment - personExpenditure.Medical - personExpenditure.Birth - personExpenditure.Unemployment - personExpenditure.IndustrialInjury - personExpenditure.HousingAccumulationFund
	personTax := calculatePersonTax(insuranceSalary, personRule)
	taxSalary := insuranceSalary - personTax.MonthlyAverage
	medicalInput := baseMoney * (personRule.MedicalInput + companyRule.MedicalInput) / 1000

	personResult := PersonResult{
		Expenditure:             personExpenditure,
		MedicalInput:            medicalInput,
		InsuranceSalary:         insuranceSalary,
		MonthlyAverageTax:       personTax.MonthlyAverage,
		MonthlyTax:              personTax.Monthly,
		MonthlyAverageTaxSalary: taxSalary,
		MonthlyTaxSalary:        calculateMonthlyTaxSalary(insuranceSalary, personTax.Monthly),
		RealIncome:              taxSalary + personExpenditure.HousingAccumulationFund + companyExpenditure.HousingAccumulationFund + medicalInput,
	}

	companyResult := CompanyResult{
		Expenditure: companyExpenditure,
	}

	return Result{
		Person:  personResult,
		Company: companyResult,
	}
}

func calculateMonthlyTaxSalary(insuranceSalary int, monthlyTax []int) []int {
	monthlyTaxSalary := make([]int, 0, len(monthlyTax))
	for _, tax := range monthlyTax {
		monthlyTaxSalary = append(monthlyTaxSalary, insuranceSalary-tax)
	}
	return monthlyTaxSalary
}

type PersonTax struct {
	Annually       int
	Monthly        []int
	MonthlyAverage int
}

func calculatePersonTax(insuranceSalary int, rule Rule) PersonTax {
	taxes := getTaxTable()
	base := insuranceSalary*12 - rule.FreeTax
	totalTax := 0
	monthlyTax := make([]int, 0, 12)
	for i := len(taxes) - 1; i >= 0; i-- {
		if base > taxes[i].From {
			totalTax += (base - taxes[i].From) * taxes[i].Rate / 1000
			base = taxes[i].From
		}
	}

	cumulativeSalary := 0
	cumulativeTax := 0
	for i := 0; i < 12; i++ {
		cumulativeSalary += insuranceSalary - rule.FreeTax/12
		for j := len(taxes) - 1; j >= 0; j-- {
			if cumulativeSalary > taxes[j].From {
				tax := cumulativeSalary*taxes[j].Rate/1000 - taxes[j].QuickCalculationDeduction - cumulativeTax
				monthlyTax = append(monthlyTax, tax)
				cumulativeTax += tax
				break
			}
		}
	}

	return PersonTax{
		Annually:       totalTax,
		MonthlyAverage: totalTax / 12,
		Monthly:        monthlyTax,
	}
}

func calculateExpenditure(salary int, rule Rule) Expenditure {
	baseMoney := getBaseMoney(salary, rule)

	expenditure := Expenditure{
		Endowment:               baseMoney * rule.Endowment / 1000,
		Medical:                 baseMoney*rule.Medical/1000 + rule.MedicalPlus,
		Birth:                   baseMoney * rule.Birth / 1000,
		Unemployment:            baseMoney * rule.Unemployment / 1000,
		IndustrialInjury:        baseMoney * rule.IndustrialInjury / 1000,
		HousingAccumulationFund: baseMoney * rule.HousingAccumulationFund / 1000,
	}
	expenditure.Total = expenditure.Endowment + expenditure.Medical + expenditure.Birth + expenditure.Unemployment + expenditure.IndustrialInjury + expenditure.HousingAccumulationFund
	return expenditure
}

func getBaseMoney(salary int, rule Rule) int {
	baseMoney := salary
	if baseMoney > rule.BaseLimit {
		baseMoney = rule.BaseLimit
	}
	return baseMoney
}
