package storage

import "testing"

func TestDLNARendererProfileRebasePreservesCustomEdits(t *testing.T) {
	ctx, store := testDBStore(t)
	original := requireRendererProfile(t, ctx, store, "vlc")
	input := rendererProfileInput(original)
	input.Name = "VLC Custom"
	input.Notes = "local edits"
	edited, err := store.UpdateDLNARendererProfile(ctx, "vlc", input)
	if err != nil {
		t.Fatalf("update renderer profile: %v", err)
	}
	if !edited.Customized {
		t.Fatalf("edited profile should be customized: %#v", edited)
	}
	_, err = store.pool.Exec(ctx, `
		update app.dlna_renderer_profile_defaults
		set source_version = source_version + 1,
		    name = 'VLC Seed Upgrade'
		where id = 'vlc'
	`)
	if err != nil {
		t.Fatalf("upgrade seed default: %v", err)
	}

	rebased, err := store.RebaseDLNARendererProfile(ctx, "vlc")
	if err != nil {
		t.Fatalf("rebase renderer profile: %v", err)
	}
	if rebased.Name != "VLC Custom" || rebased.Notes != "local edits" || !rebased.Customized {
		t.Fatalf("rebased profile lost custom edits: %#v", rebased)
	}
	if rebased.SourceVersion != original.SourceVersion+1 || rebased.Source != "user" {
		t.Fatalf("rebased source state = %#v", rebased)
	}

	reset, err := store.ResetDLNARendererProfile(ctx, "vlc")
	if err != nil {
		t.Fatalf("reset renderer profile: %v", err)
	}
	if reset.Name != "VLC Seed Upgrade" || reset.Customized {
		t.Fatalf("reset profile = %#v", reset)
	}
}
