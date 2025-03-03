// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package text

import (
	"bytes"
	"context"

	"codeberg.org/gruf/go-byteutil"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// FromPlain fulfils FormatFunc by parsing
// the given plaintext input into a FormatResult.
func (f *Formatter) FromPlain(
	ctx context.Context,
	parseMention gtsmodel.ParseMentionFunc,
	authorID string,
	statusID string,
	input string,
) *FormatResult {
	// Initialize standard block parser
	// that wraps result in <p> tags.
	plainTextParser := parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(newPlaintextParser(), 500),
		),
	)

	return f.fromPlain(
		ctx,
		plainTextParser,
		false, // emojiOnly = false
		parseMention,
		authorID,
		statusID,
		input,
	)
}

// FromPlainNoParagraph fulfils FormatFunc by parsing
// the given plaintext input into a FormatResult.
//
// Unlike FromPlain, it will not wrap the resulting
// HTML in <p> tags, making it useful for parsing
// short fragments of text that oughtn't be formally
// wrapped as a paragraph.
func (f *Formatter) FromPlainNoParagraph(
	ctx context.Context,
	parseMention gtsmodel.ParseMentionFunc,
	authorID string,
	statusID string,
	input string,
) *FormatResult {
	// Initialize block parser that
	// doesn't wrap result in <p> tags.
	plainTextParser := parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(newPlaintextParserNoParagraph(), 500),
		),
	)

	return f.fromPlain(
		ctx,
		plainTextParser,
		false, // emojiOnly = false
		parseMention,
		authorID,
		statusID,
		input,
	)
}

// FromPlainEmojiOnly fulfils FormatFunc by parsing
// the given plaintext input into a FormatResult.
//
// Unlike FromPlain, it will only parse emojis with
// the custom renderer, leaving aside mentions and tags.
func (f *Formatter) FromPlainEmojiOnly(
	ctx context.Context,
	parseMention gtsmodel.ParseMentionFunc,
	authorID string,
	statusID string,
	input string,
) *FormatResult {
	// Initialize block parser that
	// doesn't wrap result in <p> tags.
	plainTextParser := parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(newPlaintextParserNoParagraph(), 500),
		),
	)

	return f.fromPlain(
		ctx,
		plainTextParser,
		true, // emojiOnly = true
		parseMention,
		authorID,
		statusID,
		input,
	)
}

// fromPlain parses the given input text
// using the given plainTextParser, and
// returns the result.
func (f *Formatter) fromPlain(
	ctx context.Context,
	plainTextParser parser.Parser,
	emojiOnly bool,
	parseMention gtsmodel.ParseMentionFunc,
	authorID string,
	statusID string,
	input string,
) *FormatResult {
	result := new(FormatResult)

	// Instantiate goldmark parser for
	// plaintext, using custom renderer
	// to add hashtag/mention links.
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithHardWraps(),
		),
		// Use whichever plaintext
		// parser we were passed.
		goldmark.WithParser(plainTextParser),
		goldmark.WithExtensions(
			&customRenderer{
				ctx,
				f.db,
				parseMention,
				authorID,
				statusID,
				emojiOnly,
				result,
			},
			extension.Linkify, // Turns URLs into links.
		),
	)

	// Convert input string to bytes
	// without performing any allocs.
	bInput := byteutil.S2B(input)

	// Parse input into HTML.
	var htmlBytes bytes.Buffer
	if err := md.Convert(
		bInput,
		&htmlBytes,
	); err != nil {
		log.Errorf(ctx, "error formatting plaintext input to HTML: %s", err)
	}

	// Clean and shrink HTML.
	result.HTML = byteutil.B2S(htmlBytes.Bytes())
	result.HTML = SanitizeToHTML(result.HTML)
	result.HTML = MinifyHTML(result.HTML)

	return result
}
