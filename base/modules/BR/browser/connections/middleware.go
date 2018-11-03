package connections

/*Middleware : middleware*/
func Middleware(C1 chan []byte) chan []byte {
	C2 := make(chan []byte, 2048)
	if data, ok := <-C1; ok {
		C2 <- data
	}
	return C2
}
