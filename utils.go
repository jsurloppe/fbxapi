package fbxapi

import (
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicAttack(err *error) {
	if r := recover(); r != nil {
		if safeErr, ok := r.(error); ok {
			*err = safeErr
		} else {
			// call medic!
			panic(r)
		}
	}
}

func recoverAsErr(err error) {

}

func dataIsNil(data interface{}) bool {
	return data == nil
}

func APIVersionToInt(apiVersion string) (int, error) {
	f, err := strconv.ParseFloat(apiVersion, 32)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}
