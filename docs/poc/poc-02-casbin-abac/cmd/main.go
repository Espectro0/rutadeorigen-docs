package main

import (
	"fmt"

	"github.com/Espectro0/rutadeorigen-docs/internal/auth"
	"github.com/Espectro0/rutadeorigen-docs/internal/domain"
)

func main() {
	enforcer, err := auth.NewEnforcer()
	if err != nil {
		panic(err)
	}

	admin := domain.User{
		ID:   "u1",
		Role: "admin",
	}
	producer1 := domain.User{
		ID:         "u2",
		Role:       "producer",
		BusinessID: "b1",
	}
	producer2 := domain.User{
		ID:         "u3",
		Role:       "producer",
		BusinessID: "b2",
	}
	consumer := domain.User{
		ID:   "u4",
		Role: "consumer",
	}

	batchDraft := domain.Batch{
		ID:         "batch1",
		BusinessID: "b1",
		Status:     "draft",
	}
	batchPublished := domain.Batch{
		ID:         "batch2",
		BusinessID: "b1",
		Status:     "published",
	}

	cases := []struct {
		description string
		sub         domain.User
		obj         domain.Batch
		act         string
		expected    bool
	}{
		{"Admin edits any batch", admin, batchDraft, "edit", true},
		{"Producer edits their own draft batch", producer1, batchDraft, "edit", true},
		{"Producer edits their own published batch", producer1, batchPublished, "edit", false},
		{"Producer edits another producer's draft batch", producer2, batchDraft, "edit", false},
		{"Consumer edits any batch", consumer, batchDraft, "edit", false},
	}

	fails := 0
	for _, c := range cases {
		ok, err := enforcer.Enforce(c.sub, c.obj, c.act)
		if err != nil {
			panic(err)
		}
		status := "OK"
		if ok != c.expected {
			status = "FAIL"
			fails++
		}
		fmt.Printf("[%s] %-55s -> allowed=%-5v [expected=%-5v]\n", status, c.description, ok, c.expected)
	}

	if fails > 0 {
		fmt.Printf("\n%d test(s) failed.\n", fails)
	} else {
		fmt.Println("\nAll tests passed.")
	}
}
