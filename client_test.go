package aidev

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	opts := defaultOptions()
	opts.User = "admin"
	opts.Password = "202cb962ac59075b964b07152d234b70"
	// opts.Debug = true
	cl := NewClient(opts)
	if err := cl.RenewToken(); err != nil {
		t.Fatal(err)
	}

	t.Run("GetToken", func(t *testing.T) {
		token, err := cl.BaseAPI.GetToken("admin", "202cb962ac59075b964b07152d234b70")
		if err != nil {
			t.Error(err)
		} else {
			log.Println(token)
		}
	})

	t.Run("GetDates", func(t *testing.T) {
		list, err := cl.BaseAPI.GetDates("2020-07-05")
		if err != nil {
			t.Error(err)
		}
		log.Printf("%+v", list)

		// Generate a sample QR code
		qr, err := QRCode(list[0].Token)
		if err != nil {
			t.Fatal(err)
		}

		// Save temp file
		_ = ioutil.WriteFile("docs/sample_qr.png", qr, 0400)
	})

	t.Run("GetCabins", func(t *testing.T) {
		list, err := cl.BaseAPI.GetCabins()
		if err != nil {
			t.Error(err)
		}
		log.Printf("%+v", list)
	})

	t.Run("GetStudies", func(t *testing.T) {
		res, err := cl.BaseAPI.GetStudies()
		if err != nil {
			t.Error(err)
		} else {
			log.Printf("%+v", res)
		}
	})

	t.Run("GetAvailableTime", func(t *testing.T) {
		res, err := cl.BaseAPI.GetAvailableTime("1", "2020-07-05")
		if err != nil {
			t.Error(err)
		}
		log.Printf("%+v", res)
	})

	t.Run("AddPerson", func(t *testing.T) {
		id, err := cl.BaseAPI.AddPerson(PersonInput{
			Name:   "Rick Sanchez",
			CURP:   "MOSM780130HVZNTR07",
			Age:    "21",
			Gender: "1",
		})
		if err != nil {
			t.Error(err)
		}
		log.Println(id)
	})

	t.Run("AddAppointment", func(t *testing.T) {
		res, err := cl.BaseAPI.AddAppointment(DateInput{
			CabinID:   "1",
			StudyID:   "1",
			PersonID:  "7",
			Date:      "2020-07-05",
			HourStart: "11:00:00",
			HourEnd:   "12:00:00",
		})
		if err != nil {
			t.Error(err)
		} else {
			log.Printf("%+v", res)
		}
	})

	t.Run("GetResults", func(t *testing.T) {
		res, err := cl.BaseAPI.GetResults(7)
		if err != nil {
			t.Error(err)
		} else {
			log.Printf("%+v", res)
		}
	})
}
