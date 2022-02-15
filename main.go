package main

import (
	"fmt"
	"log"
	"salary-calculator/calculator"
	"strings"

	"github.com/writethesky/cmd"
)

type Params struct {
	Salary int `name:"s" usage:"月薪" require:"true" type:"int"`
}

func main() {
	cmd := cmd.NewCMD("薪资计算器", "计算中")
	params := new(Params)
	if !cmd.Parse(params) {
		log.Fatalln("参数解析失败")
	}

	personRule := calculator.Rule{
		Endowment:               80,
		Medical:                 20,
		MedicalInput:            20,
		MedicalPlus:             3,
		Birth:                   0,
		Unemployment:            5,
		IndustrialInjury:        0,
		HousingAccumulationFund: 120,
		BaseLimit:               28221,
		FreeTax:                 60000,
	}

	companyRule := calculator.Rule{
		Endowment:               200,
		Medical:                 100,
		MedicalInput:            8,
		MedicalPlus:             0,
		Birth:                   8,
		Unemployment:            5,
		IndustrialInjury:        10,
		HousingAccumulationFund: 120,
		BaseLimit:               28221,
		FreeTax:                 0,
	}

	result := calculator.Calculate(params.Salary, personRule, companyRule)
	cmd.StopLoading()
	fmt.Println()
	printResult(params.Salary, result)

}

func printResult(salary int, result calculator.Result) {
	fmt.Println("五险一金\n-------------------------------------------------")
	fmt.Printf("%s %s %s %s\n", getStringOfLength("", 20), getStringOfLength("个人", 10), getStringOfLength("企业", 10), getStringOfLength("总计", 10))

	fmt.Printf("%s %s %s %d\n", getStringOfLength("养老保险", 20), getStringOfLength(result.Person.Endowment, 10), getStringOfLength(result.Company.Endowment, 10), result.Person.Endowment+result.Company.Endowment)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("医疗保险", 20), getStringOfLength(result.Person.Medical, 10), getStringOfLength(result.Company.Medical, 10), result.Person.Medical+result.Company.Medical)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("生育保险", 20), getStringOfLength(result.Person.Birth, 10), getStringOfLength(result.Company.Birth, 10), result.Person.Birth+result.Company.Birth)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("失业保险", 20), getStringOfLength(result.Person.Unemployment, 10), getStringOfLength(result.Company.Unemployment, 10), result.Person.Unemployment+result.Company.Unemployment)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("工伤保险", 20), getStringOfLength(result.Person.IndustrialInjury, 10), getStringOfLength(result.Company.IndustrialInjury, 10), result.Person.IndustrialInjury+result.Company.IndustrialInjury)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("住房公积金", 20), getStringOfLength(result.Person.HousingAccumulationFund, 10), getStringOfLength(result.Company.HousingAccumulationFund, 10), result.Person.HousingAccumulationFund+result.Company.HousingAccumulationFund)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("总计", 20),
		getStringOfLength(result.Person.Expenditure.Total, 10),
		getStringOfLength(result.Company.Expenditure.Total, 10), 0)
	fmt.Println("\n税后工资\n---------------------------")
	fmt.Printf("%s %d\n", getStringOfLength("五险一金后工资", 20), result.Person.InsuranceSalary)
	fmt.Printf("%s %d\n", getStringOfLength("应缴纳个税(月均)", 20), result.Person.MonthlyAverageTax)
	fmt.Printf("%s %d\n", getStringOfLength("应缴纳个税(每月)", 20), result.Person.MonthlyTax)
	fmt.Printf("%s %d\n", getStringOfLength("税后工资(月均)", 20), result.Person.MonthlyAverageTaxSalary)
	fmt.Printf("%s %d\n", getStringOfLength("税后工资(每月)", 20), result.Person.MonthlyTaxSalary)
	fmt.Println("\n实际收入\n---------------------------")
	fmt.Printf("%s %d\n", getStringOfLength("医保收入", 20), result.Person.MedicalInput)
	fmt.Printf("%s %d\n", getStringOfLength("住房公积金收入", 20), result.Person.HousingAccumulationFund+result.Company.HousingAccumulationFund)
	fmt.Printf("%s %d\n", getStringOfLength("实际收入(月均)", 20), result.Person.RealIncome)

	fmt.Println("\n企业成本\n---------------------------")

	fmt.Printf("%s %d\n", getStringOfLength("员工工资", 20), salary)
	fmt.Printf("%s %d\n", getStringOfLength("五险一金", 20), result.Company.Expenditure.Total)
	fmt.Printf("%s %d\n", getStringOfLength("总计", 20), salary+result.Company.Expenditure.Total)
}

func getStringOfLength(any interface{}, length int) string {
	str := fmt.Sprintf("%v", any)

	englishLength := (3*len([]rune(str)) - len([]byte(str))) / 2
	chineseLength := len([]rune(str)) - englishLength

	originLength := chineseLength*2 + englishLength

	repeatCount := length - originLength

	if repeatCount < 0 {
		panic(fmt.Sprintf("getStringOfLength(%s, %d) length太小了", str, length))
	}
	return str + strings.Repeat(" ", repeatCount)
}
