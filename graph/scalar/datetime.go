package scalar

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type DateTime time.Time

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *DateTime) UnmarshalGQL(v interface{}) error {
	dt, ok := v.(string)
	if !ok {
		return fmt.Errorf("DateTime must be a string")
	}

	godt, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return err
	}

	*t = DateTime(godt)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t DateTime) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(strconv.Quote(time.Time(t).Format(time.RFC3339))))
	//w.Write([]byte(strconv.Quote(time.Time(t).UTC().Format(time.RFC3339))))
}
