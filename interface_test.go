package mockdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopulateInterface(t *testing.T) {

	type person struct {
		ID    int
		Name  string
		Email string
	}

	john := person{
		ID:    1,
		Name:  "John Connor",
		Email: "john@connor.com",
	}

	sarah := person{
		ID:    2,
		Name:  "Sarah Connor",
		Email: "sarah@sky.net",
	}

	target := person{}

	// Populate directly from object
	populateInterface(john, &target)
	require.Equal(t, 1, target.ID)
	require.Equal(t, "John Connor", target.Name)
	require.Equal(t, "john@connor.com", target.Email)

	// Overwrite and populate from pointer
	populateInterface(&sarah, &target)
	require.Equal(t, 2, target.ID)
	require.Equal(t, "Sarah Connor", target.Name)
	require.Equal(t, "sarah@sky.net", target.Email)
}
