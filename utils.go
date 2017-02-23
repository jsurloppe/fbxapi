package fbxapi

import (
	"log"
	"os"
	"strconv"
)

var Logr *log.Logger

func init() {
	output := os.Stderr
	flags := log.Flags()
	prefix := log.Prefix()
	Logr = log.New(output, prefix, flags)
}

func drebug(format string, args ...interface{}) {
	Logr.Printf(format, args)
}

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
