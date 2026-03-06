package figlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func intPtr(v int) *int { return &v }

// ---------------------------------------------------------------------------
// getSmushingRules
// ---------------------------------------------------------------------------

func TestGetSmushingRules(t *testing.T) {
	tests := []struct {
		name       string
		oldLayout  int
		fullLayout *int
		want       FittingRules
	}{
		{
			// oldLayout=-1, no fullLayout → hLayout=fullWidth, vLayout=fullWidth, no rules
			name:       "oldLayout -1 (full width)",
			oldLayout:  -1,
			fullLayout: nil,
			want: FittingRules{
				HLayout: layoutFullWidth,
				VLayout: layoutFullWidth,
			},
		},
		{
			// oldLayout=0, no fullLayout → hLayout=fitting, vLayout=fullWidth, no rules
			name:       "oldLayout 0 (fitting)",
			oldLayout:  0,
			fullLayout: nil,
			want: FittingRules{
				HLayout: layoutFitting,
				VLayout: layoutFullWidth,
			},
		},
		{
			// oldLayout=15 (bits 8+4+2+1 = hRule4+hRule3+hRule2+hRule1), no fullLayout.
			// Rules are set but hLayoutSet is false → default branch → rules present → controlledSmushing
			name:       "oldLayout 15 (controlled smushing h rules 1-4, no full layout)",
			oldLayout:  15,
			fullLayout: nil,
			want: FittingRules{
				HLayout: layoutControlledSmushing,
				HRule1:  true,
				HRule2:  true,
				HRule3:  true,
				HRule4:  true,
				VLayout: layoutFullWidth,
			},
		},
		{
			// Standard font: oldLayout=15, fullLayout=24463
			// Bit decode of 24463: sets hLayout=lSmushing, hRule1-4, vLayout=lSmushing, vRule1-5
			// Both layout values upgrade from lSmushing → lControlledSmushing because rules are present
			name:       "Standard font (oldLayout=15, fullLayout=24463)",
			oldLayout:  15,
			fullLayout: intPtr(24463),
			want: FittingRules{
				HLayout: layoutControlledSmushing,
				HRule1:  true,
				HRule2:  true,
				HRule3:  true,
				HRule4:  true,
				HRule5:  false,
				HRule6:  false,
				VLayout: layoutControlledSmushing,
				VRule1:  true,
				VRule2:  true,
				VRule3:  true,
				VRule4:  true,
				VRule5:  true,
			},
		},
		{
			// fullLayout=0 decodes to no bits → hLayoutSet=false, no h-rules fired.
			// Falls back to oldLayout=15 default branch: no rules set → lSmushing.
			// vLayoutSet=false, no v-rules → vLayout=lFullWidth.
			// Contrast with fullLayout=nil,oldLayout=15 which fires rule bits → lControlledSmushing.
			name:       "fullLayout=0 overrides oldLayout rule bits",
			oldLayout:  15,
			fullLayout: intPtr(0),
			want: FittingRules{
				HLayout: layoutSmushing,
				VLayout: layoutFullWidth,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSmushingRules(tt.oldLayout, tt.fullLayout)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// getHorizontalFittingRules
// ---------------------------------------------------------------------------

func TestGetHorizontalFittingRules(t *testing.T) {
	meta := FontMetadata{
		FittingRules: FittingRules{
			HLayout: layoutControlledSmushing,
			HRule1:  true,
			HRule2:  true,
			HRule3:  true,
			HRule4:  true,
			HRule5:  false,
			HRule6:  false,
		},
	}

	tests := []struct {
		name   string
		layout KerningMethods
		wantOK bool
		wantFR FittingRules
	}{
		{
			name:   "default — returns font's own h rules",
			layout: KerningDefault,
			wantOK: true,
			wantFR: FittingRules{
				HLayout: layoutControlledSmushing,
				HRule1:  true, HRule2: true, HRule3: true, HRule4: true,
			},
		},
		{
			name:   "full — overrides to fullWidth, clears rules",
			layout: KerningFull,
			wantOK: true,
			wantFR: FittingRules{HLayout: layoutFullWidth},
		},
		{
			name:   "fitted — overrides to fitting",
			layout: KerningFitted,
			wantOK: true,
			wantFR: FittingRules{HLayout: layoutFitting},
		},
		{
			name:   "controlled smushing — all h rules enabled",
			layout: KerningControlledSmushing,
			wantOK: true,
			wantFR: FittingRules{
				HLayout: layoutControlledSmushing,
				HRule1:  true, HRule2: true, HRule3: true,
				HRule4: true, HRule5: true, HRule6: true,
			},
		},
		{
			name:   "universal smushing",
			layout: KerningUniversalSmushing,
			wantOK: true,
			wantFR: FittingRules{HLayout: layoutSmushing},
		},
		{
			name:   "unrecognised value → false",
			layout: "nonsense",
			wantOK: false,
			wantFR: FittingRules{},
		},
		{
			name:   "empty string → false",
			layout: "",
			wantOK: false,
			wantFR: FittingRules{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getHorizontalFittingRules(tt.layout, meta)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantFR, got)
		})
	}
}

// ---------------------------------------------------------------------------
// getVerticalFittingRules
// ---------------------------------------------------------------------------

func TestGetVerticalFittingRules(t *testing.T) {
	meta := FontMetadata{
		FittingRules: FittingRules{
			VLayout: layoutControlledSmushing,
			VRule1:  true,
			VRule2:  true,
			VRule3:  true,
			VRule4:  true,
			VRule5:  true,
		},
	}

	tests := []struct {
		name   string
		layout KerningMethods
		wantOK bool
		wantFR FittingRules
	}{
		{
			name:   "default — returns font's own v rules",
			layout: KerningDefault,
			wantOK: true,
			wantFR: FittingRules{
				VLayout: layoutControlledSmushing,
				VRule1:  true, VRule2: true, VRule3: true, VRule4: true, VRule5: true,
			},
		},
		{
			name:   "full — overrides to fullWidth",
			layout: KerningFull,
			wantOK: true,
			wantFR: FittingRules{VLayout: layoutFullWidth},
		},
		{
			name:   "fitted",
			layout: KerningFitted,
			wantOK: true,
			wantFR: FittingRules{VLayout: layoutFitting},
		},
		{
			name:   "controlled smushing — all v rules enabled",
			layout: KerningControlledSmushing,
			wantOK: true,
			wantFR: FittingRules{
				VLayout: layoutControlledSmushing,
				VRule1:  true, VRule2: true, VRule3: true, VRule4: true, VRule5: true,
			},
		},
		{
			name:   "universal smushing",
			layout: KerningUniversalSmushing,
			wantOK: true,
			wantFR: FittingRules{VLayout: layoutSmushing},
		},
		{
			name:   "unrecognised value → false",
			layout: "nonsense",
			wantOK: false,
			wantFR: FittingRules{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getVerticalFittingRules(tt.layout, meta)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantFR, got)
		})
	}
}

// ---------------------------------------------------------------------------
// reworkFontOpts
// ---------------------------------------------------------------------------

func TestReworkFontOpts(t *testing.T) {
	baseMeta := FontMetadata{
		HardBlank:      "$",
		Height:         6,
		PrintDirection: LeftToRight,
		FittingRules: FittingRules{
			HLayout: layoutControlledSmushing,
			HRule1:  true, HRule2: true, HRule3: true, HRule4: true,
			VLayout: layoutControlledSmushing,
			VRule1:  true, VRule2: true, VRule3: true, VRule4: true, VRule5: true,
		},
	}

	t.Run("nil opts → Width=-1, rest inherited from meta", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, nil)
		assert.Equal(t, -1, out.Width)
		assert.Equal(t, baseMeta.FittingRules, out.FittingRules)
		assert.Equal(t, LeftToRight, out.PrintDirection)
	})

	t.Run("Width=0 is normalised to -1 (no limit)", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{Width: 0})
		assert.Equal(t, -1, out.Width)
	})

	t.Run("explicit Width is preserved", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{Width: 80})
		assert.Equal(t, 80, out.Width)
	})

	t.Run("ShowHardBlanks propagated", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{ShowHardBlanks: true})
		assert.True(t, out.ShowHardBlanks)
	})

	t.Run("WhitespaceBreak propagated", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{WhitespaceBreak: true})
		assert.True(t, out.WhitespaceBreak)
	})

	t.Run("HorizontalLayout full overrides h rules to fullWidth", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{HorizontalLayout: KerningFull})
		assert.Equal(t, layoutFullWidth, out.FittingRules.HLayout)
		assert.False(t, out.FittingRules.HRule1)
	})

	t.Run("HorizontalLayout default restores font h rules", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{HorizontalLayout: KerningDefault})
		assert.Equal(t, layoutControlledSmushing, out.FittingRules.HLayout)
		assert.True(t, out.FittingRules.HRule1)
	})

	t.Run("VerticalLayout full overrides v rules to fullWidth", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{VerticalLayout: KerningFull})
		assert.Equal(t, layoutFullWidth, out.FittingRules.VLayout)
		assert.False(t, out.FittingRules.VRule1)
	})

	t.Run("PrintDirection DefaultDirection does not override meta value", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{PrintDirection: DefaultDirection})
		assert.Equal(t, LeftToRight, out.PrintDirection)
	})

	t.Run("PrintDirection RightToLeft overrides meta value", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{PrintDirection: RightToLeft})
		assert.Equal(t, RightToLeft, out.PrintDirection)
	})

	t.Run("unrecognised HorizontalLayout leaves meta rules unchanged", func(t *testing.T) {
		out := reworkFontOpts(baseMeta, &FigletOptions{HorizontalLayout: "garbage"})
		assert.Equal(t, baseMeta.FittingRules.HLayout, out.FittingRules.HLayout)
		assert.Equal(t, baseMeta.FittingRules.HRule1, out.FittingRules.HRule1)
	})
}
