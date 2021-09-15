package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/writethesky/cmd"
)

type Params struct {
	Salary int `name:"s" usage:"月薪" require:"true" type:"int"`
}

type Rule struct {
	Endowment               int // 养老保险比例
	Medical                 int //医疗保险比例
	MedicalInput            int //医疗保险进入个人账户比例
	MedicalPlus             int //大病保险附加
	Birth                   int //生育保险比例
	Unemployment            int //失业保险比例
	IndustrialInjury        int //工伤保险比例
	HousingAccumulationFund int //住房公积金比例
	FreeTax                 int //每年免税额度
}

func main() {
	cmd := cmd.NewCMD("薪资计算器", "计算中")
	params := new(Params)
	if !cmd.Parse(params) {
		log.Fatalln("参数解析失败")
	}

	personRule := Rule{
		Endowment:               80,
		Medical:                 20,
		MedicalPlus:             3,
		Birth:                   0,
		Unemployment:            2,
		IndustrialInjury:        0,
		HousingAccumulationFund: 120,
		FreeTax:                 60000,
	}

	companyRule := Rule{
		Endowment:               200,
		Medical:                 100,
		MedicalInput:            28,
		MedicalPlus:             0,
		Birth:                   8,
		Unemployment:            10,
		IndustrialInjury:        10,
		HousingAccumulationFund: 120,
		FreeTax:                 0,
	}

	personResult := calculator(params.Salary, personRule)
	companyResult := calculator(params.Salary, companyRule)
	cmd.StopLoading()
	fmt.Println()
	printResult(params.Salary, personResult, companyResult)

}

func printResult(salary int, personResult, companyResult Calculate) {
	fmt.Println("五险一金")
	fmt.Printf("%s %s %s %s\n", getStringOfLength("", 20), getStringOfLength("个人", 10), getStringOfLength("企业", 10), getStringOfLength("总计", 10))

	fmt.Printf("%s %s %s %d\n", getStringOfLength("养老保险", 20), getStringOfLength(personResult.Endowment, 10), getStringOfLength(companyResult.Endowment, 10), personResult.Endowment+companyResult.Endowment)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("医疗保险", 20), getStringOfLength(personResult.Medical, 10), getStringOfLength(companyResult.Medical, 10), personResult.Medical+companyResult.Medical)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("生育保险", 20), getStringOfLength(personResult.Birth, 10), getStringOfLength(companyResult.Birth, 10), personResult.Birth+companyResult.Birth)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("失业保险", 20), getStringOfLength(personResult.Unemployment, 10), getStringOfLength(companyResult.Unemployment, 10), personResult.Unemployment+companyResult.Unemployment)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("工伤保险", 20), getStringOfLength(personResult.IndustrialInjury, 10), getStringOfLength(companyResult.IndustrialInjury, 10), personResult.IndustrialInjury+companyResult.IndustrialInjury)
	fmt.Printf("%s %s %s %d\n", getStringOfLength("住房公积金", 20), getStringOfLength(personResult.HousingAccumulationFund, 10), getStringOfLength(companyResult.HousingAccumulationFund, 10), personResult.HousingAccumulationFund+companyResult.HousingAccumulationFund)
	fmt.Println("\n税后工资")
	fmt.Printf("%s %d\n", getStringOfLength("五险一金后工资", 20), personResult.InsuranceSalary)
	fmt.Printf("%s %d\n", getStringOfLength("应缴纳个税", 20), personResult.Tax)
	fmt.Printf("%s %d\n", getStringOfLength("税后工资", 20), personResult.TaxSalary)
	fmt.Println("\n实际收入")
	fmt.Printf("%s %d\n", getStringOfLength("医保收入", 20), companyResult.MedicalInput)
	fmt.Printf("%s %d\n", getStringOfLength("住房公积金收入", 20), personResult.HousingAccumulationFund+companyResult.HousingAccumulationFund)
	fmt.Printf("%s %d\n", getStringOfLength("实际收入", 20), personResult.HousingAccumulationFund+companyResult.HousingAccumulationFund+personResult.TaxSalary+companyResult.MedicalInput)

	fmt.Println("\n企业成本")

	fmt.Printf("%s %d\n", getStringOfLength("员工工资", 20), salary)
	fmt.Printf("%s %d\n", getStringOfLength("五险一金", 20), companyResult.Endowment+companyResult.Medical+companyResult.Birth+companyResult.Unemployment+companyResult.IndustrialInjury+companyResult.HousingAccumulationFund)
	fmt.Printf("%s %d\n", getStringOfLength("总计", 20), salary+companyResult.Endowment+companyResult.Medical+companyResult.Birth+companyResult.Unemployment+companyResult.IndustrialInjury+companyResult.HousingAccumulationFund)
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

type Calculate struct {
	Endowment               int // 养老保险
	Medical                 int //医疗保险
	MedicalInput            int //医疗保险进入个人账户
	Birth                   int //生育保险
	Unemployment            int //失业保险
	IndustrialInjury        int //工伤保险
	HousingAccumulationFund int //住房公积金
	InsuranceSalary         int //五险一金后工资
	Tax                     int // 应缴纳个税
	TaxSalary               int //税后工资
}

type Tax struct {
	From int
	To   int
	Rate int
}

func calculator(salary int, rule Rule) Calculate {
	calculate := Calculate{
		Endowment:               salary * rule.Endowment / 1000,
		Medical:                 salary*rule.Medical/1000 + rule.MedicalPlus,
		MedicalInput:            salary * rule.MedicalInput / 1000,
		Birth:                   salary * rule.Birth / 1000,
		Unemployment:            salary * rule.Unemployment / 1000,
		IndustrialInjury:        salary * rule.IndustrialInjury / 1000,
		HousingAccumulationFund: salary * rule.HousingAccumulationFund / 1000,
	}
	calculate.InsuranceSalary = salary - calculate.Endowment - calculate.Medical - calculate.Birth - calculate.Unemployment - calculate.IndustrialInjury - calculate.HousingAccumulationFund

	taxes := getTaxTable()
	base := calculate.InsuranceSalary*12 - rule.FreeTax

	for i := len(taxes) - 1; i >= 0; i-- {
		if base > taxes[i].From {

			base -= taxes[i].From
			calculate.Tax += base * taxes[i].Rate / 1000
		}
	}

	calculate.Tax /= 12
	calculate.TaxSalary = calculate.InsuranceSalary - calculate.Tax

	return calculate
}

func getTaxTable() []Tax {
	return []Tax{
		{
			From: 0,
			To:   36000,
			Rate: 30,
		},
		{
			From: 36000,
			To:   144000,
			Rate: 100,
		},
		{
			From: 144000,
			To:   300000,
			Rate: 200,
		},
		{
			From: 300000,
			To:   420000,
			Rate: 250,
		},
		{
			From: 420000,
			To:   660000,
			Rate: 300,
		},
		{
			From: 660000,
			To:   960000,
			Rate: 350,
		},
		{
			From: 960000,
			To:   96000000000000000,
			Rate: 450,
		},
	}
}
