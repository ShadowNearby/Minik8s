package controller

func Run() {
	var serviceController ServiceController
	go StartController(&serviceController)
	select {}
}
