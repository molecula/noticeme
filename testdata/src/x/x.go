package x

func x() int {
	return 0
}

func directCall() {
	x() // want `unused value of type int`
}

func directCallUsing() int {
	return x()
}

var xMap = map[string]func() int{"foo": x}

func unusedMap() {
	xMap["foo"]() // want `unused value of type int`
}

func usedMap() int {
	return xMap["foo"]()
}
