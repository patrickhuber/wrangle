package interpolate

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/stores"
	memstore "github.com/patrickhuber/wrangle/internal/stores/memory"
	"github.com/stretchr/testify/require"
)

type fakeConfigService struct {
	cfg config.Config
}

func (f *fakeConfigService) Get() (config.Config, error) {
	return f.cfg, nil
}

type fakeStoreService struct {
	stores map[string]stores.Store
}

func (f *fakeStoreService) Get(name string) (stores.Store, error) {
	s, ok := f.stores[name]
	if !ok {
		return nil, fmt.Errorf("store '%s' not found", name)
	}
	return s, nil
}

func (f *fakeStoreService) List() ([]stores.Store, error) {
	var list []stores.Store
	for _, s := range f.stores {
		list = append(list, s)
	}
	return list, nil
}

func TestInterpolate(t *testing.T) {
	cases := []struct {
		name        string
		cfg         config.Config
		stores      map[string]stores.Store
		wantEnv     map[string]string
		wantPkgs    []config.Package
		wantErrLike string
	}{
		{
			name: "resolves env variables",
			cfg: config.Config{Spec: config.Spec{
				Stores:      []config.Store{{Name: "secrets"}},
				Environment: map[string]string{"PASSWORD": "((secret/password))"},
			}},
			stores: func() map[string]stores.Store {
				key, _ := stores.ParseKey("secret/password")
				s := memstore.NewStore()
				_ = s.Set(key, "hunter2")
				return map[string]stores.Store{"secrets": s}
			}(),
			wantEnv: map[string]string{"PASSWORD": "hunter2"},
		},
		{
			name: "resolves non-env fields",
			cfg: config.Config{Spec: config.Spec{
				Stores:      []config.Store{{Name: "vars"}},
				Environment: map[string]string{"PLAIN": "nochange"},
				Packages:    []config.Package{{Name: "app", Version: "((pkg/version))"}},
			}},
			stores: func() map[string]stores.Store {
				key, _ := stores.ParseKey("pkg/version")
				s := memstore.NewStore()
				_ = s.Set(key, "1.2.3")
				return map[string]stores.Store{"vars": s}
			}(),
			wantEnv:  map[string]string{"PLAIN": "nochange"},
			wantPkgs: []config.Package{{Name: "app", Version: "1.2.3"}},
		},
		{
			name: "unresolved variables return error",
			cfg: config.Config{Spec: config.Spec{
				Stores:      []config.Store{{Name: "secrets"}},
				Environment: map[string]string{"PASSWORD": "((secret/missing))"},
			}},
			stores: func() map[string]stores.Store {
				s := memstore.NewStore()
				return map[string]stores.Store{"secrets": s}
			}(),
			wantErrLike: "unable to resolve the following variables",
		},
		{
			name: "no stores keeps template intact",
			cfg: config.Config{Spec: config.Spec{
				Environment: map[string]string{"PASSWORD": "((secret/password))"},
			}},
			stores:  map[string]stores.Store{},
			wantEnv: map[string]string{"PASSWORD": "((secret/password))"},
		},
		{
			name: "respects store dependencies",
			cfg: config.Config{Spec: config.Spec{
				Stores: []config.Store{
					{Name: "b"},
					{Name: "a", Dependencies: []string{"b"}},
				},
				Environment: map[string]string{"VAL": "((b/value))"},
			}},
			stores: func() map[string]stores.Store {
				key, _ := stores.ParseKey("b/value")
				s := memstore.NewStore()
				_ = s.Set(key, "from-b")
				return map[string]stores.Store{"a": memstore.NewStore(), "b": s}
			}(),
			wantEnv: map[string]string{"VAL": "from-b"},
		},
		{
			name: "detects dependency cycles",
			cfg: config.Config{Spec: config.Spec{
				Stores: []config.Store{
					{Name: "a", Dependencies: []string{"b"}},
					{Name: "b", Dependencies: []string{"a"}},
				},
				Environment: map[string]string{"VAL": "((a/value))"},
			}},
			stores: func() map[string]stores.Store {
				return map[string]stores.Store{"a": memstore.NewStore(), "b": memstore.NewStore()}
			}(),
			wantErrLike: "store dependency cycle detected",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(&fakeConfigService{cfg: tc.cfg}, &fakeStoreService{stores: tc.stores})
			result, err := svc.Execute()

			if tc.wantErrLike != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrLike)
				return
			}

			require.NoError(t, err)

			if tc.wantEnv != nil {
				require.Equal(t, tc.wantEnv, result.Spec.Environment)
			}

			if tc.wantPkgs != nil {
				require.Equal(t, tc.wantPkgs, result.Spec.Packages)
			}
		})
	}
}

func TestInterpolateNoStoresSkipsResolution(t *testing.T) {
	cfg := config.Config{
		Spec: config.Spec{
			Stores: []config.Store{},
			Environment: map[string]string{
				"PASSWORD": "((secret/password))",
			},
		},
	}

	svc := NewService(&fakeConfigService{cfg: cfg}, &fakeStoreService{stores: map[string]stores.Store{}})

	result, err := svc.Execute()
	require.NoError(t, err)
	require.Equal(t, "((secret/password))", result.Spec.Environment["PASSWORD"])
}
