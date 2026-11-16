package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

func Step(text string) {
	utils.LogWithLabelAndTimeStamp("Step", text)
}

func Check(text string) {
	utils.LogWithLabelAndTimeStamp("Check", text)
}

func Postcondition(text string) {
	utils.LogWithLabelAndTimeStamp("Postcondition", text)
}
