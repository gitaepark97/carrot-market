package util

import "database/sql"

func CreateNullableString(val string) sql.NullString {
	if val == "" {
		return sql.NullString{
			String: val,
			Valid:  false,
		}
	}

	return sql.NullString{
		String: val,
		Valid:  true,
	}
}

func CreateNullableInt32(val *int32) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}

	return sql.NullInt32{
		Int32: *val,
		Valid: true,
	}
}
