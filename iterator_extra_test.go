package mockdb

import (
	"sort"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/stretchr/testify/require"
)

func TestIterator_ResetCountCloseError(t *testing.T) {

	it := NewIterator(getTestData())

	require.Equal(t, 10, it.Count())
	require.Nil(t, it.Error())

	// Advance the iterator
	person := testPerson{}
	require.True(t, it.Next(&person))
	require.Equal(t, 1, it.Counter)

	// Reset returns to the beginning
	it.Reset()
	require.Equal(t, 0, it.Counter)

	// Close prevents any further reads
	require.Nil(t, it.Close())
	require.False(t, it.Next(&person))
}

func TestIterator_SortDescending(t *testing.T) {

	it := NewIterator(getTestData(), option.SortDesc("name"))
	sort.Sort(it)

	person := testPerson{}
	require.True(t, it.Next(&person))
	require.Equal(t, "Stonewall Jackson", person.Name)
}

func TestIterator_SecondarySort(t *testing.T) {

	// Two records share the same Age, so the secondary sort (name) decides order.
	records := []data.Object{
		&testPerson{PersonID: "1", Name: "Zachary", Age: 30},
		&testPerson{PersonID: "2", Name: "Adam", Age: 30},
		&testPerson{PersonID: "3", Name: "Mike", Age: 20},
	}

	it := NewIterator(records, option.SortAsc("age"), option.SortAsc("name"))
	sort.Sort(it)

	person := testPerson{}

	require.True(t, it.Next(&person))
	require.Equal(t, "Mike", person.Name) // age 20 sorts first

	require.True(t, it.Next(&person))
	require.Equal(t, "Adam", person.Name) // age 30, name "Adam" before "Zachary"

	require.True(t, it.Next(&person))
	require.Equal(t, "Zachary", person.Name)
}

func TestIterator_Less_RangeGuards(t *testing.T) {

	it := NewIterator(getTestData(), option.SortAsc("name"))

	// Out-of-range indexes return false rather than panicking
	require.False(t, it.Less(999, 0))
	require.False(t, it.Less(0, 999))
}

func TestIterator_Less_NoSortOptions(t *testing.T) {

	// With no sort options, no element sorts before another
	it := NewIterator(getTestData())
	require.False(t, it.Less(0, 1))
}

func TestIterator_Swap(t *testing.T) {

	records := []data.Object{
		&testPerson{PersonID: "A", Name: "First"},
		&testPerson{PersonID: "B", Name: "Second"},
	}

	it := NewIterator(records)
	it.Swap(0, 1)

	require.Equal(t, "B", it.Data[0].ID())
	require.Equal(t, "A", it.Data[1].ID())
}

func TestSafeFieldInterface_ByFieldName(t *testing.T) {

	person := testPerson{Name: "Joe Jackson", Age: 42}

	// Match by struct field name (case-insensitive)
	value, ok := safeFieldInterface(&person, "Age")
	require.True(t, ok)
	require.Equal(t, 42, value)
}

func TestSafeFieldInterface_ByBsonTag(t *testing.T) {

	person := testPerson{Email: "joe@jackson.com"}

	// Match by bson tag
	value, ok := safeFieldInterface(&person, "email")
	require.True(t, ok)
	require.Equal(t, "joe@jackson.com", value)
}

func TestSafeFieldInterface_NotFound(t *testing.T) {
	person := testPerson{}
	_, ok := safeFieldInterface(&person, "nonexistent")
	require.False(t, ok)
}

func TestSafeFieldInterface_Nil(t *testing.T) {
	_, ok := safeFieldInterface(nil, "name")
	require.False(t, ok)
}

func TestSafeFieldInterface_NotStruct(t *testing.T) {
	// A non-struct value has no fields
	_, ok := safeFieldInterface("just a string", "name")
	require.False(t, ok)
}
