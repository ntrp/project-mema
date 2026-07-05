package storage

import "testing"

func TestDatabaseStatusUsesGeneratedQuery(t *testing.T) {
	ctx, store := testDBStore(t)

	status, err := store.GetDatabaseStatus(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if status.Type != "PostgreSQL" || status.Version == "" {
		t.Fatalf("database status = %#v", status)
	}
}
