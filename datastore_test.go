package mockdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatastore(t *testing.T) {
	ds1 := New()
	s1, _ := ds1.Session(context.TODO())
	require.NotNil(t, s1)
	s1.Close()
}

func TestSampleDataset(t *testing.T) {
	ds2 := getSampleDataset()
	s2, _ := ds2.Session(context.TODO())
	require.NotNil(t, s2)
	s2.Close()
}

func getSampleDataset() Server {

	return Server{
		"Person": {
			&testPerson{PersonID: "michael", Name: "Michael Jackson", Email: "mike@jackson.com"},
			&testPerson{PersonID: "jermaine", Name: "Jermaine Jackson", Email: "jer@jackson.com"},
			&testPerson{PersonID: "latoya", Name: "Latoya Jackson", Email: "lat@jackson.com"},
			&testPerson{PersonID: "janet", Name: "Janet Jackson", Email: "jan@jackson.com"},
		},
	}
}
