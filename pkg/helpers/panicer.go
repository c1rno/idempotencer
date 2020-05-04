package helpers

func Panicer(err error) {
	if err != nil {
		panic(err)
	}
}
