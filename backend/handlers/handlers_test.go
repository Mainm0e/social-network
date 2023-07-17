package handlers

import (
	"backend/db"
	"backend/utils"
	"errors"
	"fmt"
	"testing"
)

/*
createRandomUser generates random user data for testing purposes. It creates a user with random values for db.User struct fields (expect for UserId,nickname and avatar)
It returns the created user as a db.User struct and any error encountered during the process.
*/
func createRandomUser(t *testing.T) (RegisterData, error) {
	firstName, err := utils.RandomString(6)
	if err != nil {
		t.Fatalf("error in random name")
	}
	lastName, err := utils.RandomString(6)
	if err != nil {
		t.Fatalf("error in random name")
	}
	email, err := utils.GenerateRandomEmail()
	if err != nil {
		t.Fatalf("error in random email")
	}
	password, err := utils.RandomPassword(8)
	if err != nil {
		t.Fatalf("error in random password")
	}
	aboutMe, err := utils.RandomString(100)
	if err != nil {
		t.Fatalf("error in random about me")
	}

	user := RegisterData{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: utils.RandomDateBeforeNow(),
		Email:     email,
		Password:  password,
		AboutMe:   aboutMe,
	}
	return user, nil
}

/*
insertRandomUser generates a random user using createRandomUser and inserts it into the database for testing purposes.
It returns the inserted user as a db.User struct and any error encountered during the process.
*/
func insertRandomUser(t *testing.T) (RegisterData, error) {
	user, err := createRandomUser(t)
	fmt.Println(user)
	if err != nil {
		return user, errors.New("createRandomUser got error: " + err.Error())
	}
	data := RegisterData{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Email:     user.Email,
		Password:  user.Password,
		AboutMe:   user.AboutMe,
		Avatar:    user.Avatar,
		NickName:  user.NickName,
	}
	err = data.register()
	if err != nil {
		return user, errors.New("register got error: " + err.Error())
	}
	return user, nil
}

/*
TestLogin performs tests for the login function.using random user data created by createRandomUser and insertRandomUser,
to make sure functionality independent of the existed users data in the database.
*/
func TestLogin(t *testing.T) {
	db.Check("../db/database.db", "./backend/database/test.sql")

	user, err := insertRandomUser(t)
	if err != nil {
		t.Errorf("insertRandomUser got error: %v:", err)
	}
	randomEmail, err := utils.GenerateRandomEmail()
	if err != nil {
		t.Errorf("GenerateRandomEmail got error: %v:", err)
	}
	randomPassword, err := utils.RandomPassword(8)
	if err != nil {
		t.Errorf("RandomPassword got error: %v:", err)
	}
	var tests = []struct {
		Email, Password any
	}{
		//test correct input
		{user.Email, user.Password},
		//test incorrect input (wrong password)
		{user.Email, randomPassword},
		//test incorrect email (user not found)
		{randomEmail, user.Password},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.Email, tt.Password)
		t.Run(testname, func(t *testing.T) {
			var data LoginData
			data.Email = tt.Email.(string)
			data.Password = tt.Password.(string)

			_, err := data.login()

			if err != nil {
				t.Errorf("login failed got error: %v", err)
			}
		})

	}

}

/*
TestRegister performs tests for the register function.using random user data created by createRandomUser and insertRandomUser,
to make sure functionality independent of the existed users data in the database.
*/
func TestRegister(t *testing.T) {
	db.Check("../db/database.db", "./backend/database/test.sql")
	user, err := insertRandomUser(t)
	if err != nil {
		t.Errorf("insertRandomUser got error: %v:", err)
	}
	randomUserInfo, err := createRandomUser(t)
	if err != nil {
		t.Errorf("createRandomUser got error: %v:", err)
	}
	var tests = []struct {
		user RegisterData
		want error
	}{
		//test correct input
		{randomUserInfo, nil},
		//test incorrect input (user already exists)
		//want not nil error
		{user, errors.New("user with this email already exists")},
	}
	for _, tt := range tests {

		testName := fmt.Sprintf("email: %s, password:%s", tt.user.Email, tt.user.Password)
		t.Run(testName, func(t *testing.T) {
			// marshal the email and password into json base on LoginData struct and some random data
			data := RegisterData{
				FirstName: tt.user.FirstName,
				LastName:  tt.user.LastName,
				BirthDate: tt.user.BirthDate,
				Email:     tt.user.Email,
				Password:  tt.user.Password,
				AboutMe:   tt.user.AboutMe,
				Avatar:    tt.user.Avatar,
				NickName:  tt.user.NickName,
			}
			err := data.register()

			if err != nil {
				t.Errorf("register failed got error: %v, want: %v ", err, tt.want)
			}
		})

	}
}
