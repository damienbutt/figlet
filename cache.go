package figlet

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

var (
	figFonts = make(map[string]*FigletFont)
	mu       sync.RWMutex
)

// storeFigFont stores a parsed font in the cache under the given name.
// Called from ParseFont after successfully parsing .flf data.
func storeFigFont(name string, font *FigletFont) {
	mu.Lock()
	defer mu.Unlock()
	figFonts[name] = font
}

// LoadFont loads a font by name. It first resolves any known aliases via
// getFontName, then checks the in-memory cache, then tries the embedded
// FontFS, and finally falls back to the filesystem at figDefaults.FontPath.
func LoadFont(name string) (*FontMetadata, error) {
	actualName := getFontName(name)

	// Fast path: cache hit
	mu.RLock()
	if f, ok := figFonts[actualName]; ok {
		mu.RUnlock()
		return f.options, nil
	}

	mu.RUnlock()

	// Try embedded FontFS first
	embeddedPath := "fonts/" + actualName + ".flf"
	data, err := fs.ReadFile(FontFS, embeddedPath)
	if err == nil {
		return ParseFont(actualName, string(data))
	}

	// Fallback: filesystem at configured font path
	fsPath := filepath.Join(figDefaults.FontPath, actualName+".flf")
	data, err = os.ReadFile(fsPath)
	if err != nil {
		return nil, fmt.Errorf("font not found: %q (also checked embedded fonts)", name)
	}

	return ParseFont(actualName, string(data))
}

// LoadedFonts returns the names of all currently cached fonts.
func LoadedFonts() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(figFonts))
	for k := range figFonts {
		names = append(names, k)
	}

	return names
}

// ClearLoadedFonts removes all fonts from the cache.
func ClearLoadedFonts() {
	mu.Lock()
	defer mu.Unlock()
	figFonts = make(map[string]*FigletFont)
}
