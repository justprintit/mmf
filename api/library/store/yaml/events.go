package yaml

import (
	"fmt"
	"log"
	"strings"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/library/types"
)

func (store *Store) logNewUser(w *types.Library, u *types.User) {
	if u.Username != u.Name {
		log.Printf("%c %s: %s (%s)", '+', "User", u.Path(), u.Username)
	} else {
		log.Printf("%c %s: %s", '+', "User", u.Path())
	}
}

func (store *Store) logNewGroup(w *types.Library, g *types.Group) {
	log.Printf("%c %s: %s", '+', "Group", g.Path())
}

func (store *Store) logUpdate(w *types.Library, before, after interface{}, prefix string, args ...interface{}) {
	if len(args) > 0 {
		prefix = fmt.Sprintf(prefix, args)
	}

	if s0, ok := before.(string); ok {
		// %q
		s1, _ := after.(string)
		if len(s0) > 0 {
			log.Printf("%c %s: %q", '-', prefix, s0)
		} else {
			log.Printf("%c %s: %q", '-', prefix, s1)
		}
	} else {
		// %v
		log.Printf("%c %s: %v", '-', prefix, before)
		log.Printf("%c %s: %v", '+', prefix, after)
	}
}

func (store *Store) logUserUpdate(w *types.Library, u *types.User, field string, before, after interface{}) {
	store.logUpdate(w, before, after, "%s: %s: %s", "User", u.Path(), field)
}

func (store *Store) logGroupUpdate(w *types.Library, g *types.Group, field string, before, after interface{}) {
	store.logUpdate(w, before, after, "%s: %s: %s", "Group", g.Path(), field)
}

func (store *Store) logErrorPrefixed(w *types.Library, err error, prefix string, args ...interface{}) {
	// prefix
	if len(args) > 0 {
		prefix = fmt.Sprintf(prefix, args...)
	}
	prefix = strings.TrimSpace(prefix)

	if err == nil {
		if len(prefix) > 0 {
			log.Printf("<E>%s", prefix)
		}
	} else if v, ok := errors.AsValidator(err); !ok {
		goto single
	} else if ve := v.Errors(); len(ve) == 1 {
		err = ve[0]
		goto single
	} else if len(ve) > 1 {
		s := make([]string, 1, len(ve)+1)

		if len(prefix) > 0 {
			s[0] = fmt.Sprintf("<E>%s", prefix)
		}

		for i, err := range ve {
			s = append(s, fmt.Sprintf("<E> %v: %s", i, err))
		}

		log.Print(strings.Join(s, "\n"))
	}

	return
single:
	// single error
	if len(prefix) > 0 {
		log.Printf("<E>%s: %s", prefix, err)
	} else {
		log.Printf("<E>%s", err)
	}
}

func (store *Store) logError(w *types.Library, err error) {
	store.logErrorPrefixed(w, err, "")
}

func (store *Store) logUserError(w *types.Library, u *types.User, err error) {
	store.logErrorPrefixed(w, err, "%s: %s", "User", u.Path())
}

func (store *Store) logGroupError(w *types.Library, g *types.Group, err error) {
	store.logErrorPrefixed(w, err, "%s: %s", "Group", g.Path())
}
