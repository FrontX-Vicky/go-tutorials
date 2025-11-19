package main

import (
	"fmt"
	"slices"
)

func multiply(a, b int) (int, error) {
	if a <= 0 || b <= 0 {
		return 0, fmt.Errorf("a and b must be > 0")
	}
	return a * b, nil
}

func main() {
	name := "vicky"
	age := 25
	weight := 75.43
	married := false

	languages := []string{"python", "go", "PHP"}

	personal := map[string]string{
		"dob":     "14/02/1999",
		"surname": "more",
	}

	fmt.Println(name, age, weight, married)
	fmt.Println("Languages:", languages)
	fmt.Println("Personal:", personal)

	result, err := multiply(9, 9)
	if err != nil {
		fmt.Println("multiply error:", err)
		return
	}
	fmt.Println("Multiply result:", result)

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println("Numbers:", numbers)

	numbers = append(numbers, 6)
	fmt.Println("After append:", numbers)

	numbers = slices.Delete(numbers, 2, 3) // remove index 2 (3rd element)
	fmt.Println("After delete:", numbers)

	student := map[string]string{
		"name":  "vijay",
		"age":   "25",
		"grade": "8th",
		"city":  "mumbai",
	}
	fmt.Println("Student:", student)

	student["city"] = "delhi"
	fmt.Println("Updated student:", student)

	var city string

	fmt.Println("Enter your name, age and city (space separated):")
	_, err = fmt.Scan(&name, &age, &city)
	if err != nil {
		fmt.Println("input error:", err)
		return
	}

	sentence := fmt.Sprintf("Hello %s, age %d from %s. Welcome to Go!", name, age, city)
	fmt.Println(sentence)
}
