package iif

func String(expr bool, trueValue, falseValue string) string {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Int(expr bool, trueValue, falseValue int) int {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Int64(expr bool, trueValue, falseValue int64) int64 {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Uint(expr bool, trueValue, falseValue uint) uint {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Uint64(expr bool, trueValue, falseValue uint64) uint64 {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Float32(expr bool, trueValue, falseValue float32) float32 {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}

func Float64(expr bool, trueValue, falseValue float64) float64 {
	if expr {
		return trueValue
	} else {
		return falseValue
	}
}
