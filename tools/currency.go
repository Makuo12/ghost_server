package tools

//"strconv"

//"github.com/bojanz/currency"

//func NairaToDollar(rate string, naira string) (string, error) {
//	money, err := strconv.ParseFloat(naira, 64)

//	if err != nil {
//		log.Printf("There an error in NairaToDollar: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	rateMoney, err := strconv.ParseFloat(rate, 64)

//	if err != nil {
//		log.Printf("There an error in NairaToDollar: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	result := ConvertInt64ToString(int64(money / rateMoney))
//	return result, nil
//}

//func DollarToNaira(rate string, dollar string)(string, error){

//	money, err := strconv.ParseFloat(dollar, 64)

//	if err != nil {
//		log.Printf("There an error in DollarToNaira: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	rateMoney, err := strconv.ParseFloat(rate, 64)

//	if err != nil {
//		log.Printf("There an error in DollarToNaira: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	result := ConvertInt64ToString(int64(money*rateMoney))
//	return result, nil
//}

//func CADToDollar(rate string, CAD string) (string, error) {
//	money, err := strconv.ParseFloat(CAD, 64)

//	if err != nil {
//		log.Printf("There an error in CADToDollar: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	rateMoney, err := strconv.ParseFloat(rate, 64)
//	if err != nil {
//		log.Printf("There an error in CADToDollar: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	result := ConvertInt64ToString(int64(money / rateMoney))
//	return result, nil
//}

//func DollarToCAD(rate string, dollar string) (string, error) {
//	money, err := strconv.ParseFloat(dollar, 64)

//	if err != nil {
//		log.Printf("There an error in DollarToCAD: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	rateMoney, err := strconv.ParseFloat(rate, 64)

//	if err != nil {
//		log.Printf("There an error in DollarToCAD: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	result := ConvertInt64ToString(int64(money * rateMoney))
//	return result, nil
//}

//func NairaToCAD(dollarToNairaRate string, naira string, dollarToCADRate string) (string, error) {
//	// first we convert it to dollar
//	result, err := NairaToDollar(dollarToNairaRate, naira)

//	if err != nil {
//		log.Printf("There an error in nairaToCAD: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	resultTwo, err := DollarToCAD(dollarToCADRate, result)
//	return resultTwo, err
//}

//func CADToNaira(CADToDollarRate string, CAD string, dollarToNairaRate string) (string, error) {
//	// first we convert it to dollar
//	result, err := CADToDollar(CADToDollarRate, CAD)

//	if err != nil {
//		log.Printf("There an error in CADToNaira: %v\n", err)
//		err = fmt.Errorf("there was an internal server error")
//		return "", err
//	}
//	resultTwo, err := DollarToNaira(dollarToNairaRate, result)
//	return resultTwo, err
//}
