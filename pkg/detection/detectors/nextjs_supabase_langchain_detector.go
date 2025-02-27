// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.
// Deprecated: This package contains deprecated detectors that will be removed in a future version.
package detectors

// DEPRECATED: This detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology stacks.
// Please migrate any direct references to this detector to use StackDetector instead.
// Planned removal: v1.x.0 (next major/minor release)

import (
	"context"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// NextjsSupabaseLangchainDetector detects Next.js + Supabase + LangChain stack
type NextjsSupabaseLangchainDetector struct {
	detection.BaseDetector
	nextjsDetector    *NextjsDetector
	supabaseDetector  *SupabaseDetector
	langchainDetector *LangchainDetector
	pgvectorDetector  *PgvectorDetector
	tailwindDetector  *TailwindDetector
	stripeDetector    *StripeDetector
}

// NewNextjsSupabaseLangchainDetector creates a new detector for Next.js + Supabase + LangChain stack
func NewNextjsSupabaseLangchainDetector() *NextjsSupabaseLangchainDetector {
	detector := &NextjsSupabaseLangchainDetector{
		nextjsDetector:    NewNextjsDetector(),
		supabaseDetector:  NewSupabaseDetector(),
		langchainDetector: NewLangchainDetector(),
		pgvectorDetector:  NewPgvectorDetector(),
		tailwindDetector:  NewTailwindDetector(),
		stripeDetector:    NewStripeDetector(),
	}
	detector.BaseDetector = *detection.NewBaseDetector("Next.js + Supabase + LangChain Detector", 0.8)
	return detector
}

// Detect implements the Detector interface
func (d *NextjsSupabaseLangchainDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("NextjsSupabaseLangchainDetector")

	// Continue with existing implementation
	info := &detection.ProjectInfo{
		Name:     filepath.Base(dir),
		Type:     "unknown",
		Metadata: make(map[string]interface{}),
	}

	// Detect individual components
	nextjsInfo, _ := d.nextjsDetector.Detect(ctx, dir)
	supabaseInfo, _ := d.supabaseDetector.Detect(ctx, dir)
	langchainInfo, _ := d.langchainDetector.Detect(ctx, dir)

	// Detect optional components
	pgvectorInfo, _ := d.pgvectorDetector.Detect(ctx, dir)
	tailwindInfo, _ := d.tailwindDetector.Detect(ctx, dir)
	stripeInfo, _ := d.stripeDetector.Detect(ctx, dir)

	// If we have Next.js + Supabase + LangChain
	if nextjsInfo.Type == "nextjs" &&
		supabaseInfo.Type == "supabase" &&
		langchainInfo.Type == "langchain" {

		// Calculate combined confidence score
		combinedConf := (nextjsInfo.Confidence + supabaseInfo.Confidence + langchainInfo.Confidence) / 3.0

		// Determine if we have the vector-enabled variant
		stackType := "nextjs-supabase-langchain"
		if pgvectorInfo.Type == "pgvector" && pgvectorInfo.Confidence > 0.7 {
			stackType = "nextjs-supabase-langchain-vector"
		}

		// Set the project info
		info.Type = stackType
		info.Confidence = combinedConf
		info.Language = "JavaScript"
		info.Framework = "Next.js"
		info.LLMProvider = "LangChain"

		// Add metadata about detected components
		info.Metadata["has_pgvector"] = pgvectorInfo.Type == "pgvector" && pgvectorInfo.Confidence > 0.7
		info.Metadata["has_tailwind"] = tailwindInfo.Type == "tailwind" && tailwindInfo.Confidence > 0.7
		info.Metadata["has_stripe"] = stripeInfo.Type == "stripe" && stripeInfo.Confidence > 0.7

		return info, nil
	}

	return info, nil
}
