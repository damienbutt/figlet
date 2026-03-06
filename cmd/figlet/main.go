package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	figlet "github.com/damienbutt/figlet"
	"github.com/damienbutt/figlet/internal/version"
)

func main() {
	app := &cli.Command{
		Name:      "figlet",
		Usage:     "Go FIGdriver. Generates ASCII art from text using FIGlet fonts.",
		ArgsUsage: "[text]",
		Version:   version.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "font",
				Aliases: []string{"f"},
				Value:   "Standard",
				Usage:   "font to use",
			},
			&cli.IntFlag{
				Name:    "width",
				Aliases: []string{"w"},
				Value:   80,
				Usage:   "output width",
			},
			&cli.StringFlag{
				Name:    "horizontalLayout",
				Aliases: []string{"h"},
				Value:   "default",
				Usage:   "horizontal layout",
			},
			&cli.StringFlag{
				Name:    "verticalLayout",
				Aliases: []string{"v"},
				Value:   "default",
				Usage:   "vertical layout",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list available fonts",
			},
			&cli.StringFlag{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "show font information",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Handle --list
			if cmd.Bool("list") {
				fonts, err := figlet.Fonts()
				if err != nil {
					return err
				}

				fmt.Println("Available fonts:")
				for _, f := range fonts {
					fmt.Printf("  %s\n", f)
				}

				return nil
			}

			// Handle --info <font>
			if fontName := cmd.String("info"); fontName != "" {
				meta, err := figlet.LoadFont(fontName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error loading font '%s': %v\n", fontName, err)
					os.Exit(1)
				}

				fmt.Printf("Font: %s\n", fontName)
				out, _ := json.MarshalIndent(meta, "", "  ")
				fmt.Printf("Options: %s\n", out)

				return nil
			}

			// Require text argument
			if cmd.NArg() == 0 {
				return cli.ShowAppHelp(cmd)
			}

			text := cmd.Args().First()
			font := cmd.String("font")
			width := cmd.Int("width")
			horizontalLayout := cmd.String("horizontalLayout")
			verticalLayout := cmd.String("verticalLayout")

			opts := &figlet.FigletOptions{
				Font:             figlet.FontName(font),
				HorizontalLayout: figlet.KerningMethods(horizontalLayout),
				VerticalLayout:   figlet.KerningMethods(verticalLayout),
				Width:            width,
			}

			result, err := figlet.Text(text, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Font '%s' not found. Use --list to see available fonts.\n", font)
				os.Exit(1)
			}

			fmt.Println(result)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
