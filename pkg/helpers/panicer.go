package helpers

func Panicer(errs ...interface{}) {
	for _, err := range errs {
		terr, ok := err.(error)
		if ok && terr != nil {
			panic(terr)
		}
	}
}
