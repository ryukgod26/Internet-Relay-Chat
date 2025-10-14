package main

import (
	"fmt"
	"strings"
)

func main() {
	var totalTickets uint16
	var userTickets uint16
	var firstName string
	var lastName string
	var email string
	totalTickets = 20
	flightName := "Flight"
	var bookings = []string{"Armin", "Eren", "Mikasa"}
	fmt.Printf("We have %d Tickets.\n", totalTickets)
	fmt.Printf("Type:%T\n", totalTickets)
	fmt.Printf("Flight Name is %s\n", flightName)
	fmt.Println("Current Booking ", bookings, ".")
	for {
		fmt.Print("Enter Your First Name:")
		fmt.Scan(&firstName)

		fmt.Print("Enter Your Last Name:")
		fmt.Scan(&lastName)

		fmt.Print("Enter You E-mail:")
		fmt.Scan(&email)

		if !(len(firstName) >= 2) && !(len(lastName) >= 2) && !(strings.Contains(email, "@")) {
			fmt.Println("First Name and Last Name Must Have More than 2 characters and Email Must be Valid.")
			continue
		}

		fmt.Print("Enter Number of Tickets:")
		fmt.Scan(&userTickets)

		if userTickets > totalTickets {
			fmt.Printf("We only have %d tickets left. You cannot Book %d Tickets.", totalTickets, userTickets)
			continue
		}

		totalTickets -= userTickets

		bookings = append(bookings, firstName+" "+lastName)

		firstNames := []string{}
		for _, booking := range bookings {
			names := strings.Fields(booking)
			firstNames = append(firstNames, names[0])
		}
		if totalTickets == 0 {
			fmt.Println("Sorry All the tickets are Booked. You cannot Book More Tickets.")
			break
		}

		fmt.Printf("Username: %s\n", firstName+" "+lastName)
		fmt.Printf("User Tickets: %d\n", userTickets)
		fmt.Println("Current Booking ", firstNames, ".")
	}
}
