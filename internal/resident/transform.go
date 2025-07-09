package resident

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/iamsuudi/digital-id-server/database/sqlc"
)

func toPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: !t.IsZero()}
}

func toPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toGender(s string) sqlc.Gender {
	return sqlc.Gender(s)
}
func toReligion(s string) sqlc.Religion {
	return sqlc.Religion(s)
}
func toMarital(s string) sqlc.MaritalStatus {
	return sqlc.MaritalStatus(s)
}
func toDocumentType(s string) sqlc.DocumentType {
	return sqlc.DocumentType(s)
}
