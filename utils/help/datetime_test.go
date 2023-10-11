package help

import (
	"database/sql"
	"testing"
)

func TestParseUnixWithNullTime(t *testing.T) {
	t.Log(ParseUnixWithNullTime(0))
	t.Log(ParseUnixWithNullTime(1666237264518))
	t.Log(ParseUnixWithNullTime(1666237223))
}

func TestNullTimeToUnix(t *testing.T) {
	t.Log(NullTimeToUnix(sql.NullTime{Time: ParseUnix(1666237223), Valid: true}))
}
