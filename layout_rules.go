package figlet

// Layout modes used throughout the smushing and fitting logic.
const (
	layoutFullWidth int = iota
	layoutFitting
	layoutSmushing
	layoutControlledSmushing
)

// getSmushingRules decodes the bit-packed oldLayout / fullLayout integers
// from a FIGlet font header into a FittingRules value.
func getSmushingRules(oldLayout int, fullLayout *int) FittingRules {
	codes := []struct {
		code  int
		field FittingProperties
		val   any
	}{
		{16384, FitVLayout, layoutSmushing},
		{8192, FitVLayout, layoutFitting},
		{4096, FitVRule5, true},
		{2048, FitVRule4, true},
		{1024, FitVRule3, true},
		{512, FitVRule2, true},
		{256, FitVRule1, true},
		{128, FitHLayout, layoutSmushing},
		{64, FitHLayout, layoutFitting},
		{32, FitHRule6, true},
		{16, FitHRule5, true},
		{8, FitHRule4, true},
		{4, FitHRule3, true},
		{2, FitHRule2, true},
		{1, FitHRule1, true},
	}

	var (
		hLayout, vLayout       int
		hLayoutSet, vLayoutSet bool
		hRule1, hRule2, hRule3 bool
		hRule4, hRule5, hRule6 bool
		vRule1, vRule2, vRule3 bool
		vRule4, vRule5         bool
	)

	val := oldLayout
	if fullLayout != nil {
		val = *fullLayout
	}

	for _, c := range codes {
		if val >= c.code {
			val -= c.code

			switch c.field {
			case FitHLayout:
				if !hLayoutSet {
					hLayout = c.val.(int)
					hLayoutSet = true
				}
			case FitVLayout:
				if !vLayoutSet {
					vLayout = c.val.(int)
					vLayoutSet = true
				}
			case FitHRule1:
				hRule1 = true
			case FitHRule2:
				hRule2 = true
			case FitHRule3:
				hRule3 = true
			case FitHRule4:
				hRule4 = true
			case FitHRule5:
				hRule5 = true
			case FitHRule6:
				hRule6 = true
			case FitVRule1:
				vRule1 = true
			case FitVRule2:
				vRule2 = true
			case FitVRule3:
				vRule3 = true
			case FitVRule4:
				vRule4 = true
			case FitVRule5:
				vRule5 = true
			}
		}
	}

	// Resolve hLayout if not set by bit decoding
	if !hLayoutSet {
		switch oldLayout {
		case 0:
			hLayout = layoutFitting
		case -1:
			hLayout = layoutFullWidth
		default:
			if hRule1 || hRule2 || hRule3 || hRule4 || hRule5 || hRule6 {
				hLayout = layoutControlledSmushing
			} else {
				hLayout = layoutSmushing
			}
		}
	} else if hLayout == layoutSmushing {
		if hRule1 || hRule2 || hRule3 || hRule4 || hRule5 || hRule6 {
			hLayout = layoutControlledSmushing
		}
	}

	// Resolve vLayout if not set by bit decoding
	if !vLayoutSet {
		if vRule1 || vRule2 || vRule3 || vRule4 || vRule5 {
			vLayout = layoutControlledSmushing
		} else {
			vLayout = layoutFullWidth
		}
	} else if vLayout == layoutSmushing {
		if vRule1 || vRule2 || vRule3 || vRule4 || vRule5 {
			vLayout = layoutControlledSmushing
		}
	}

	return FittingRules{
		HLayout: hLayout,
		HRule1:  hRule1,
		HRule2:  hRule2,
		HRule3:  hRule3,
		HRule4:  hRule4,
		HRule5:  hRule5,
		HRule6:  hRule6,
		VLayout: vLayout,
		VRule1:  vRule1,
		VRule2:  vRule2,
		VRule3:  vRule3,
		VRule4:  vRule4,
		VRule5:  vRule5,
	}
}

