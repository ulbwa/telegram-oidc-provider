package entity

import (
	"net/url"
)

type User struct {
	FirstName string
	LastName  *string
	Username  *string
	PhotoUrl  *url.URL
	IsPremium *bool
}

func NewUser(firstName string, lastName, username *string, photoUrl *url.URL, isPremium *bool) (*User, error) {
	if err := validateUserFirstName(firstName); err != nil {
		return nil, err
	}
	if lastName != nil {
		if err := validateUserLastName(*lastName); err != nil {
			return nil, err
		}
	}
	if username != nil {
		if err := validateUserUsername(*username); err != nil {
			return nil, err
		}
	}
	if photoUrl != nil {
		if err := validateUrl(photoUrl); err != nil {
			return nil, err
		}
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		PhotoUrl:  photoUrl,
		IsPremium: isPremium,
	}, nil
}

func (u *User) FullName() string {
	if u.LastName != nil {
		return u.FirstName + " " + *u.LastName
	}
	return u.FirstName
}

func (u *User) WithFirstName(firstName string) (*User, error) {
	if err := validateUserFirstName(firstName); err != nil {
		return nil, err
	}
	newUser := *u
	newUser.FirstName = firstName
	return &newUser, nil
}

func (u *User) WithLastName(lastName *string) (*User, error) {
	if lastName != nil {
		if err := validateUserLastName(*lastName); err != nil {
			return nil, err
		}
	}
	newUser := *u
	newUser.LastName = lastName
	return &newUser, nil
}

func (u *User) WithUsername(username *string) (*User, error) {
	if username != nil {
		if err := validateUserUsername(*username); err != nil {
			return nil, err
		}
	}
	newUser := *u
	newUser.Username = username
	return &newUser, nil
}

func (u *User) WithPhotoUrl(photoUrl *url.URL) (*User, error) {
	if photoUrl != nil {
		if err := validateUrl(photoUrl); err != nil {
			return nil, err
		}
	}
	newUser := *u
	newUser.PhotoUrl = photoUrl
	return &newUser, nil
}

func (u *User) WithIsPremium(isPremium *bool) *User {
	newUser := *u
	newUser.IsPremium = isPremium
	return &newUser
}
