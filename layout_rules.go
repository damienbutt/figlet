package figlet

// Layout modes — mirror the internal/layout constants as plain ints to
// avoid cross-package casting everywhere inside this package.
const (
	lFullWidth          = 0
	lFitting            = 1
	lSmushing           = 2
	lControlledSmushing = 3
)

// getSmushingRules decodes the bit-packed oldLayout / fullLayout integers
// from a FIGlet font header into a FittingRules value.
func getSmushingRules(oldLayout int, fullLayout *int) FittingRules {
	type entry struct {
		code  int
		field string
		val   int // 0 = false, 1 = true, layout value otherwise
	}

	codes := []struct {
		code  int
		field string
		val   any
	}{
		{16384, "vLayout", lSmushing},
		{8192, "vLayout", lFitting},
		{4096, "vRule5", true},
		{2048, "vRule4", true},
		{1024, "vRule3", true},
		{512, "vRule2", true},
		{256, "vRule1", true},
		{128, "hLayout", lSmushing},
		{64, "hLayout", lFitting},
		{32, "hRule6", true},
		{16, "hRule5", true},
		{8, "hRule4", true},
		{4, "hRule3", true},
		{2, "hRule2", true},
		{1, "hRule1", true},
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
			case "hLayout":
				if !hLayoutSet {
					hLayout = c.val.(int)
					hLayoutSet = true
				}
			case "vLayout":
				if !vLayoutSet {
					vLayout = c.val.(int)
					vLayoutSet = true
				}
			case "hRule1":
				hRule1 = true
			case "hRule2":
				hRule2 = true
			case "hRule3":
				hRule3 = true
			case "hRule4":
				hRule4 = true
			case "hRule5":
				hRule5 = true
			case "hRule6":
				hRule6 = true
			case "vRule1":
				vRule1 = true
			case "vRule2":
				vRule2 = true
			case "vRule3":
				vRule3 = true
			case "vRule4":
				vRule4 = true
			case "vRule5":
				vRule5 = true
			}
		} else if c.field != "vLayout" && c.field != "hLayout" {
			// rule fields default to false — already the zero value
		}
	}

	// Resolve hLayout if not set by bit decoding
	if !hLayoutSet {
		if oldLayout == 0 {
			hLayout = lFitting
		} else if oldLayout == -1 {
			hLayout = lFullWidth
		} else {
			if hRule1 || hRule2 || hRule3 || hRule4 || hRule5 || hRule6 {
				hLayout = lControlledSmushing
			} else {
				hLayout = lSmushing
			}
		}
	} else if hLayout == lSmushing {
		if hRule1 || hRule2 || hRule3 || hRule4 || hRule5 || hRule6 {
			hLayout = lControlledSmushing
		}
	}

	// Resolve vLayout if not set by bit decoding
	if !vLayoutSet {
		if vRule1 || vRule2 || vRule3 || vRule4 || vRule5 {
			vLayout = lControlledSmushing
		} else {
			vLayout = lFullWidth
		}
	} else if vLayout == lSmushing {
		if vRule1 || vRule2 || vRule3 || vRule4 || vRule5 {
			vLayout = lControlledSmushing
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
	case "default":
		return FittingRules{
			HLayout: fr.HLayout,
			HRule1:  fr.HRule1,
			HRule2:  fr.HRule2,
			HRule3:  fr.HRule3,
			HRule4:  fr.HRule4,
			HRule5:  fr.HRule5,
			HRule6:  fr.HRule6,
		}, true
	case "full":
		return FittingRules{HLayout: lFullWidth}, true
	case "fitted":
		return FittingRules{HLayout: lFitting}, true
	case "controlled smushing":
		return FittingRules{
			HLayout: lControlledSmushing,
			HRule1:  true, HRule2: true, HRule3: true,
			HRule4: true, HRule5: true, HRule6: true,
		}, true
	case "universal smushing":
		return FittingRules{HLayout: lSmushing}, true
	default:
		return FittingRules{}, false
	}
}

// getVerticalFittingRules returns the FittingRules override for a
// user-requested vertical layout mode.
func getVerticalFittingRules(layout KerningMethods, opts FontMetadata) (FittingRules, bool) {
	fr := opts.FittingRules

	switch layout {
	case "default":
		return FittingRules{
			VLayout: fr.VLayout,
			VRule1:  fr.VRule1,
			VRule2:  fr.VRule2,
			VRule3:  fr.VRule3,
			VRule4:  fr.VRule4,
			VRule5:  fr.VRule5,
		}, true
	case "full":
		return FittingRules{VLayout: lFullWidth}, true
	case "fitted":
		return FittingRules{VLayout: lFitting}, true
	case "controlled smushing":
		return FittingRules{
			VLayout: lControlledSmushing,
			VRule1:  true, VRule2: true, VRule3: true, VRule4: true, VRule5: true,
		}, true
	case "universal smushing":
		return FittingRules{VLayout: lSmushing}, true
	default:
		return FittingRules{}, false
	}
}

// reworkFontOpts merges user-supplied FigletOptions on top of font metadata
// to produce the InternalOptions used throughout rendering.
func reworkFontOpts(meta FontMetadata, opts *FigletOptions) InternalOptions {
	myOpts := InternalOptions{FontMetadata: meta}

	if opts == nil {
		myOpts.Width = -1
		return myOpts
	}

	myOpts.ShowHardBlanks = opts.ShowHardBlanks
	myOpts.Width = opts.Width
	if myOpts.Width == 0 {
		myOpts.Width = -1
	}

	myOpts.WhitespaceBreak = opts.WhitespaceBreak

	if opts.HorizontalLayout != "" {
		if params, ok := getHorizontalFittingRules(opts.HorizontalLayout, meta); ok {
			myOpts.FittingRules.HLayout = params.HLayout
			myOpts.FittingRules.HRule1 = params.HRule1
			myOpts.FittingRules.HRule2 = params.HRule2
			myOpts.FittingRules.HRule3 = params.HRule3
			myOpts.FittingRules.HRule4 = params.HRule4
			myOpts.FittingRules.HRule5 = params.HRule5
			myOpts.FittingRules.HRule6 = params.HRule6
		}
	}

	if opts.VerticalLayout != "" {
		if params, ok := getVerticalFittingRules(opts.VerticalLayout, meta); ok {
			myOpts.FittingRules.VLayout = params.VLayout
			myOpts.FittingRules.VRule1 = params.VRule1
			myOpts.FittingRules.VRule2 = params.VRule2
			myOpts.FittingRules.VRule3 = params.VRule3
			myOpts.FittingRules.VRule4 = params.VRule4
			myOpts.FittingRules.VRule5 = params.VRule5
		}
	}

	if opts.PrintDirection != DefaultDirection {
		myOpts.PrintDirection = opts.PrintDirection
	}

	return myOpts
}