// getHorizontalFittingRules returns the FittingRules override for a
// user-requested horizontal layout mode. Returns zero-value and false if
// layout is empty or unrecognised.
func getHorizontalFittingRules(layout KerningMethods, opts FontMetadata) (FittingRules, bool) {
	fr := opts.FittingRules

	switch layout {
	case KerningDefault:
		return FittingRules{
			HLayout: fr.HLayout,
			HRule1:  fr.HRule1,
			HRule2:  fr.HRule2,
			HRule3:  fr.HRule3,
			HRule4:  fr.HRule4,
			HRule5:  fr.HRule5,
			HRule6:  fr.HRule6,
		}, true
	case KerningFull:
		return FittingRules{HLayout: layoutFullWidth}, true
	case KerningFitted:
		return FittingRules{HLayout: layoutFitting}, true
	case KerningControlledSmushing:
		return FittingRules{
			HLayout: layoutControlledSmushing,
			HRule1:  true, HRule2: true, HRule3: true,
			HRule4: true, HRule5: true, HRule6: true,
		}, true
	case KerningUniversalSmushing:
		return FittingRules{HLayout: layoutSmushing}, true
	default:
		return FittingRules{}, false
	}
}

// getVerticalFittingRules returns the FittingRules override for a
// user-requested vertical layout mode.
func getVerticalFittingRules(layout KerningMethods, opts FontMetadata) (FittingRules, bool) {
	fr := opts.FittingRules

	switch layout {
	case KerningDefault:
		return FittingRules{
			VLayout: fr.VLayout,
			VRule1:  fr.VRule1,
			VRule2:  fr.VRule2,
			VRule3:  fr.VRule3,
			VRule4:  fr.VRule4,
			VRule5:  fr.VRule5,
		}, true
	case KerningFull:
		return FittingRules{VLayout: layoutFullWidth}, true
	case KerningFitted:
		return FittingRules{VLayout: layoutFitting}, true
	case KerningControlledSmushing:
		return FittingRules{
			VLayout: layoutControlledSmushing,
			VRule1:  true, VRule2: true, VRule3: true, VRule4: true, VRule5: true,
		}, true
	case KerningUniversalSmushing:
		return FittingRules{VLayout: layoutSmushing}, true
	default:
		return FittingRules{}, false
	}
}

// reworkFontOpts merges user-supplied FigletOptions on top of font metadata
// to produce the InternalOptions used throughout rendering.
func reworkFontOpts(meta FontMetadata, opts *FigletOptions) InternalOptions {
	result := InternalOptions{FontMetadata: meta}

	if opts == nil {
		result.Width = -1
		return result
	}

	result.ShowHardBlanks = opts.ShowHardBlanks
	result.Width = opts.Width
	if result.Width == 0 {
		result.Width = -1
	}

	result.WhitespaceBreak = opts.WhitespaceBreak

	if opts.HorizontalLayout != "" {
		if params, ok := getHorizontalFittingRules(opts.HorizontalLayout, meta); ok {
			result.FittingRules.HLayout = params.HLayout
			result.FittingRules.HRule1 = params.HRule1
			result.FittingRules.HRule2 = params.HRule2
			result.FittingRules.HRule3 = params.HRule3
			result.FittingRules.HRule4 = params.HRule4
			result.FittingRules.HRule5 = params.HRule5
			result.FittingRules.HRule6 = params.HRule6
		}
	}

	if opts.VerticalLayout != "" {
		if params, ok := getVerticalFittingRules(opts.VerticalLayout, meta); ok {
			result.FittingRules.VLayout = params.VLayout
			result.FittingRules.VRule1 = params.VRule1
			result.FittingRules.VRule2 = params.VRule2
			result.FittingRules.VRule3 = params.VRule3
			result.FittingRules.VRule4 = params.VRule4
			result.FittingRules.VRule5 = params.VRule5
		}
	}

	if opts.PrintDirection != DefaultDirection {
		result.PrintDirection = opts.PrintDirection
	}

	return result
}
