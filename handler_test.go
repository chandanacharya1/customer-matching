package main_test

import (
	"fmt"
	"github.com/chandanacharya1/customer-matching/middleware"
	"github.com/chandanacharya1/customer-matching/models"
	"testing"
)

func TestPartnerListempty(t *testing.T) {
	var customer models.Customer
	customer.Material = "rubber"
	customer.AddressLat = 20.50871548742554
	customer.AddressLong = 12.375536452185349
	customer.SquareMeters = 800
	customer.PhoneNumber = 125478956

	partners := middleware.GetMatchingPartners(customer)

	if len(partners) != 0 {
		t.Fatalf("Fail: the list should be empty")
	}

}

func TestOrderOfMatchedList(t *testing.T) {
	var customer models.Customer
	customer.Material = "tiles"
	customer.AddressLat = 52.50879681532554
	customer.AddressLong = 13.375567271135349
	customer.SquareMeters = 800
	customer.PhoneNumber = 125478956

	partners := middleware.GetMatchingPartners(customer)
	if len(partners) != 3 {
		failString := fmt.Sprintln("Fail: the list should be 3 and is", len(partners))
		t.Fatalf(failString)
	}

	if partners[0].PartnerID != 6 {
		failString := fmt.Sprintln("Fail: the first partner id in the list should be 6 and is", partners[0].PartnerID)
		t.Fatalf(failString)
	}

	if partners[1].PartnerID != 3 {
		failString := fmt.Sprintln("Fail: the second partner id in the list should be 3 and is", partners[1].PartnerID)
		t.Fatalf(failString)
	}
	if partners[2].PartnerID != 2 {
		failString := fmt.Sprintln("Fail: the third partner id in the list should be 2 and is", partners[2].PartnerID)
		t.Fatalf(failString)
	}

}

func TestGetAllPartners(t *testing.T) {
	partners := middleware.GetAllPartners()
	if len(partners) != 7 {
		failString := fmt.Sprintln("Partner List not correct, should have 7 and has ", len(partners))
		t.Fatalf(failString)
	}

}

func TestGetPartnerDetail(t *testing.T) {
	partner := middleware.GetPartnerFromId(2)
	if len(partner.PartnerName) == 0 {
		failString := fmt.Sprintln("Partner 552 not found")
		t.Fatalf(failString)
	}

	if partner.PartnerName != "Home24" {
		failString := fmt.Sprintln("Partner Name not correct, should be 'Mr.Carpet' and is ", partner.PartnerName)
		t.Fatalf(failString)
	}

}
