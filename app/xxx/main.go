package main

import "fmt"

func main() {
	defer func(x int) { fmt.Println(x) }(func() int { fmt.Println("build"); return 123 }())
	fmt.Println("main")
	//myApp := app.New()
	//myWindow := myApp.NewWindow("Form Layout")
	//
	//label1 := canvas.NewText("Label 1", color.Black)
	//value1 := canvas.NewText("Value", color.Black)
	//label2 := canvas.NewText("Label 2", color.Black)
	//value2 := canvas.NewText("Something", color.Black)
	//grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2)
	//myWindow.SetContent(grid)
	//myWindow.ShowAndRun()
}
