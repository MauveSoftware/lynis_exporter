package main

import "strconv"

type converter func(string) (float64, error)

func boolConverter() converter {
	return func(x string) (float64, error) {
		b, err := strconv.ParseBool(x)
		if err != nil {
			return 0, err
		}

		if b {
			return float64(1), nil
		}

		return float64(0), nil
	}
}

func floatConverter() converter {
	return func(x string) (float64, error) {
		return strconv.ParseFloat(x, 64)
	}
}

func converterByName(name string) converter {
	switch name {
	case "bool":
		return boolConverter()
	}

	return floatConverter()
}
