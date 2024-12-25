package utils

import (
	"github.com/shopspring/decimal"
)

func calculateOverpowerBase(score int64, internalLevel decimal.Decimal) decimal.Decimal {
	levelBase := internalLevel.Mul(decimal.NewFromInt(10000)).BigInt().Int64()

	var op int64 = 0

	if score >= 1_007_500 {
		op = levelBase + 20_000 + (score-1_007_500)*3
	} else if score >= 1_005_000 {
		op = levelBase + 15_000 + (score-1_005_000)*2
	} else if score >= 1_000_000 {
		op = levelBase + 10_000 + (score - 1_000_000)
	} else if score >= 975_000 {
		op = levelBase + (score-975_000)*2/5
	} else if score >= 900_000 {
		op = levelBase - 50_000 + (score-900_000)*2/3
	} else if score >= 800_000 {
		op = (levelBase-50_000)/2 + (((score - 800_000) * ((levelBase - 50_000) / 2)) / 100_000)
	} else if score >= 500_000 {
		op = (((levelBase - 50_000) / 2) * (score - 500_000)) / 300_000
	}

	if score >= 975_000 {
		return decimal.NewFromInt(op).DivRound(decimal.NewFromInt(1_000), 2).Div(decimal.NewFromInt(2))
	}

	return decimal.NewFromInt(op).DivRound(decimal.NewFromInt(1_000), 2).Div(decimal.NewFromInt(5))
}

func calculateOverpowerMax(internalLevel decimal.Decimal) decimal.Decimal {
	return internalLevel.Mul(decimal.NewFromInt(5)).Add(decimal.NewFromInt(15))
}

func CalculateOverpower(score int64, internalLevel decimal.Decimal, lamp string) decimal.Decimal {
	playOp := calculateOverpowerBase(score, internalLevel)

	if score == 1_010_000 {
		playOp = calculateOverpowerMax(internalLevel)
	} else if lamp == "ALL JUSTICE" || lamp == "ALL JUSTICE CRITICAL" {
		playOp = playOp.Add(decimal.NewFromInt(1))
	} else if lamp == "FULL COMBO" {
		playOp = playOp.Add(decimal.NewFromFloat(0.5))
	}

	return playOp
}
