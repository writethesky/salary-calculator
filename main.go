package main

import (
	"fmt"
	"log"

	"github.com/writethesky/cmd"
)

type Params struct {
	Salary int `name:"s" usage:"月薪" require:"true" type:"int"`
}

type Rule struct {
	Endowment               int // 养老保险比例
	Medical                 int //医疗保险比例
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

	fmt.Printf("%+v\n", calculator(params.Salary, personRule))

}

type Calculate struct {
	Endowment               int // 养老保险
	Medical                 int //医疗保险
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
