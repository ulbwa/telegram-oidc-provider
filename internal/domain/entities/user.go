package entities

import (
	"net/url"
	"time"
)

// User represents a Telegram user.
type User struct {
	Id        int64
	FirstName string
	LastName  *string
	Username  *string
	PhotoUrl  *url.URL
	IsPremium *bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// NewUser creates a new User instance with validation.
func NewUser(id int64, firstName string, lastName, username *string, photoUrl *url.URL, isPremium *bool) (*User, error) {
	if err := validateUserId(id); err != nil {
		return nil, err
	}
	if err := validateFirstName(firstName); err != nil {
		return nil, err
	}
	if lastName != nil {
		if err := validateLastName(*lastName); err != nil {
			return nil, err
		}
	}
	if username != nil {
		if err := validateUsername(*username); err != nil {
			return nil, err
		}
	}
	if photoUrl != nil {
		if err := validateUrl(photoUrl); err != nil {
			return nil, err
		}
	}

	return &User{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		PhotoUrl:  photoUrl,
		IsPremium: isPremium,
		CreatedAt: time.Now(),
	}, nil
}

func (u *User) ModifiedAt() time.Time {
	if u.UpdatedAt != nil {
		return *u.UpdatedAt
	}
	return u.CreatedAt
}

func (u *User) FullName() string {
	if u.LastName != nil {
		return u.FirstName + " " + *u.LastName
	}
	return u.FirstName
}

func (u *User) Touch() {
	now := time.Now()
	u.UpdatedAt = &now
}

func (u *User) SetFirstName(firstName string) error {
	if err := validateFirstName(firstName); err != nil {
		return err
	}
	if u.FirstName == firstName {
		return nil
	}
	u.FirstName = firstName
	u.Touch()
	return nil
}

func (u *User) SetLastName(lastName *string) error {
	if lastName != nil {
		if err := validateLastName(*lastName); err != nil {
			return err
		}
	}
	if (u.LastName == nil && lastName == nil) || (u.LastName != nil && lastName != nil && *u.LastName == *lastName) {
		return nil
	}
	u.LastName = lastName
	u.Touch()
	return nil
}

func (u *User) SetUsername(username *string) error {
	if username != nil {
		if err := validateUsername(*username); err != nil {
			return err
		}
	}
	if (u.Username == nil && username == nil) || (u.Username != nil && username != nil && *u.Username == *username) {
		return nil
	}
	u.Username = username
	u.Touch()
	return nil
}

func (u *User) SetPhotoUrl(photoUrl *url.URL) error {
	if (u.PhotoUrl == nil && photoUrl == nil) || (u.PhotoUrl != nil && photoUrl != nil && *u.PhotoUrl == *photoUrl) {
		return nil
	}
	if photoUrl != nil {
		if err := validateUrl(photoUrl); err != nil {
			return err
		}
	}
	u.PhotoUrl = photoUrl
	u.Touch()
	return nil
}

func (u *User) SetIsPremium(isPremium *bool) {
	if (u.IsPremium == nil && isPremium == nil) || (u.IsPremium != nil && isPremium != nil && *u.IsPremium == *isPremium) {
		return
	}
	u.IsPremium = isPremium
	u.Touch()
}
