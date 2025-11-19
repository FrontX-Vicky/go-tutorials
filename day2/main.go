package main

import "fmt"

// TODO: Define the User struct
type User struct {
	Name string 
	Age int 
	Skills []string
}



// TODO: Add methods to User (Greet, IsAdult, AddSkill)

func (u User) Greet() string {
	return "Hello, " + u.Name
}

func (u User) IsAdult() bool {
	return u.Age >= 18
}

func (u *User) AddSkill(skill string){
	u.Skills = append(u.Skills, skill) 
}


// TODO: Define the Profile interface

type Profile interface {
	PrintProfile() 
}

// TODO: Implement Profile for User
func (u User) PrintProfile()  {
	fmt.Printf("Name: %s, Age: %d, Skills: %v\n", u.Name, u.Age, u.Skills)
}

// TODO (Optional): Define the Employee struct and override PrintProfile

type Employee struct {
	User
	Position string
}

func (e Employee) PrintProfile()  {
	fmt.Printf("Name: %s, Age: %d, Skills: %v, Position: %s\n", e.Name, e.Age, e.Skills, e.Position)
}

func PrintProfiles(profiles []Profile) {
	for _, profile := range profiles{ 
		profile.PrintProfile()
	}
}
func main() {
	// TODO: Create a User instance and test methods
	user := User{Name: "Alice", Age: 30, Skills: []string{"Go", "Python"}}
	fmt.Println(user.Greet())
	fmt.Println(user.IsAdult())
	user.AddSkill("java")


	fmt.Println(user.Skills)
	// fmt.Println(user.PrintProfile())
	// TODO: (Optional) Create an Employee instance and test methods
	employee := Employee{
		User : User{Name: "Bob", Age: 25, Skills: []string{"Java", "C++"}},
		Position: "Developer",
	}
	// fmt.Println(employee.PrintProfile())
	
	profiles := []Profile{user, employee}
	PrintProfiles(profiles)

	fmt.Println("Day 2 tasks completed!")
}
