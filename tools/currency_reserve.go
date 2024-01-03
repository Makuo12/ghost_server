package tools

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
)

func NairaToDollarFloat(rate string, naira string) (float64, error) {
	money, err := strconv.ParseFloat(naira, 64)

	if err != nil {
		log.Printf("There an error in NairaToDollarFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	rateMoney, err := strconv.ParseFloat(rate, 64)

	if err != nil {
		log.Printf("There an error in NairaToDollarFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	result := money / rateMoney
	return result, nil
}

func DollarToNairaFloat(rate string, dollar string) (float64, error) {

	money, err := strconv.ParseFloat(dollar, 64)

	if err != nil {
		log.Printf("There an error in DollarToNairaFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	rateMoney, err := strconv.ParseFloat(rate, 64)

	if err != nil {
		log.Printf("There an error in DollarToNairaFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	result := money * rateMoney
	return result, nil
}

func CADToDollarFloat(rate string, CAD string) (float64, error) {
	money, err := strconv.ParseFloat(CAD, 64)

	if err != nil {
		log.Printf("There an error in CADToDollarFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	rateMoney, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		log.Printf("There an error in CADToDollarFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	result := money / rateMoney
	return result, nil
}

func DollarToCADFloat(rate string, dollar string) (float64, error) {
	money, err := strconv.ParseFloat(dollar, 64)

	if err != nil {
		log.Printf("There an error in DollarToCADFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	rateMoney, err := strconv.ParseFloat(rate, 64)

	if err != nil {
		log.Printf("There an error in DollarToCADFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	result := money * rateMoney
	return result, nil
}

func NairaToCADFloat(dollarToNairaRate string, naira string, dollarToCADRate string) (float64, error) {
	// first we convert it to dollar
	result, err := NairaToDollarFloat(dollarToNairaRate, naira)

	if err != nil {
		log.Printf("There an error in nairaToCADFloat: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	resultTwo, err := DollarToCADFloat(dollarToCADRate, ConvertFloatToString(result))
	return resultTwo, err
}

func CADToNairaFloat(CADToDollarRate string, CAD string, dollarToNairaRate string) (float64, error) {
	// first we convert it to dollar
	result, err := CADToDollarFloat(CADToDollarRate, CAD)

	if err != nil {
		log.Printf("There an error in CADToNaira: %v\n", err)
		err = fmt.Errorf("there was an internal server error")
		return 0, err
	}
	resultTwo, err := DollarToNairaFloat(dollarToNairaRate, ConvertFloatToString(result))
	return resultTwo, err
}

func ConvertPrice(price string, hostCurrency string, userCurrency string, dollarToNaira string, dollarToCAD string, ID uuid.UUID) (priceConvert float64, err error) {
	if price == "0.00" || price == "0" || ServerStringEmpty(price) {
		priceConvert = 0.00
		return
	}
	if hostCurrency == "USD" {
		switch userCurrency {
		case "USD":
			// USD TO USD
			priceConvert = ConvertStringToFloat(price)
			return
		case "NGN":
			// USD TO NGN
			priceConvert, err = DollarToNairaFloat(dollarToNaira, price)
			if err != nil {
				log.Printf("There an error in USD TO NGN: %v, ID: %v\n", err, ID)
			}
		case "CAD":
			priceConvert, err = DollarToCADFloat(dollarToCAD, price)
			if err != nil {
				log.Printf("There an error in USD TO CAD: %v, ID: %v\n", err, ID)
			}
		}
	} else if hostCurrency == "NGN" {
		switch userCurrency {
		case "NGN":
			// NGN TO NGN
			priceConvert = ConvertStringToFloat(price)

		case "USD":
			// NGN TO USD
			priceConvert, err = NairaToDollarFloat(dollarToNaira, price)
			if err != nil {
				log.Printf("There an error in NGN TO USD: %v, ID: %v\n", err, ID)
			}
		case "CAD":
			// NGN TO CAD
			priceConvert, err = NairaToCADFloat(dollarToNaira, price, dollarToCAD)
			if err != nil {
				log.Printf("There an error in NGN TO CAD: %v, ID: %v\n", err, ID)
			}
		}
	} else if hostCurrency == "CAD" {
		switch userCurrency {
		case "CAD":
			// CAD TO CAD
			priceConvert = ConvertStringToFloat(price)

		case "USD":
			// CAD TO USD
			priceConvert, err = CADToDollarFloat(dollarToCAD, price)
			if err != nil {
				log.Printf("There an error in CAD TO USD: %v, ID: %v\n", err, ID)
			}
		case "NGN":
			// CAD TO NGN
			priceConvert, err = CADToNairaFloat(dollarToCAD, price, dollarToNaira)
			if err != nil {
				log.Printf("There an error in CAD TO NGN: %v, ID: %v\n", err, ID)
			}
		}
	}else{
		err = fmt.Errorf("currency not found")
	}
	return
}
