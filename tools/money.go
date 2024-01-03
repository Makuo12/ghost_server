package tools

import (
	"fmt"
	"math"
	"strconv"
)

func GetBalanceFrom(balance string, amount string) (string, error) {
	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return "", err
	}
	am, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	result := bal - am
	resultString := strconv.FormatFloat(result, 'f', -1, 64)
	return resultString, err
}

func GetBalanceTo(balance string, amount string) (string, error) {
	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return "", err
	}
	am, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	result := bal + am
	resultString := strconv.FormatFloat(result, 'f', -1, 64)
	return resultString, err
}

func PaystackMoneyToDB(amount int) (amountString string) {
	// Divide by 100 when storing in database
	amountFloat := float64(amount) / 100
	amountString = ConvertFloatToString(amountFloat)
	return
}

// MoneyStringToInt this converts money to its lowest form either in kobo or cents
func MoneyStringToInt(moneyStr string) int64 {
	// Convert the string to a float64
	moneyFloat, err := strconv.ParseFloat(moneyStr, 64)
	if err != nil {
		return 0
	}
	// Multiply by 100 to convert to cents and cast to int
	cents := int64(moneyFloat * 100)
	return cents
}

// MoneyStringToInt this converts money to its lowest form either in kobo or cents
func MoneyFloatToInt(moneyFloat float64) int64 {
	// Multiply by 100 to convert to cents and cast to int
	cents := int64(moneyFloat * 100)
	return cents
}

func IntToMoneyString(cents int64) string {
	// Divide by 100 to convert cents to dollars
	dollars := float64(cents) / 100.0

	// Format as a string with two decimal places
	return fmt.Sprintf("%.2f", dollars)
}

func RemoveGateCharge(amount int, fees int) (amountString string) {
	currentAmount := amount - fees
	amountFloat := float64(currentAmount) / 100
	amountString = ConvertFloatToString(amountFloat)
	return
}

func HandleDiscount(fee float64, percent int) (number float64) {
	return fee - (fee * (float64(percent) / 100))
}

func ConvertToPaystackCharge(charge string) int {
	paystackCharge := ConvertStringToFloat(charge)
	return int(math.Ceil(paystackCharge * 100))
}

func ConvertToPaystackPayout(payout string) int {
	paystackPayout := ConvertStringToFloat(payout)
	return int(math.Floor(paystackPayout * 100))
}

func ConvertToPaystackChargeString(charge string) string {
	return ConvertInt64ToString(int64(ConvertToPaystackCharge(charge)))
}

// Converts account number to only four characters
func AccountNumberToFour(a string) (s string) {
	if ServerStringEmpty(a) {
		s = a
		return
	}
	for i := 0; i < len(a); i++ {
		s += string(a[i])
		if i == 3 {
			return
		}
	}
	return
}
