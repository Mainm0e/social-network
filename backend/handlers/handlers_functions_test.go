package handlers

import (
	"backend/db"
	"backend/util"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

/*
createRandomUser generates random user data for testing purposes. It creates a user with random values for db.User struct fields (expect for UserId,nickname and avatar)
It returns the created user as a db.User struct and any error encountered during the process.
*/
func createRandomUser(t *testing.T) (db.User, error) {
	firstName, err := util.RandomString(6)
	if err != nil {
		t.Fatalf("error in random name")
	}
	lastName, err := util.RandomString(6)
	if err != nil {
		t.Fatalf("error in random name")
	}
	email, err := util.GenerateRandomEmail()
	if err != nil {
		t.Fatalf("error in random email")
	}
	password, err := util.RandomPassword(8)
	if err != nil {
		t.Fatalf("error in random password")
	}
	aboutMe, err := util.RandomString(100)
	if err != nil {
		t.Fatalf("error in random about me")
	}

	user := db.User{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: util.RandomDateBeforeNow(),
		Email:     email,
		Password:  password,
		AboutMe: &sql.NullString{
			String: aboutMe,
			Valid:  true,
		},
		CreationTime: util.RandomDateTimeAfterNow(),
		Privacy:      "public",
	}
	return user, nil
}

/*
insertRandomUser generates a random user using createRandomUser and inserts it into the database for testing purposes.
It returns the inserted user as a db.User struct and any error encountered during the process.
*/
func insertRandomUser(t *testing.T) (db.User, error) {
	user, err := createRandomUser(t)
	if err != nil {
		return user, errors.New("createRandomUser got error: " + err.Error())
	}
	//marshal the user into json
	data, err := json.Marshal(user)
	if err != nil {
		return user, errors.New("marshaling got error: %v: " + err.Error())
	}
	_, err = register(data)
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
	randomEmail, err := util.GenerateRandomEmail()
	if err != nil {
		t.Errorf("GenerateRandomEmail got error: %v:", err)
	}
	randomPassword, err := util.RandomPassword(8)
	if err != nil {
		t.Errorf("RandomPassword got error: %v:", err)
	}
	var tests = []struct {
		Email, Password any
		want            bool
	}{
		//test correct input
		{user.Email, user.Password, true},
		//test incorrect input (wrong password)
		{user.Email, randomPassword, false},
		//test incorrect email (user not found)
		{randomEmail, user.Password, false},
		//test incorrect input format
		{23234, randomPassword, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.Email, tt.Password)
		t.Run(testname, func(t *testing.T) {
			data, err := json.Marshal(tt)
			if err != nil {
				t.Errorf("marshaling got error: %v:", err)
			}
			bo, err := login(data)

			if bo != tt.want {
				t.Errorf("login got: %v, want: %v error:%v.", bo, tt.want, err)
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
		user db.User
		want bool
	}{
		//test correct input
		{randomUserInfo, true},
		//test incorrect input (user already exists)
		{user, false},
	}
	for _, tt := range tests {

		testName := fmt.Sprintf("email: %s, password:%s", tt.user.Email, tt.user.Password)
		t.Run(testName, func(t *testing.T) {
			// marshal the email and password into json base on LoginData struct and some random data
			data, err := json.Marshal(tt.user)
			if err != nil {
				t.Errorf("marshaling got error: %v:", err)
			}
			bo, err := register(data)

			if bo != tt.want {
				t.Errorf("register failed got: %v, want: %v error:%v.", bo, tt.want, err)
			}
		})

	}
}
