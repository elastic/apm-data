package modelpb

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestStacktraceToModelJSON(t *testing.T) {
	vars, varsMap := randomStruct(t)

	testCases := map[string]struct {
		proto    *StacktraceFrame
		expected *modeljson.StacktraceFrame
	}{
		"empty": {
			proto:    &StacktraceFrame{},
			expected: &modeljson.StacktraceFrame{},
		},
		"no pointers": {
			proto: &StacktraceFrame{
				Filename:       "frame1_filename",
				Classname:      "frame1_classname",
				ContextLine:    "frame1_contextline",
				Module:         "frame1_module",
				Function:       "frame1_function",
				AbsPath:        "frame1_abspath",
				SourcemapError: "frame1_sourcemaperror",
				Original: &Original{
					AbsPath:      "orig1_abspath",
					Filename:     "orig1_filename",
					Classname:    "orig1_classname",
					Function:     "orig1_function",
					LibraryFrame: true,
				},
				LibraryFrame:        true,
				SourcemapUpdated:    true,
				ExcludeFromGrouping: true,
			},
			expected: &modeljson.StacktraceFrame{
				Sourcemap: &modeljson.StacktraceFrameSourcemap{
					Error:   "frame1_sourcemaperror",
					Updated: true,
				},
				Line: &modeljson.StacktraceFrameLine{
					Context: "frame1_contextline",
				},
				Filename:  "frame1_filename",
				Classname: "frame1_classname",
				Module:    "frame1_module",
				Function:  "frame1_function",
				AbsPath:   "frame1_abspath",
				Original: &modeljson.StacktraceFrameOriginal{
					AbsPath:      "orig1_abspath",
					Filename:     "orig1_filename",
					Classname:    "orig1_classname",
					Function:     "orig1_function",
					LibraryFrame: true,
				},
				LibraryFrame:        true,
				ExcludeFromGrouping: true,
			},
		},
		"full": {
			proto: &StacktraceFrame{
				Vars:           vars,
				Lineno:         uintPtr(1),
				Colno:          uintPtr(2),
				Filename:       "frame_filename",
				Classname:      "frame_classname",
				ContextLine:    "frame_contextline",
				Module:         "frame_module",
				Function:       "frame_function",
				AbsPath:        "frame_abspath",
				SourcemapError: "frame_sourcemaperror",
				Original: &Original{
					AbsPath:      "orig_abspath",
					Filename:     "orig_filename",
					Classname:    "orig_classname",
					Lineno:       uintPtr(3),
					Colno:        uintPtr(4),
					Function:     "orig_function",
					LibraryFrame: true,
				},
				PreContext:          []string{"pre"},
				PostContext:         []string{"post"},
				LibraryFrame:        true,
				SourcemapUpdated:    true,
				ExcludeFromGrouping: true,
			},
			expected: &modeljson.StacktraceFrame{
				Sourcemap: &modeljson.StacktraceFrameSourcemap{
					Error:   "frame_sourcemaperror",
					Updated: true,
				},
				Vars: varsMap,
				Line: &modeljson.StacktraceFrameLine{
					Number:  uintPtr(1),
					Column:  uintPtr(2),
					Context: "frame_contextline",
				},
				Filename:  "frame_filename",
				Classname: "frame_classname",
				Module:    "frame_module",
				Function:  "frame_function",
				AbsPath:   "frame_abspath",
				Original: &modeljson.StacktraceFrameOriginal{
					AbsPath:      "orig_abspath",
					Filename:     "orig_filename",
					Classname:    "orig_classname",
					Lineno:       uintPtr(3),
					Colno:        uintPtr(4),
					Function:     "orig_function",
					LibraryFrame: true,
				},
				Context: &modeljson.StacktraceFrameContext{
					Pre:  []string{"pre"},
					Post: []string{"post"},
				},
				LibraryFrame:        true,
				ExcludeFromGrouping: true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.StacktraceFrame
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}

}
