package main

import (
	"fmt"
	"learning_go/myutils"
	"strings"
	"sync"
)

type Customer struct{
	userTickets uint16
	firstName string
	lastName string
	email string
}


var wg = sync.WaitGroup()

func main() {
	myutils.Test()
	var totalTickets uint16
	var userTickets uint16
	var firstName string
	var lastName string
	var email string
	totalTickets = 20
	flightName := "Flight"
	var bookings = make([]Customer,0)
	
	wg.Add(1)
	sendTicket(userTickets,firstName,lastName)
	for {
		fmt.Printf("We have %d Tickets.\n", totalTickets)
		fmt.Printf("Type:%T\n", totalTickets)
		fmt.Printf("Flight Name is %s\n", flightName)
		fmt.Println("Current Booking ", bookings, ".")
		userTickets,firstName,lastName,email = getUserInput()
		if !(len(firstName) >= 2) && !(len(lastName) >= 2) && !(strings.Contains(email, "@")) {
			fmt.Println("First Name and Last Name Must Have More than 2 characters and Email Must be Valid.")
			continue
		}



		if userTickets > totalTickets {
			fmt.Printf("We only have %d tickets left. You cannot Book %d Tickets.", totalTickets, userTickets)
			continue
		}
		// bookings = append(bookings, firstName+" "+lastName)
		totalTickets -= userTickets
		var userData = Customer{
			firstName: firstName,
			lastName: lastName,
			email: email,
			userTickets: userTickets,
		}


		/*
		var userData = make(map[string]string)
		userData["firstName"] = firstName
		userData["lastName"] = lastName
		userData["email"] = email
		userData["userTickets"] = strconv.FormatUint(uint64(userTickets),10)
		*/
		bookings = append(bookings, userData)

		 firstNames := getFirstNames(bookings)
			if totalTickets == 0 {
			fmt.Println("Sorry All the tickets are Booked. You cannot Book More Tickets.")
			break
		}

		fmt.Printf("Username: %s\n", firstName+" "+lastName)
		fmt.Printf("User Tickets: %d\n", userTickets)
		fmt.Println("Current Booking ", firstNames, ".")

	}
	wg.Wait()
}

func getUserInput() (uint16,string,string,string) {
	var userTickets uint16
	var firstName string
	var lastName string
	var email string
	fmt.Print("Enter Your First Name:")
	fmt.Scan(&firstName)

	fmt.Print("Enter Your Last Name:")
	fmt.Scan(&lastName)

	fmt.Print("Enter You E-mail:")
	fmt.Scan(&email)

	fmt.Print("Enter Number of Tickets:")
	fmt.Scan(&userTickets)

	return userTickets,firstName,lastName,email
}

func getFirstNames(bookings []Customer) []string{


		firstNames := []string{}
		for _, booking := range bookings {
			firstNames = append(firstNames, booking.firstName)
		}

		return firstNames
}

func sendTicket(userTickents uint16,firstName string, lastname string){
var message = fmt.Snprintf("Your %d tickets has been Booked Mr.or Mrs.%s %s",userTickets,firstName,lastName)
fmt.Println("Information:",message)
}
