package mockdb

import (
	"testing"

	"github.com/benpate/data/journal"
	"github.com/benpate/exp"
	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {

	john := testPerson{
		Name:  "John Connor",
		Email: "john@connor.com",
		Age:   42,
		Journal: journal.Journal{
			CreateDate: 0,
		},
	}

	criteria := exp.New("name", "=", "John Connor")
	fn := MatcherFunc(&john)
	assert.True(t, fn(criteria))
}

func TestMatch(t *testing.T) {

	john := testPerson{
		PersonID: "42",
		Name:     "John Connor",
		Email:    "john@connor.com",
		Age:      42,
		Journal: journal.Journal{
			CreateDate: 0,
		},
	}

	sarah := testPerson{
		PersonID: "43",
		Name:     "Sarah Connor",
		Email:    "sarah@sky.net",
		Age:      43,
		Journal: journal.Journal{
			CreateDate: 1,
		},
	}

	kyle := testPerson{
		PersonID: "44",
		Name:     "Kyle Reese",
		Email:    "kyle@resistance.org",
		Age:      44,
		Journal: journal.Journal{
			CreateDate: 2,
		},
	}

	{
		// Test INTEGER equality
		exp := exp.New("_id", "=", "42")
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER inequality
		exp := exp.New("_id", "!=", "44")
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER less than
		exp := exp.New("age", "<", 44)
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER greater than
		exp := exp.New("age", ">", 44)
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER greater or equal
		exp := exp.New("age", ">=", 44)
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.True(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER less or equal
		exp6 := exp.New("age", "<=", 43)
		assert.True(t, exp6.Match(MatcherFunc(&john)))
		assert.True(t, exp6.Match(MatcherFunc(&sarah)))
		assert.False(t, exp6.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INTEGER type mismatch
		exp6 := exp.New("_id", ">=", "Michael Jackson")
		assert.False(t, exp6.Match(MatcherFunc(&john)))
		assert.False(t, exp6.Match(MatcherFunc(&sarah)))
		assert.False(t, exp6.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 equality
		exp := exp.New("journal.createDate", "=", 0)
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 inequality
		exp := exp.New("journal.createDate", "!=", 1)
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.True(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 less than
		exp := exp.New("journal.createDate", "<", 2)
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 greater than
		exp := exp.New("journal.createDate", ">", 1)
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.True(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 greater or equal
		exp := exp.New("journal.createDate", ">=", 1)
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.True(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 less or equal
		exp := exp.New("journal.createDate", "<=", 3)
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.True(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test INT64 comparisons
		assert.True(t, exp.New("journal.createDate", "=", int64(0)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "!=", int64(1)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<", int64(1)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<=", int64(0)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<=", int64(1)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">", int64(-1)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">=", int64(-1)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">=", int64(0)).Match(MatcherFunc(&john)))

		assert.False(t, exp.New("journal.createDate", "=", int64(1)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "!=", int64(0)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "<", int64(-1)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "<=", int64(-1)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", ">", int64(1)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", ">=", int64(1)).Match(MatcherFunc(&john)))
	}

	{
		// Test INT64 type mismatch
		exp := exp.New("journal.createDate", "<=", "STRING")
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	// Test multiple fields
	{
		exp := exp.New("name", "=", "John Connor").AndEqual("_id", "42")
		assert.True(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test multiple fields
		exp := exp.New("name", ">", "John Connor").AndLessThan("_id", "44")
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))

		// Test pointers
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.True(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test missing fields
		exp := exp.New("missing-field", ">", "John Connor")
		assert.False(t, exp.Match(MatcherFunc(&john)))
		assert.False(t, exp.Match(MatcherFunc(&sarah)))
		assert.False(t, exp.Match(MatcherFunc(&kyle)))
	}

	{
		// Test missing fields
		exp1 := exp.New("_id", "=", "John Connor")
		assert.False(t, exp1.Match(MatcherFunc(&john)))
		assert.False(t, exp1.Match(MatcherFunc(&sarah)))
		assert.False(t, exp1.Match(MatcherFunc(&kyle)))
	}

	{
		// Test string comparisons
		assert.True(t, exp.New("name", "=", "John Connor").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", ">=", "John Connor").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", "<=", "John Connor").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", "!=", "A").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", "<", "Klaus").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", "<=", "Kaus").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", ">", "Ignacio").Match(MatcherFunc(&john)))
		assert.True(t, exp.New("name", ">=", "Ignacio").Match(MatcherFunc(&john)))

		assert.False(t, exp.New("name", "=", "Sarah Connor").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", "<", "John Connor").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", ">", "John Connor").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", ">=", "Klaus").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", "<=", "Ignacio").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", "!=", "John Connor").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", "<", "Ignacio").Match(MatcherFunc(&john)))
		assert.False(t, exp.New("name", ">", "Klaus").Match(MatcherFunc(&john)))
	}

	{
		// Test INT / INT64 type mismatch
		assert.True(t, exp.New("age", "=", int64(42)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "=", 0).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", "=", int64(43)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "=", 1).Match(MatcherFunc(&john)))

		assert.True(t, exp.New("age", "<", int64(43)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<", 1).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", "<", int64(40)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "<", 0).Match(MatcherFunc(&john)))

		assert.True(t, exp.New("age", "<=", int64(42)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<=", 0).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", "<=", int64(40)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "<=", -1).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("age", "<=", int64(43)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", "<=", 1).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", "<=", int64(40)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", "<=", -1).Match(MatcherFunc(&john)))

		assert.True(t, exp.New("age", ">", int64(40)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">", -1).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", ">", int64(44)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", ">", 1).Match(MatcherFunc(&john)))

		assert.True(t, exp.New("age", ">=", int64(42)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">=", 0).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", ">=", int64(43)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", ">=", 1).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("age", ">=", int64(42)).Match(MatcherFunc(&john)))
		assert.True(t, exp.New("journal.createDate", ">=", 0).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("age", ">=", int64(44)).Match(MatcherFunc(&john)))
		assert.False(t, exp.New("journal.createDate", ">=", 1).Match(MatcherFunc(&john)))
	}
}
